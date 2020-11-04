package Services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	Cfg "gold-store/Config"
	G "gold-store/Globals"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"log"
	"math"
	"strconv"
	"time"
)

func CheckCashPoint(amount float64, reqPoint float64, user Mod.User) bool {
	cashPointData := R.CheckCashPoint(user)
	cashP := amount / cashPointData.Rate

	log.Println("Actual Point:", cashP, "----- Point got from request:", reqPoint)
	if cashP != reqPoint {
		return false
	} else if cashP > cashPointData.CashPoint {
		return false
	} else {
		return true
	}
}

func IsCouponValid(coupon Mod.Coupon, userID uint) bool {
	if coupon.ID > 0 {
		usage := R.CouponUsageCount(userID, coupon.ID)
		if usage <= 0 {
			now, _ := time.Parse("2006-01-02 15:04", time.Now().Format("2006-01-02 15:04"))
			if diffStart := now.Sub(coupon.StartDate.UTC()); diffStart.Seconds() >= 0 {
				if diffEnd := now.Sub(coupon.EndDate.UTC()); diffEnd.Seconds() >= 0 {
					coupon.Status = 0
					R.SaveCoupon(coupon)
					G.Msg.Fail = "The code is not valid or expired."
					return false
				} else {
					return true
				}
			} else {
				G.Msg.Fail = "The code is not valid or expired."
				return false
			}
		} else {
			G.Msg.Fail = "You have used the code already."
			return false
		}
	} else {
		G.Msg.Fail = "The code is not valid or expired."
		return false
	}
}


func ShippingAddressValidation(c *gin.Context) (bool, Mod.OrderShippingAddress) {
	var shipAddress Mod.OrderShippingAddress
	err := c.ShouldBind(&shipAddress)

	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Please fill all the required fields with valid values."
		return false, shipAddress

	}

	_, err = strconv.Atoi(shipAddress.Phone)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Not a valid phone number."
		return false, shipAddress
	}

	_, err = strconv.Atoi(shipAddress.ZipCode)
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Not a valid zip code."
		return false, shipAddress
	}
	return true, shipAddress
}


func CheckoutProcessing(c *gin.Context, user Mod.User, order Mod.Order, fromApi bool) (Mod.Order,bool) {
	payment := c.PostForm("payment")
	if payment == "1" {
		order.PayMethodID = 1
		order.PaidAmount = 0.0
		order := OrderProcessing(c, user, order, fromApi)

		if order.GrandTotal <= 0.0 {
			log.Println("Checkout failed. Check here")
			return order, false
		}

		order, success := CheckoutDBProcessing(c, order, payment, fromApi)
		if success {
			success = SendOrderEmail(order, c, "")
			log.Println("Order Place Email Sending Success:", success)

			return order, true
		} else {
			log.Println("Checkout failed. Check here")
			return order, false
		}
	} else if payment == "2" {
		order := OrderProcessing(c, user, order, fromApi)
		if order.GrandTotal <= 0.0 || order.DeliveryType < 1 || order.DeliveryType > 2 {
			log.Println("Checkout failed. Check here")
			return order, false
		}

		order, success := CheckoutDBProcessing(c, order, payment, fromApi)

		if success {
			/*if order.PayMethodID != 2 && order.DeliveryType == 2 {
				success = SendOrderEmail(order, c, "")
				log.Println("Order Email Sending Success:", success)
			}*/
			return order, true
		} else {
			log.Println("Checkout failed. Check here")
			return order, false
		}
	}
	log.Println("Checkout failed. Check here")
	return order, false
}

