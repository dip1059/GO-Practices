package Services

import (
	"github.com/gin-gonic/gin"
	"github.com/rikvdh/go-mollie"
	Cfg "gold-store/Config"
	G "gold-store/Globals"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	"log"
	"net/http"
	"strconv"
)

func PaymentRedirect(c *gin.Context) {
	user := M.GetAuthUser(c, G.FStore)
	log.Println("mollie payment redirected on order_id:",c.Query("order_id")," logged user:", user.Email, user.ID)
	ord, _ := strconv.Atoi(c.Query("order_id"))

	var order Mod.Order
	order.ID = uint(ord)
	order = R.OnlyOrder(order)
	if order.PaymentStatus == 1 && order.DeliveryType == 1 {
		G.Msg.Success = "Thank You, Your Order Has Been Placed And Completed Successfully."
	} else if order.PaymentStatus == 0 {
		G.Msg.Success = "Thank You, Your Order Has Been Placed Successfully. Your payment is under processing."
	} else if order.PaymentStatus == 3 {
		G.Msg.Success = "Sorry, Your Order Has Been Placed Successfully but the payment is cancelled. Make a new order please."
	} else {
		G.Msg.Success = "Thank You, Your Order Has Been Placed Successfully."
	}
	c.Redirect(http.StatusFound, "/account/orders")
	return
}

func PaymentWebhook(c *gin.Context) {

	id := c.PostForm("id")
	log.Println(id)
	payment := GetPayment(id)
	log.Println(*payment)

	metaData := payment.Metadata.(map[string]interface{})
	orderID, err := strconv.Atoi(metaData["order_id"].(string))
	if err != nil {
		log.Println(err.Error())
	}
	var order Mod.Order
	order.ID = uint(orderID - 10000)
	order = R.Order(order)

	if payment.Status == mollie.StatusPaid {
		order.PaidAmount = payment.Amount
		order.PaymentStatus = 1
		if order.DeliveryType == 1 {
			order.OrderStatus = 1
		}

		//GenerateCertificate(order)
		go func() {
			GenerateCertificate(order, c)
			//time.Sleep(time.Second * 20)
			//log.Println("go routine worked")
		}()

	} else {
		order.PaymentStatus = 3
		order.OrderStatus = 2
	}

	db := Cfg.DBConnect()
	if !R.SaveUpdatedOrder(db, order) {
		log.Println("Order not updated")
	}
	defer db.Close()

	var mo Mod.MollieOrder
	if !R.UpdateMollieOrder(mo, map[string]interface{}{"payment_status":payment.Status, "payment_method":payment.Method}, "payment_id = ?", id) {
		log.Println("MollieOrder not updated")
	}

	return
}