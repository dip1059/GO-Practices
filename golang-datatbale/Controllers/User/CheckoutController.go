package User

import (
	"github.com/gin-gonic/gin"
	Cfg "gold-store/Config"
	G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	S "gold-store/Services"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CheckoutGet(c *gin.Context) {
	user, guest := M.IsGuest(c, G.FStore)
	user, authUser := M.IsAuthUser(c, G.FStore)
	sessionX, _ := G.Store.Get(c.Request, "cart")
	total := len(sessionX.Values)
	//total, _ := strconv.Atoi(c.Param("total"))
	if authUser && total > 0 {
		if user.DefaultAddress.ID == 0 {
			user.DefaultAddress.FirstName = user.FirstName
			user.DefaultAddress.LastName = user.LastName
			user.DefaultAddress.Phone = user.Phone
			user.DefaultAddress.Country = user.CountryCode
			user.DefaultAddress.City = user.City
			user.DefaultAddress.Address = user.Address
			user.DefaultAddress.ZipCode = user.ZipCode
		}
		user.DefaultAddress.Country = G.Country[user.DefaultAddress.Country]

		var countries []Country
		var country Country
		for key, value := range G.Country {
			country.Key = key
			country.Value = value
			countries = append(countries, country)
		}

		lenWish := R.CountWishlist(user)
		session, _ := G.Store.Get(c.Request, "cart")
		lenCart := len(session.Values)
		finalCart := S.ProcessCart(c)
		DataAmount := int64(finalCart.GrandTotal * 100.0)
		stripePK := G.Adm["Stripe_Publishable_Key"].Value.String
		//brainTkzKey := G.Adm["Brain_Tree_Tokenization_Key"].Value.String
		bankRefCode := "REF"+H.RandomString(11)

		var payMethods []Mod.PayMethod
		payMethods = R.AllMethod(payMethods, "status=1 and id<>?", 7)

		var banks []Mod.Bank
		banks = R.AllBank(banks, "status=1")

		var wbsts []Mod.WebsiteSetting
		wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

		var menus []Mod.Menu
		menus = R.Menus(menus, "status=?", 1)

		G.LangVal = H.GetCookie("secret", nil,"lang", c)
		if G.LangVal == "" {
			G.LangVal = "en"
			H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
		}
		var lang G.Lang
		lang.LangValue = G.LangVal

		redsysOrdId := H.RandomNumber(10)

		c.HTML(http.StatusOK, "checkout.html", map[string]interface{}{
			"AppEnv": G.AppEnv, "User": user, "Banner": "Checkout", "Carts": lenCart,"Lang":lang,"Banks": banks,
			"Msg": G.Msg, "Adm": G.Adm, "Title": "Checkout", "DataAmount": DataAmount, "FinalCart": finalCart,
			/*"SecretKey":sessionId,*/ /*"RedsysUrl":G.Redsys.Url,*/ "StripePK":stripePK,"ProductID":0,"FreeBids":false,
			"LenWish": lenWish, "PayMethods": payMethods,  "UniqueCode":bankRefCode, "RedsysOrdId":redsysOrdId,
			"Wbsts": wbsts, "Menus":menus,"Active": "Cart", "Countries": countries, "Quantity":0/*"WalletApp":G.WalletApp,*/ })
		G.Msg.Success = ""
		G.Msg.Fail = ""
	} else if guest && total > 0 {
		if user.ID > 0 && user.ActiveStatus == 0 {
			G.Msg.Fail = "Please Activate Your Account."
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		}else {
			c.Redirect(http.StatusFound, "/login?val=checkout")
		}
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}




func DownloadCheckoutPDF(c *gin.Context) {
	user, _ :=M.IsAuthUser(c,G.FStore)

	bankID, _ := strconv.Atoi(c.Param("bankId"))
	var bank Mod.Bank
	bank.ID = uint(bankID)
	bank = R.Bank(bank)

	finalCart := S.ProcessCart(c)

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts)

	ref := c.Param("ref")

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	Cfg.LoadAdminSettings()

	pdf := S.NewRequestPdf("")

	err := pdf.ParseTemplate("Views/PdfTemp/checkout-bank-details-pdf.html", map[string]interface{}{
		"User":user, "Wbsts":wbsts, "FinalCart":finalCart, "Bank":bank,
		"UniqueCode":ref, "Adm":G.Adm, "Lang":lang, "AppEnv":G.AppEnv})
	if err != nil {
		log.Println(err.Error())
		return
	}

	t := time.Now().Unix()
	fileName := strconv.FormatInt(int64(t), 10)
	success := pdf.GeneratePDF("Storage/Temp/"+fileName+".pdf")
	if success {
		log.Println("PDF Generated Successfully.")
		c.Redirect(http.StatusFound, "/assets/Storage/Temp/"+fileName+".pdf")
		/*err := os.Remove("Storage/Temp/"+fileName+".pdf")
		if err != nil {
			log.Println(err.Error())
		}*/
	}

}


func RedsysCheckoutProcess(c *gin.Context) {
	S.RedsysProcess(c)
}

func RedsysCheckoutSuccess(c *gin.Context) {
	//order, user := S.RedsysSuccess(c)

	/*filename, success := GenerateInvoice(order)
	log.Println("Invoice Generation Success:", success)*/
	//success := SendOrderEmail(order, user, c, "")
	//log.Println("Order Placing Email Sending Success:", success)
}

func RedsysCheckoutFail(c *gin.Context) {
	S.RedsysFail(c)
}


func LusopayCheckoutProcess(c *gin.Context) {
	S.LusopayProcess(c)
}


/*func CheckoutPost(c *gin.Context) {
	user, success := M.IsAuthUser(c,G.FStore)
	if !success {
		return
	}

	var coupon Mod.Coupon
	var couponValid bool
	code := c.PostForm("coupon_code")
	if code != "" {
		coupon = R.Coupon(coupon, "code = BINARY ? and status=? and type=?",code, 1, 1)
		couponValid = S.IsCouponValid(coupon, user.ID)
		if !couponValid {
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}
	}

	log.Println("Checkout started for User:",user.Email, user.ID)
	payMethodID, _ := strconv.Atoi(c.PostForm("pay_method_id"))
	log.Println("PayMethodID:", payMethodID)
	var payMethod Mod.PayMethod
	payMethod.ID = uint(payMethodID)

	if payMethodID == 7 || !R.PayMethodExists(payMethod, "status=?", 1) {
		G.Msg.Fail = "Payment method not found."
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}

	proID, _ := strconv.Atoi(c.PostForm("product_id"))
	log.Println("productID:",proID)
	cashP, _ := strconv.ParseFloat(c.PostForm("cash_point"), 64)

	var stripeToken string
	var stripeEmail string
	var ch *stripe.Charge

	var trx *braintree.Transaction

	var order Mod.Order
	var finalCart G.FinalCart

	if proID < 0 {
		return
	}

	if proID == 0 {
		finalCart = S.ProcessCart(c)
	} else {
		finalCart = S.BuyNowProcess(uint(proID), c)
	}

	if finalCart.GrandTotal == 0.0 {
		log.Println( "Order Total Amount:", int64(finalCart.GrandTotal * 100.0))
		return
	}
	flag := 0

	if payMethodID == 2 || !R.PayMethodExists(payMethod, "status=?", 1) {
		S.MolliePayment(finalCart.GrandTotal,user.Email,1, c)
		return
	}

	if couponValid {
		discount := (finalCart.GrandTotal * coupon.Amount) / 100.0
		order.CouponDiscount = math.Floor(discount*100.0)/100.0
		order.IsDiscountCoupon = true

		finalCart.GrandTotal = finalCart.GrandTotal - order.CouponDiscount
		finalCart.GrandTotal = math.Ceil(finalCart.GrandTotal*100.0)/100.0
	}

	if payMethodID == 6 {
		if !S.CheckCashPoint(finalCart.GrandTotal, cashP, user) {
			G.Msg.Fail = "You are applying fraudulency. This is a warning."
			c.Redirect(http.StatusFound, "/")
			return
		}
	}

	if payMethodID == 1 {
		stripeToken = c.PostForm("stripeToken")
		stripeEmail = c.PostForm("stripeEmail")
		log.Println("stripeToken:", stripeToken," --- stripeEmail:",stripeEmail)

		//ch, success = S.StripeCreateCharge(c)

		log.Println("Stripe Charge ID:", ch.ID)
		log.Println("Stripe Charge Amount:",ch.Amount, "----- Order Total Amount:", int64(finalCart.GrandTotal * 100.0))

		if success && ch.Amount < int64(finalCart.GrandTotal * 100.0) {
			G.Msg.Fail = "You Did Something Fraudulent. You Didn't Pay The Right Amount. Your Order Is Cancelled."
			order.PaymentStatus = 2
			order.PaidAmount = float64(ch.Amount) / 100.0
			order.OrderStatus = 2
		} else if success && ch.Paid && ch.Status == "succeeded" {
			flag = 1
			order.PaymentStatus = 1
			order.PaidAmount = float64(ch.Amount) / 100
			order.UniqueCode = ch.ID
		}
	}

	if payMethodID == 2 || payMethodID == 3 {
		trx, success = S.CreateTransaction(c)
		if success && trx.Amount.Unscaled != int64(math.Ceil(finalCart.SubTotal*100.0)) {
			G.Msg.Fail = "You Did Something Fraudulent. You Didn't Pay The Right Amount. Your Order is Cancelled."
			order.PaymentStatus = 2
			order.PaidAmount = float64(trx.Amount.Unscaled / 100)
			order.OrderStatus = 2
		} else if success {
			flag = 1
			order.PaymentStatus = 1
			order.PaidAmount = float64(trx.Amount.Unscaled / 100)
		}
		//if success {
		//	if string(trx.PaymentInstrumentType) == "credit_card" {
		//		order.BrainTreeMethodID.Int64 = 1
		//		order.BrainTreeMethodID.Valid = true
		//	} else if string(trx.PaymentInstrumentType) == "paypal_account" {
		//		order.BrainTreeMethodID.Int64 = 2
		//		order.BrainTreeMethodID.Valid = true
		//	}
		//}
	}

	var lusopay Mod.Lusopay

	if success || payMethodID == 1 || payMethodID == 4 || payMethodID == 6 {
		order.UserID = user.ID
		order.PayMethodID = uint(payMethodID)
		order.SubTotal = finalCart.SubTotal
		order.GrandTotal = finalCart.GrandTotal
		order.TotalGrmAmount = finalCart.TotalGrmAmount
		order.Fees = finalCart.Fees
		order.VAT = finalCart.VAT
		order.TotalFees = finalCart.TotalFees
		order.TotalVAT = finalCart.TotalVAT
		order.TotalDiscount = finalCart.TotalDiscount
		if payMethodID == 1 {
			flag = 1
			bankId, _ := strconv.Atoi(c.PostForm("bank_id"))
			log.Println("bankID:", bankId)
			var bank Mod.Bank
			bank.ID = uint(bankId)
			if !R.BankExists(bank, "status=?", 1) {
				G.Msg.Fail = "Bank not found."
				c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
				return
			}

			order.UniqueCode = "REF"+H.RandomString(11)
			log.Println("bank_unique_code:", order.UniqueCode)
			order.BankID.Int64 = int64(bankId)
			order.BankID.Valid = true
		}

		if payMethodID == 4 {
			flag = 1
			order.UniqueCode = c.PostForm("reference")
		}

		if payMethodID == 6 {
			flag = 1
			order.PaidAmount = cashP
			order.PaymentStatus = 1
		}

		//order, success = R.AddOrder(order)
		if !success {
			if G.Msg.Fail == "" {
				G.Msg.Fail = "Some Error Occurred. Please Try Again Later."
			}
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}

		if payMethodID == 4 {
			flag = 1
			lusopay.OrderID = order.ID
			lusopay.Entity = c.PostForm("entity")
			lusopay.Reference = order.UniqueCode
			lusopay.Amount, _ = strconv.ParseFloat(c.PostForm("amount"), 64)

			log.Println("lusopay data insertion:", R.AddLusopay(lusopay))
		}

		var orderDetail Mod.OrderDetail
		var orderDetails []Mod.OrderDetail
		for _, cart := range finalCart.Carts {
			orderDetail.OrderID = order.ID
			orderDetail.ProductID = cart.Product.ID
			orderDetail.KaratID = cart.Product.KaratID
			orderDetail.ProductKaratAmount = cart.Product.Karat.Amount
			orderDetail.ProductTitle = cart.Product.Title
			orderDetail.ProductType = cart.Product.Type
			orderDetail.ProductPrice = cart.Product.Price
			orderDetail.ProductGrmAmount = cart.Product.GrmAmount
			orderDetail.ProductImgUrl = cart.Product.ImgUrl
			orderDetail.Quantity = cart.Quantity
			orderDetail.Total = cart.Total
			orderDetail.TotalDiscount = cart.TotalDiscount
			orderDetail.TotalWithDiscount = cart.TotalWithDiscount
			orderDetails = append(orderDetails, orderDetail)
		}
		orderDetails, success = R.AddOrderDetails(orderDetails)
		if !success {
			if G.Msg.Fail == "" {
				G.Msg.Fail = "Some Error Occurred. Please Try Again Later."
			}
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			return
		}

		order.OrderDetails = orderDetails

		var coupOrd Mod.CouponOrder
		if couponValid {
			coupOrd. OrderID = order.ID
			coupOrd.UserID = user.ID
			coupOrd.CouponID = coupon.ID
			//coupOrd, success = R.AddCouponOrder(coupOrd)

			if !success {
				if G.Msg.Fail == "" {
					log.Println("Coupon Order insertion failed")
					G.Msg.Fail = "Some Error Occurred. Please Try Again Later."
				}
				c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
				return
			}
			order.CouponOrder = coupOrd
		}

		if payMethodID == 5 {
			var stTx Mod.StripeTransaction
			stTx.OrderID = order.ID
			stTx.StripeToken = stripeToken
			stTx.ChargeID = ch.ID
			stTx.CustomerEmail = stripeEmail
			stTx.Currency = string(ch.Currency)
			stTx.Amount = int(ch.Amount)
			if ch.Paid == true {
				stTx.PaidStatus = 1
			}
			stTx.Status = ch.Status
			//if success = R.AddStripeTransaction(stTx); !success {
			//	if G.Msg.Fail == "" {
			//		G.Msg.Fail = "Some Error Occurred. Please Try Again Later."
			//	}
			//	c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			//	return
			//}
		}

		if payMethodID == 2 || payMethodID == 3 {
			var btTx Mod.BrainTreeTransaction
			btTx.OrderID = order.ID
			btTx.TxID = trx.Id
			btTx.CustomerEmail = user.Email
			btTx.Currency = trx.CurrencyISOCode
			btTx.Amount = int(trx.Amount.Unscaled)
			btTx.PaymentInstrumentType = string(trx.PaymentInstrumentType)
			btTx.Status = string(trx.Status)

			//if success = R.AddBTTx(btTx); !success {
			//	if G.Msg.Fail == "" {
			//		G.Msg.Fail = "Some Error Occurred. Please Try Again Later."
			//	}
			//	c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
			//	return
			//}
		}
		if proID <= 0 {
			S.DestroyCart(c)
		}
		if flag == 1 {
			if payMethodID == 6 {
				for {
					if R.DeductCashPoint(user, cashP) {
						log.Println("Cash Point", cashP, "deduction successful for user", user.ID)
						break
					} else {
						log.Println("Cash Point", cashP, "deduction failed for user", user.ID)
					}
				}
			}

			if payMethodID == 1 {
				G.Msg.Success = `Thank you for your purchase. Kindly download the document with the bank details on your order below.`
			} else {
				G.Msg.Success = `Thank you for your purchase. Your order has been placed successfully.`
			}
			//filename, success := GenerateInvoice(order)
			//log.Println("Invoice Generation Success:", success)

			success = SendOrderEmail(order, user, c, "")
			log.Println("Order Placing Email Sending Success:", success)
		}
		c.Redirect(http.StatusFound, "/account/orders")
	} else {
		if G.Msg.Fail == "" {
			G.Msg.Fail = "Some Error Occurred. Please Try Again Later."
		}
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
	}

}*/




//new checkout
func CheckoutPost(c *gin.Context) {
	user := M.GetAuthUser(c, G.FStore)

	if user.RoleID != 2 {
		return
	}

	success, shipAdd := S.ShippingAddressValidation(c)
	if !success {
		if c.Request.Header["Referer"][0] == G.AppEnv.Url+"/buy-now" {
			c.Redirect(http.StatusFound, "/")
		} else {
			c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		}
		return
	}

	log.Println("Checkout started for User:",user.Email, user.ID)

	var order Mod.Order
	order.OrderShippingAddress = shipAdd
	order, success = S.CheckoutProcessing(c, user, order, false)
	if success {
		if order.PayMethodID == 2 {
			c.Redirect(http.StatusFound, order.MollieOrder.PaymentUrl)
		} else {
			if order.DeliveryType == 2 {
				G.Msg.Success = "Thank You, Your Order Has Been Placed Successfully."
			} else if order.DeliveryType == 1 {
				G.Msg.Success = "Thank You, Your Order Has Been Placed And Completed Successfully."
			}
			c.Redirect(http.StatusFound, "/account/orders")
			//S.GenerateCertificate(order)
			go func() {
				S.GenerateCertificate(order, c)
				//time.Sleep(time.Second * 20)
				//log.Println("go routine worked")
			}()
		}
		return
	} else {
		if G.Msg.Fail == "" {
			G.Msg.Fail = "Some error occurred. Please reload and try again later."
			if c.Request.Header["Referer"][0] == G.AppEnv.Url+"/buy-now" {
				c.Redirect(http.StatusFound, "/")
				return
			}
		}
		c.Redirect(http.StatusFound, c.Request.Header["Referer"][0])
		return
	}
}




func SendOrderEmail(order Mod.Order, user Mod.User, c *gin.Context, filename string) bool{
	//order = R.Order(order)
	orderID := 10000+order.ID

	var wbsts []Mod.WebsiteSetting
	wbsts = R.WebsiteSettings(wbsts, "status=?", 1)

	G.LangVal = H.GetCookie("secret", nil,"lang", c)
	if G.LangVal == "" {
		G.LangVal = "en"
		H.SetCookie("secret", nil, G.LangVal, "lang", 60*60*24*365, c)
	}
	var lang G.Lang
	lang.LangValue = G.LangVal

	htmlString, err := H.ParseTemplate("Views/Email/order-email.html", map[string]interface{}{
		"User":user, "Wbsts":wbsts, "Order":order, "OrderID":orderID, "Adm":G.Adm, "Lang":lang, "AppEnv":G.AppEnv})

	if err != nil {
		log.Println(err.Error())
		return false
	}

	To := []string{user.Email}
	Subject := "Order Placing Details"
	HtmlString := htmlString
	if !S.SendEmail(To, Subject, HtmlString, filename) {
		return false
	}
	return true
}


func GenerateInvoice(order Mod.Order) (string, bool){
	pdf := S.NewRequestPdf("")
	/*err := pdf.ParseTemplate("Views/PdfTemp/checkout-details-pdf.html", map[string]interface{}{
		"User":user, "Wbsts":wbsts, "FinalCart":finalCart, "Bank":bank,
		"UniqueCode":ref, "Adm":G.Adm, "Lang":lang, "AppEnv":G.AppEnv})*/

	var invoice Mod.Invoice
	invoice.OrderID = order.ID
	invoice = R.AddInvoice(invoice)

	err := pdf.ParseTemplate("Views/PdfTemp/invoice.html", invoice.ID)
	if err != nil {
		log.Println(err.Error())
		return "", false
	}

	t := time.Now().Unix()
	fileName := strconv.FormatInt(int64(t), 10)
	success := pdf.GeneratePDF("Storage/Temp/"+fileName+".pdf")
	if success {
		log.Println("PDF Generated Successfully.")
		//c.Redirect(http.StatusFound, "/assets/Storage/Temp/"+fileName+".pdf")
	}
	return "Storage/Temp/"+fileName+".pdf", true
}


func CheckCashPoint(c *gin.Context) {
	user, success := M.IsAuthUser(c,G.FStore)
	if !success {
		return
	}
	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)
	cashPointData := R.CheckCashPoint(user)

	cashP := amount / cashPointData.Rate

	if cashP > cashPointData.CashPoint {
		c.JSON(http.StatusOK, gin.H {
			"success":false,
			"needed_point":cashP,
			"available_point":cashPointData.CashPoint,
		})
	} else {
		c.JSON(http.StatusOK, gin.H {
			"success":true,
			"needed_point":cashP,
			"available_point":cashPointData.CashPoint,
		})
	}
}