func OrderProcessing(c *gin.Context, user Mod.User, order Mod.Order, fromApi bool) Mod.Order {
	var finalCart G.FinalCart


	if !fromApi {
		proID, err := strconv.Atoi(c.PostForm("product_id"))
		if err != nil {
			log.Println(err.Error())
			return order
		}

		quantity, err := strconv.Atoi(c.PostForm("quantity"))
		if err != nil {
			log.Println(err.Error())
			return order
		}

		if c.Request.Header["Referer"][0] != G.AppEnv.Url+"/buy-now" && proID == 0 {
			finalCart = ProcessCart(c)
		} else if c.Request.Header["Referer"][0] == G.AppEnv.Url+"/buy-now" && proID > 0 && quantity > 0 {
			finalCart = BuyNowProcess(uint(proID), c)
		} else {
			return order
		}
	} else {
		finalCart = ProcessCart(c)
	}

	order.User = user

	order.UserID = user.ID
	order.SubTotal = finalCart.SubTotal
	order.GrandTotal = finalCart.GrandTotal
	order.TotalGrmAmount = finalCart.TotalGrmAmount
	order.Fees1Percent = finalCart.Fees1Percent
	order.Fees1Fixed = finalCart.Fees1Fixed
	order.Fees2Percent = finalCart.Fees2Percent
	order.Fees1Fixed = finalCart.Fees2Fixed
	order.Fees3Percent = finalCart.Fees3Percent
	order.Fees3Fixed = finalCart.Fees3Fixed
	order.TotalFees = finalCart.TotalFees
	order.TotalDiscount = finalCart.TotalDiscount

	deliveryType, err := strconv.Atoi(c.PostForm("delivery_type"))
	log.Println(deliveryType)
	if err != nil {
		log.Println(err.Error())
		return order
	}
	if deliveryType < 1 || deliveryType > 2 {
		log.Println("Checkout failed. Check here. Invalid delivery type:", deliveryType)
		return order
	}
	order.DeliveryType = deliveryType

	order = CouponDiscountProcessing(c, order, user)

	for _, cart := range finalCart.Carts {
		var orderDetail Mod.OrderDetail
		orderDetail.OrderID = order.ID
		orderDetail.ProductID = cart.Product.ID
		orderDetail.KaratID = cart.Product.KaratID
		orderDetail.ProductKaratAmount = cart.Product.Karat.Amount
		orderDetail.ProductTitle = cart.Product.Title
		orderDetail.ProductType = cart.Product.Type
		orderDetail.ProductPrice = cart.Product.Price
		orderDetail.ProductGrmAmount = cart.Product.GrmAmount
		if cart.Product.Type == 1 {
			orderDetail.TotalGrmAmount, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", (cart.Quantity * cart.Product.GrmAmount)), 64)
		} else if cart.Product.Type == 2 {
			orderDetail.TotalGrmAmount = cart.Quantity
		}
		orderDetail.ProductImgUrl = cart.Product.ImgUrl
		orderDetail.Quantity = cart.Quantity
		orderDetail.Total = cart.Total
		orderDetail.TotalDiscount = cart.TotalDiscount
		orderDetail.TotalWithDiscount = cart.TotalWithDiscount
		order.OrderDetails = append(order.OrderDetails, orderDetail)
	}
	return order
}

func CouponDiscountProcessing(c *gin.Context, order Mod.Order, user Mod.User) Mod.Order {
	var coupon Mod.Coupon
	code := c.PostForm("coupon_code")
	if code != "" {
		coupon = R.Coupon(coupon, "code = BINARY ? and status=? and type=?",code, 1, 1)
		valid := IsCouponValid(coupon, user.ID)
		if valid {
			//order.CouponDiscountPercentage = coupon.Amount
			discount := (order.GrandTotal * coupon.Amount) / 100.0
			order.CouponDiscount = math.Floor(discount*100.0)/100.0

			order.GrandTotal = order.GrandTotal - order.CouponDiscount
			order.GrandTotal = math.Round(order.GrandTotal*100.0)/100.0

			order.CouponOrder.CouponID = coupon.ID
			order.CouponOrder.UserID = user.ID
		}
	}
	return order
}

