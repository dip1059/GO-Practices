package Services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	Cfg "gold-store/Config"
	G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"log"
	"math"
	"strconv"
)

func AddCart(c *gin.Context, id string, quantity string) bool {
	session, _ := G.Store.Get(c.Request, "cart")
	//lenCart := len(session.Values)
	randomKey := H.RandomString(40)
	var cart = []string{id, quantity}
	flag := 0
	process := c.Query("process")
	
	for key,_ := range session.Values {
		val := session.Values[key].([]string)
		if cart[0] == val[0] {
			val1, _ := strconv.ParseFloat(val[1], 64)
			cart1, _ := strconv.ParseFloat(cart[1], 64)

			/*if val[1] > 20 {
				G.Msg.Fail = "You Have Reached The Maximum Quantity Limit, You Can Order Maximum 20 Of A Product."
				val[1] = 20
			}*/

			if process == "update" {
				val[1] = cart[1]
			} else {
				val[1] = fmt.Sprintf("%.1f", val1+cart1)
			}
			session.Values[key] = val
			flag = 1
			break
		}
	}
	if flag == 0 {
		session.Values[randomKey] = cart
	}

	err := session.Save(c.Request, c.Writer)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}


func ProcessCart(c *gin.Context) G.FinalCart {
	Cfg.LoadAdminSettings()

	user := M.GetAuthUser(c, G.FStore)

	session, _ := G.Store.Get(c.Request, "cart")
	var carts []G.Cart
	var err error

	var Total float64
	var TotalDiscount float64
	var TotalGrmAmount float64

	var FinalSubTotal float64
	var FinalTotalDiscount float64

	var finalCart G.FinalCart

	for key, val := range session.Values {

		/*tpe := fmt.Sprintf("%T", val)
		//log.Println(tpe)
		if tpe == "[]int" {
			DestroyCart(c)
			return finalCart
		}*/

		temp := val.([]string)
		var product Mod.Product
		var id int
		id, err = strconv.Atoi(temp[0])
		if err != nil {
			log.Println(err.Error())
		}
		product.ID = uint(id)
		product = R.Product(product, user.ID, "status = ?", 1)
		if product.Status == 0 {
			delete(session.Values, key)
			continue
		}

		var cart G.Cart
		cart.Key = key.(string)
		cart.Product = product
		cart.Quantity, err = strconv.ParseFloat(temp[1], 64)
		if err != nil {
			log.Println(err.Error())
		}

		if product.Type == 1 {
			TotalDiscount = float64(cart.Quantity) * (product.Price * product.Discount / 100.0)
			Total = float64(cart.Quantity) * product.Price

		} else if product.Type == 2 {
			TotalDiscount = float64(cart.Quantity) * ((product.Price / product.GrmAmount) * product.Discount / 100.0)

			Total = cart.Quantity * (product.Price / product.GrmAmount)
			Total = math.Ceil(Total*100.0)/100.0
			cart.Product.Price = Total
			cart.Product.GrmAmount = cart.Quantity
		}
		TotalDiscount = math.Round(TotalDiscount*100.0)/100.0
		cart.TotalDiscount = TotalDiscount
		FinalTotalDiscount += TotalDiscount

		cart.Total = Total
		TotalGrmAmount += cart.Product.GrmAmount
		cart.TotalWithDiscount =  Total - TotalDiscount
		cart.TotalWithDiscount = math.Round(cart.TotalWithDiscount*100.0)/100.0
		FinalSubTotal += cart.TotalWithDiscount
		carts = append(carts, cart)
	}

	finalCart.Carts = carts
	finalCart.TotalGrmAmount = math.Floor(TotalGrmAmount*1000.0)/1000.0
	finalCart.SubTotal = math.Ceil(FinalSubTotal*100.0)/100.0
	finalCart.TotalDiscount = FinalTotalDiscount
	//finalCart.ShippingCost = Cfg.ShippingCost
	//finalCart.VAT, err = strconv.ParseFloat(G.Adm["VAT"].Value.String, 64)
	//if err != nil {
	//	log.Println(err.Error())
	//}
	finalCart.Fees1Percent, err = strconv.ParseFloat(G.Adm["Order_Fees_1_In_Percent"].Value.String, 64)
	if err != nil {
		log.Println(err.Error())
	}
	finalCart.Fees1Fixed, err = strconv.ParseFloat(G.Adm["Order_Fees_1_Fixed_In_Euro"].Value.String, 64)
	if err != nil {
		log.Println(err.Error())
	}
	finalCart.Fees2Percent, err = strconv.ParseFloat(G.Adm["Order_Fees_2_In_Percent"].Value.String, 64)
	if err != nil {
		log.Println(err.Error())
	}
	finalCart.Fees2Fixed, err = strconv.ParseFloat(G.Adm["Order_Fees_2_Fixed_In_Euro"].Value.String, 64)
	if err != nil {
		log.Println(err.Error())
	}
	finalCart.Fees3Percent, err = strconv.ParseFloat(G.Adm["Order_Fees_3_In_Percent"].Value.String, 64)
	if err != nil {
		log.Println(err.Error())
	}
	finalCart.Fees3Fixed, err = strconv.ParseFloat(G.Adm["Order_Fees_3_Fixed_In_Euro"].Value.String, 64)
	if err != nil {
		log.Println(err.Error())
	}

	totalFeesDiscount := finalCart.Fees1Percent + finalCart.Fees2Percent + finalCart.Fees3Percent
	totalFeesFixed := finalCart.Fees1Fixed + finalCart.Fees2Fixed + finalCart.Fees3Fixed

	finalCart.TotalFees = math.Ceil(((finalCart.SubTotal * totalFeesDiscount) / 100) * 100) / 100
	finalCart.TotalFees += totalFeesFixed
	finalCart.TotalFees = math.Ceil(finalCart.TotalFees * 100) / 100

	finalCart.GrandTotal = finalCart.SubTotal + finalCart.TotalFees
	finalCart.GrandTotal = math.Round(finalCart.GrandTotal*100.0)/100.0

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Println(err.Error())
	}

	return finalCart
}