func GetDiscountCoupon(c *gin.Context) {
	user, success := M.IsAuthUser(c,G.FStore)
	if !success {
		return
	}
	var message string
	code := c.PostForm("coupon_code")
	var coupon Mod.Coupon
	coupon = R.Coupon(coupon, "code = BINARY ? and status=? and type=?",code, 1, 1)
	if coupon.ID > 0 {
		usage := R.CouponUsageCount(user.ID, coupon.ID)
		if usage <= 0 {
			now, _ := time.Parse("2006-01-02 15:04", time.Now().Format("2006-01-02 15:04"))
			if diffStart := now.Sub(coupon.StartDate.UTC()); diffStart.Seconds() >= 0 {
				if diffEnd := now.Sub(coupon.EndDate.UTC()); diffEnd.Seconds() >= 0 {
					coupon.Status = 0
					R.SaveCoupon(coupon)
					success = false
					message = "The code is not valid or expired."
					coupon = Mod.Coupon{}
				} else {
					success = true
				}
			} else {
				success = false
				message = "The code is not valid or expired."
				coupon = Mod.Coupon{}
			}
		} else {
			success = false
			message = "You have used the code already."
			coupon = Mod.Coupon{}
		}
	} else {
		success = false
		message = "The code is not valid or expired."
		coupon = Mod.Coupon{}
	}

	c.JSON(http.StatusOK, gin.H {
		"success":success,
		"coupon": coupon,
		"message":message,
	})

}