func CheckoutDBProcessing(c *gin.Context, order Mod.Order, payment string, fromApi bool) (Mod.Order,bool) {
	db := Cfg.DBConnect()
	db = db.Begin()
	order, success := R.AddOrder(db, order)
	if !success {
		db.Rollback()
		defer db.Close()
		log.Println("Checkout failed. Check here")
		return order, false
	}
	log.Println("Checkout OrderID: ", order.ID)

	if order.CouponOrder.CouponID != 0 {
		order.CouponOrder.OrderID = order.ID
		if !R.AddCouponOrder(db, order.CouponOrder) {
			db.Rollback()
			defer db.Close()
			log.Println("Checkout failed. Check here")
			return order, false
		}
	}

	for _, oddt := range order.OrderDetails {
		oddt.OrderID = order.ID
		if !R.AddOrderDetail(db, oddt) {
			defer db.Close()
			log.Println("Checkout failed. Check here")
			return order, false
		}
	}

	order.OrderShippingAddress.OrderID = order.ID
	if !R.AddOrderShippingAddress(db, order.OrderShippingAddress) {

		db.Rollback()
		defer db.Close()
		log.Println("Checkout failed. Check here")
		return order, false
	}
	order.OrderShippingAddress.Country = G.Country[order.OrderShippingAddress.Country]

	if payment == "2" {
		order, success = PaymentMethodProcessing(c, db, order)
		if !success {
			db.Rollback()
			defer db.Close()
			log.Println("Checkout failed. Check here")
			return order, false
		}
	}
	db.Commit()
	defer db.Close()

	var user Mod.User
	user.ID = order.UserID

	if !fromApi {
		proID, _ := strconv.Atoi(c.PostForm("product_id"))
		if c.Request.Header["Referer"][0] != G.AppEnv.Url+"/buy-now" && proID == 0 {
			for {
				if DestroyCart(c) {
					break
				}
			}
		}
	} else {
		for {
			if DestroyCart(c) {
				break
			}
		}
	}

	return order, true
}

func PaymentMethodProcessing(c *gin.Context, db *gorm.DB, order Mod.Order) (Mod.Order, bool) {
	payMethodID,err := strconv.Atoi(c.PostForm("pay_method_id"))
	if err != nil {
		log.Println(err.Error())
		G.Msg.Fail = "Invalid payment method."
		db.Rollback()
		return order, false
	} else if payMethodID == 0 {
		G.Msg.Fail = "Invalid payment method."
		db.Rollback()
		return order, false
	}

	var payM Mod.PayMethod
	payM.ID = uint(payMethodID)
	if !R.PayMethodExists(payM, "status=?", 1) {
		G.Msg.Fail = "Payment method not found."
		db.Rollback()
		return order, false
	}

	order.PayMethodID = uint(payMethodID)

	if order.PayMethodID == 1 {
		order, success := StripePaymentProcessing(c, db, order)
		return order, success

	} else if order.PayMethodID == 2 {
		order, success := MolliePaymentProcessing(c, db, order)
		return order, success
	}
	log.Println("Checkout failed. Check here")
	return order, false
}

func StripePaymentProcessing(c *gin.Context, db *gorm.DB, order Mod.Order) (Mod.Order,bool) {
	stripeToken := c.PostForm("stripeToken")
	stripeEmail := c.PostForm("stripeEmail")

	ch, success := StripeCreateCharge(c, int(order.GrandTotal*100.0), order.ID+10000)
	log.Println(*ch)

	if !success || ch.Status != "succeeded" {
		log.Println("Checkout failed. Check here")
		return order, false
	}
	log.Println("Stripe Charge ID:", ch.ID)
	log.Println("Stripe Charge Amount:", ch.Amount, "----- Order Total Amount * 100:", int64(order.GrandTotal*100.0))

	if success && ch.Amount < int64(order.GrandTotal*100.0) {
		G.Msg.Fail = "You Did Something Fraudulent. You Didn't Pay The Right Amount. Your Order Is Cancelled."
		order.PaymentStatus = 2
		order.PaidAmount = float64(ch.Amount) / 100.0
		order.OrderStatus = 2

	} else if success && ch.Paid && ch.Status == "succeeded" {
		order.PaymentStatus = 1
		order.PaidAmount = float64(ch.Amount) / 100
		order.UniqueCode = ch.ID
		if order.DeliveryType == 1 {
			order.OrderStatus = 1
		}

		order.StripeTransaction.OrderID = order.ID
		order.StripeTransaction.StripeToken = stripeToken
		order.StripeTransaction.ChargeID = ch.ID
		order.StripeTransaction.CustomerEmail = stripeEmail
		order.StripeTransaction.Currency = string(ch.Currency)
		order.StripeTransaction.Amount = int(ch.Amount)
		if ch.Paid == true {
			order.StripeTransaction.PaidStatus = 1
		}
		order.StripeTransaction.Status = ch.Status

		if !R.AddStripeTransaction(db, order.StripeTransaction) {
			log.Println("stripe_transactions table insertion failed")
		}

		if !R.SaveUpdatedOrder(db, order) {
			log.Println("orders table update failed")
		}
	}
	return order, true
}


