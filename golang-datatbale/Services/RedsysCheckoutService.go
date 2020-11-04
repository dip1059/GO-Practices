package Services

import (
	/*G "gold-store/Globals"
	H "gold-store/Helpers"
	M "gold-store/Middlewares"*/
	"github.com/gin-gonic/gin"
	/*R "gold-store/Repositories"
	"bytes"
	"encoding/json"
	"github.com/chekun/golaravelsession"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"io/ioutil"
	"log"
	"net/http"*/
)


type RedsysData struct {
	SignatureVersion string		`json:"signatureVersion"`
	Params string	`json:"params"`
	Signature string	`json:"signature"`
	Success bool	`json:"success"`
}


func RedsysProcess(c *gin.Context) {
	/*user, success :=M.IsAuthUser(c,G.FStore)
	if !success {
		c.Redirect(http.StatusFound, "/checkout")
		return
	}
	ord := H.RandomString(12)
	finalCart := ProcessCart(c)

	sc := securecookie.New([]byte(G.WalletApp.Key), nil)
	encOrd, err := sc.Encode("encryptedOrder", ord)
	if err != nil {
		log.Println(err.Error())
	}

	cookie, err := c.Cookie(G.SessionCookie.Name)
	if err != nil {
		log.Println(err.Error())
	}
	sessionId, err := golaravelsession.GetSessionID(cookie, G.WalletApp.Key)
	if err != nil {
		log.Println(err.Error())
	}

	req2, _ := json.Marshal(map[string]interface{}{
		"uniqueId": ord,
		"amount": finalCart.GrandTotal,
		"encOrd": encOrd,
		"secret_key": sessionId,
	})
	log.Println(string(req2))

	url := G.WalletApp.Url+"/api/redsys-secret-generate"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req2))
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(url)

	var redsysTransaction Mod.RedsysTransaction
	redsysTransaction.UserID = user.ID
	redsysTransaction.RedsysOrder = ord
	redsysTransaction.SubTotal = finalCart.SubTotal
	redsysTransaction.TotalBids = finalCart.TotalBids
	redsysTransaction.TotalFees = finalCart.TotalFees
	redsysTransaction.TotalVAT = finalCart.TotalVAT
	redsysTransaction.Total = finalCart.GrandTotal

	if !R.AddRedsysTransaction(redsysTransaction) {
		return
	}

	var redsysOrderDetail Mod.RedsysOrderDetail
	var redsysOrderDetails []Mod.RedsysOrderDetail
	for _, cart := range finalCart.Carts {
		redsysOrderDetail.RedsysOrder = ord
		redsysOrderDetail.ProductID = cart.Product.ID
		redsysOrderDetail.ProductTitle = cart.Product.Title
		redsysOrderDetail.ProductPrice = cart.Product.Price
		redsysOrderDetail.ProductGrmAmount = cart.Product.GrmAmount
		redsysOrderDetail.ProductImgUrl = cart.Product.ImgUrl
		redsysOrderDetail.Quantity = cart.Quantity
		redsysOrderDetail.Total = cart.Total
		redsysOrderDetail.TotalTax = cart.TotalTax
		redsysOrderDetail.TotalWithTax = cart.TotalWithTax
		redsysOrderDetails = append(redsysOrderDetails, redsysOrderDetail)
	}
	redsysOrderDetails, success = R.AddRedsysOrderDetails(redsysOrderDetails)
	if !success {
		return
	}

	Bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(Bytes))

	var d RedsysData
	err = json.Unmarshal(Bytes, &d)
	if err != nil {
		log.Println(err.Error())
		return
	}

	c.JSON(http.StatusOK, d)

	defer resp.Body.Close()*/
}