func DestroyCart(c *gin.Context) bool {
	session, _ := G.Store.Get(c.Request, "cart")
	session.Options.MaxAge = -1
	err := session.Save(c.Request, c.Writer)

	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true

}

func DeleteFromCart(c *gin.Context, key string) bool {
	session, _ := G.Store.Get(c.Request, "cart")
	delete(session.Values, key)
	err := session.Save(c.Request, c.Writer)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}


func BuyNowProcess(id uint, c *gin.Context) G.FinalCart {

	user := M.GetAuthUser(c, G.FStore)

	Cfg.LoadAdminSettings()
	var carts []G.Cart
	//var err error

	var Total float64
	var TotalDiscount float64
	var TotalGrmAmount float64

	var FinalSubTotal float64
	var FinalTotalDiscount float64

	var finalCart G.FinalCart

	var product Mod.Product
	var cart G.Cart
	product.ID = id
	product = R.Product(product, user.ID)
	cart.Product = product
	cart.Quantity = 1

	TotalDiscount = float64(cart.Quantity) * (product.Price * product.Discount / 100.0)
	cart.TotalDiscount = TotalDiscount
	FinalTotalDiscount += TotalDiscount
	if product.Type == 1 {
		Total = float64(cart.Quantity) * product.Price

	} else if product.Type == 2 {
		Total = float64(cart.Quantity) * (product.Price / product.GrmAmount)
		cart.Product.Price = Total
		cart.Product.GrmAmount = float64(cart.Quantity)
	}
	cart.Total = Total
	FinalSubTotal += Total
	TotalGrmAmount += cart.Product.GrmAmount
	cart.TotalWithDiscount = Total - TotalDiscount
	carts = append(carts, cart)


	finalCart.Carts = carts
	finalCart.TotalGrmAmount = TotalGrmAmount
	finalCart.SubTotal = math.Ceil(FinalSubTotal*100.0)/100.0
	finalCart.TotalDiscount = FinalTotalDiscount
	//finalCart.ShippingCost = Cfg.ShippingCost
	//finalCart.VAT, err = strconv.ParseFloat(G.Adm["VAT"].Value.String, 64)
	//if err != nil {
	//log.Println(err.Error())
	//}
	//finalCart.Fees, err = strconv.ParseFloat(G.Adm["Fees"].Value.String, 64)
	//if err != nil {
	//log.Println(err.Error())
	//}
	//finalCart.TotalVAT = math.Ceil(((finalCart.SubTotal * finalCart.VAT) / 100.0) * 100) / 100
	//finalCart.TotalFees = math.Ceil(((finalCart.SubTotal * finalCart.Fees) / 100.0) * 100) / 100
	//finalCart.GrandTotal = finalCart.SubTotal + finalCart.TotalVAT + finalCart.TotalFees - finalCart.TotalDiscount
	finalCart.GrandTotal = math.Ceil(finalCart.GrandTotal*100.0)/100.0
	return finalCart
}


func FreeCartProcess(coupon Mod.Coupon) G.FinalCart {

	Cfg.LoadAdminSettings()
	var carts []G.Cart
	//var err error

	var finalCart G.FinalCart

	var product Mod.Product
	var cart G.Cart

	cart.Quantity = 1

	product = R.Product(product, 0,"type=? and title=BINARY 'Free'",3)
	product.GrmAmount = coupon.Amount

	cart.Product = product
	carts = append(carts, cart)

	finalCart.Carts = carts
	finalCart.TotalGrmAmount = coupon.Amount
	finalCart.SubTotal = 0.0
	//finalCart.VAT, err = strconv.ParseFloat(G.Adm["VAT"].Value.String, 64)
	//if err != nil {
	//	log.Println(err.Error())
	//}
	//finalCart.Fees, err = strconv.ParseFloat(G.Adm["Fees"].Value.String, 64)
	//if err != nil {
	//	log.Println(err.Error())
	//}

	finalCart.GrandTotal = 0
	return finalCart
}