func MolliePaymentProcessing(c *gin.Context, db *gorm.DB, order Mod.Order) (Mod.Order,bool) {

	response, success := MolliePayment(order.GrandTotal, order.User.Email, order.ID, c)

	if !success {
		log.Println("Checkout failed. Check here")
		return order, false
	}
	log.Println("Mollie Payment ID:", response.ID)
	log.Println("Mollie Payment Amount:", response.Amount, "----- Order Total Amount:", order.GrandTotal)

	if success {
		order.UniqueCode = response.ID

		order.MollieOrder.OrderID = order.ID
		order.MollieOrder.PaymentID = response.ID
		order.MollieOrder.UserID = order.UserID
		order.MollieOrder.PaymentStatus = string(response.Status)
		order.MollieOrder.PaymentUrl = response.Links["paymentUrl"]

		if !R.AddMollieOrder(db, order.MollieOrder) {
			log.Println("mollie_orders table insertion failed")
		}

		if !R.SaveUpdatedOrder(db, order) {
			log.Println("orders table update failed")
		}

	}
	return order, true
}






/*func FreeBidsCheckout(c *gin.Context, order Mod.Order, coupon Mod.Coupon) bool {
	db := Cfg.DBConnect()
	db = db.Begin()
	err := db.Create(&order).Error
	if err != nil {
		db.Rollback()
		defer db.Close()
		return false
	}
	var finalCart G.FinalCart
	finalCart = FreeCartProcess(coupon)

	var orderDetail Mod.OrderDetail
	var orderDetails []Mod.OrderDetail
	for _, cart := range finalCart.Carts {
		orderDetail.OrderID = order.ID
		orderDetail.ProductID = cart.Product.ID
		orderDetail.ProductTitle = cart.Product.Title
		orderDetail.ProductPrice = cart.Product.Price
		orderDetail.ProductGrmAmount = cart.Product.GrmAmount
		orderDetail.ProductImgUrl = cart.Product.ImgUrl
		orderDetail.Quantity = cart.Quantity
		orderDetail.Total = cart.Total
		orderDetail.TotalDiscount = cart.TotalDiscount
		orderDetail.TotalWithDiscount = cart.TotalWithDiscount
		orderDetails = append(orderDetails, orderDetail)
	}

	for i,_ := range orderDetails {
		err = db.Create(&orderDetails[i]).Error
		if err != nil {
			db.Rollback()
			defer db.Close()
			return false
		}
	}
	order.OrderDetails = orderDetails

	var coupOrd Mod.CouponOrder
	coupOrd.OrderID = order.ID
	coupOrd.UserID = order.UserID
	coupOrd.CouponID = coupon.ID
	err = db.Create(&coupOrd).Error
	if err != nil {
		db.Rollback()
		defer db.Close()
		return false
	}
	order.CouponOrder = coupOrd
	order.CouponOrder.Coupon = coupon
	order.IsFree = true

	//success := SendBidToWallet(order, c)
	//log.Println(success)
	//if !success {
	//	db.Rollback()
	//	defer db.Close()
	//	return false
	//}

	G.Msg.Success = "Congratulations, Free Bids Are Sent To Your Wallet Successfully."
	filename, success := GenerateInvoice(order)
	log.Println("Invoice Generation Success:", success)
	success = SendOrderEmail(order, c, filename)
	log.Println("Order Completion Email Sending Success:", success)

	db.Commit()
	defer db.Close()
	return true
}*/