func RedsysSuccess(c *gin.Context) /*(Mod.Order, Mod.User)*/{
	/*user, authSuccess :=M.IsAuthUser(c,G.FStore)
	var success bool
	var order Mod.Order

	ord := c.Query("ord")
	encOrd := c.Query("encOrd")

	sc := securecookie.New([]byte(G.WalletApp.Key), nil)
	var value string
	err := sc.Decode("encryptedOrder", encOrd, &value)
	if err != nil {
		log.Println(err.Error())
		//c.Redirect(http.StatusFound, "/checkout")
		return order, user
	}
	log.Println(ord, value)

	if ord != value {
		//c.Redirect(http.StatusFound, "/checkout")
		return order, user
	}

	var redsysTransaction Mod.RedsysTransaction

	redsysTransaction = R.RedsysTransaction(redsysTransaction,"redsys_order=?",ord)
	if redsysTransaction.ID == 0 || redsysTransaction.PaidStatus > 0 {
		return order, user
	}


	order.UserID = redsysTransaction.UserID
	order.PayMethodID = 2
	order.SubTotal = redsysTransaction.SubTotal
	order.TotalFees = redsysTransaction.TotalFees
	order.TotalBids = redsysTransaction.TotalBids
	order.TotalVAT = redsysTransaction.TotalVAT
	order.Total = redsysTransaction.Total
	order.PaidAmount = redsysTransaction.Total
	order.PaymentStatus = 1

	if order, success = R.AddOrder(order); !success {
		if !authSuccess || user.ID != redsysTransaction.UserID {
			c.Redirect(http.StatusFound, "/")
			return order, user
		}
		if G.Msg.Fail == "" {
			G.Msg.Fail = "Some Error Occurred. Please Try Again Later."
		}
		c.Redirect(http.StatusFound, "/checkout")
		return order, user
	}


	var orderDetail Mod.OrderDetail
	var orderDetails []Mod.OrderDetail
	for _, ROD := range redsysTransaction.RedsysOrderDetails {
		orderDetail.OrderID = order.ID
		orderDetail.ProductID = ROD.ProductID
		orderDetail.ProductTitle = ROD.ProductTitle
		orderDetail.ProductPrice = ROD.ProductPrice
		orderDetail.ProductGrmAmount = ROD.ProductGrmAmount
		orderDetail.ProductImgUrl = ROD.ProductImgUrl
		orderDetail.Quantity = ROD.Quantity
		orderDetail.Total = ROD.Total
		orderDetail.TotalTax = ROD.TotalTax
		orderDetail.TotalWithTax = ROD.TotalWithTax
		orderDetails = append(orderDetails, orderDetail)
	}
	orderDetails, success = R.AddOrderDetails(orderDetails)
	if !success {
		if !authSuccess || user.ID != redsysTransaction.UserID {
			c.Redirect(http.StatusFound, "/")
			return order, user
		}
		if G.Msg.Fail == "" {
			G.Msg.Fail = "Some Error Occurred. Please Try Again Later."
		}
		c.Redirect(http.StatusFound, "/checkout")
		return order, user
	}

	order.OrderDetails = orderDetails

	redsysTransaction.PaidStatus = 1
	redsysTransaction.OrderID = order.ID
	R.SaveRedsysTransaction(redsysTransaction)

	DestroyCart(c)
	if !authSuccess || user.ID != redsysTransaction.UserID {
		c.Redirect(http.StatusFound, "/")
		return order, user
	}
	G.Msg.Success = `Thank you for your purchase. Kindly download the document with the bank details on your order below.`
	c.Redirect(http.StatusFound, "/account/orders")
	return order, user*/
}


func RedsysFail(c *gin.Context) {
	/*user, authSuccess :=M.IsAuthUser(c,G.FStore)

	//var success bool

	ord := c.Query("ord")
	encOrd := c.Query("encOrd")

	sc := securecookie.New([]byte(G.WalletApp.Key), nil)
	var value string
	err := sc.Decode("encryptedOrder", encOrd, &value)
	if err != nil {
		log.Println(err.Error())
		//c.Redirect(http.StatusFound, "/checkout")
		return
	}
	log.Println(ord, value)

	if ord != value {
		//c.Redirect(http.StatusFound, "/checkout")
		return
	}

	var redsysTransaction Mod.RedsysTransaction

	redsysTransaction = R.RedsysTransaction(redsysTransaction,"redsys_order=?",ord)
	if redsysTransaction.ID == 0 || redsysTransaction.PaidStatus > 0 {
		return
	}
	redsysTransaction.PaidStatus = 2
	R.SaveRedsysTransaction(redsysTransaction)

	if !authSuccess || user.ID != redsysTransaction.UserID {
		DestroyCart(c)
		c.Redirect(http.StatusFound, "/")
		return
	}
	G.Msg.Fail = "Checkout Failed. Please Try Again Later."
	c.Redirect(http.StatusFound, "/checkout")
	return*/
}
