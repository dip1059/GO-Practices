package Admin

import (
	"github.com/gin-gonic/gin"
	G "gold-store/Globals"
	M "gold-store/Middlewares"
	Mod "gold-store/Models"
	R "gold-store/Repositories"
	S "gold-store/Services"
	"log"
	"net/http"
	"strconv"
)


func OrderDatatable(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}

	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "orders-datatables.html", map[string]interface{}{
			"AppEnv": G.AppEnv, "User": user,  "Nav":"orders", "Title": "Orders", "Msg": G.Msg})
		G.Msg.Success = ""
		G.Msg.Fail = ""

	} else if c.Request.Method == "POST" {
		Params := S.ProcessDatatableData(c)

		var whereQuery string
		var orders []Mod.Order
		var count int

		if Params["search"].(string) != "" {
			whereQuery = "id LIKE ? OR grand_total LIKE ?"
			orders, count = R.OrderDatable("id, grand_total, user_id, pay_method_id", Params["sort"], Params["offset"],
				Params["limit"], whereQuery, Params["arg"], Params["arg"])
		} else {
			orders, count = R.OrderDatable("id, grand_total, user_id, pay_method_id",
				Params["sort"], Params["offset"], Params["limit"], whereQuery)
		}

		Params["data"] = orders
		Params["count"] = count
		S.SendDatatableData(c, Params)
	}

}


func Orders(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var orders []Mod.Order
	orders = R.Orders(orders, "created_at desc")
	c.HTML(http.StatusOK, "orders.html", map[string]interface{}{
		"AppEnv": G.AppEnv, "User": user,  "Nav":"orders", "Title": "Orders", "Msg": G.Msg, "Orders": orders, "UserID":"orders"})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func MakeOrderCancel(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var order Mod.Order
	order.ID = uint(id)
	if R.UpdateOrder(order, map[string]interface{}{"order_status": 2}) {
		G.Msg.Success = "Order Status Updated Successfully"
	} else {
		G.Msg.Fail = "Some Error Occurred,Order Status Update Failed. Please Try Again Later."
	}

	UserID := c.Param("userId")
	if UserID != "orders" {
		c.Redirect(http.StatusFound, "/user-orders/"+UserID)
	} else {
		c.Redirect(http.StatusFound, "/orders")
	}
}


func AddTrackCode(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.PostForm("id"))
	code := c.PostForm("track_code")
	var order Mod.Order
	order.ID = uint(id)
	if R.UpdateOrder(order, map[string]interface{}{"track_code": code, "order_status": 3}) {
		G.Msg.Success = "Track Code Added Successfully"

		order = R.Order(order)
		success := S.SendOrderEmail(order, c, "")
		log.Println("Order Email Sending Success:", success)
	} else {
		G.Msg.Fail = "Some Error Occurred. Please Try Again Later."
	}

	UserID := c.Param("userId")
	if UserID != "orders" {
		c.Redirect(http.StatusFound, "/user-orders/"+UserID)
	} else {
		c.Redirect(http.StatusFound, "/orders")
	}
}


func MakeOrderComplete(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var ord Mod.Order
	ord.ID = uint(id)
	ord = R.Order(ord)
	ord.CouponOrder = R.CouponOrder(ord.CouponOrder, "order_id=?", ord.ID)
	if ord.CouponOrder.Coupon.ID != 0 && ord.CouponOrder.Coupon.Type == 1 {
		ord.IsDiscountCoupon = true
	}
	success = true//S.SendBidToWallet(ord, c)
	log.Println(success)

	var order Mod.Order
	order.ID = ord.ID
	if success && R.UpdateOrder(order, map[string]interface{}{"order_status": 1}) {
		G.Msg.Success = "Order Status Updated And An Email Was Sent To Customer's Email."
		ord.OrderStatus = 1
		filename, success := S.GenerateInvoice(ord)
		/*log.Println("Invoice Generation Success:", success)*/
		success = S.SendOrderEmail(ord, c, filename)
		log.Println("Order Completion Email Sending Success:", success)
	} else {
		G.Msg.Fail = "Some Error Occurred, Order Status Update Failed And Bid Not Sent To Wallet. Please Try Again Later."
	}

	UserID := c.Param("userId")
	if UserID != "orders" {
		c.Redirect(http.StatusFound, "/user-orders/"+UserID)
	} else {
		c.Redirect(http.StatusFound, "/orders")
	}

}


func OrderDetails(c *gin.Context) {
	user, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	var order Mod.Order
	id, _ := strconv.Atoi(c.Param("id"))
	UserID := c.Param("userId")

	order.ID = uint(id)
	order = R.Order(order)
	order.OrderShippingAddress.Country = G.Country[order.OrderShippingAddress.Country]

	order.User.CountryCode = G.Country[order.User.CountryCode]
	order.OrderDetails = R.OrderDetails(order.OrderDetails)
	order.CouponOrder = R.CouponOrder(order.CouponOrder, "order_id=?", order.ID)
	lenDoc := len(order.OrderDocs)
	c.HTML(http.StatusOK, "order-details.html",map[string]interface{}{
		"AppEnv":G.AppEnv, "User":user,  "Nav":"orders", "Title":"Order-Details", "Order":order, "Msg":G.Msg, "UserID":UserID, "LenDoc":lenDoc})
	G.Msg.Success = ""
	G.Msg.Fail = ""
}


func DeleteOrder(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var order Mod.Order
	order.ID = uint(id)
	if R.DeleteOrder(order) {
		G.Msg.Success = "Order Deleted Successfully."
	} else {
		G.Msg.Fail = "Some Error Occurred, Status Update Failed. Please Try Again Later."
	}

	UserID := c.Param("userId")
	if UserID != "orders" {
		c.Redirect(http.StatusFound, "/user-orders/"+UserID)
	} else {
		c.Redirect(http.StatusFound, "/orders")
	}
}

func MakePaymentPending(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var order Mod.Order
	order.ID = uint(id)
	if R.UpdateOrder(order, map[string]interface{}{"payment_status": 0}) {
		G.Msg.Success = "Payment Status Updated Successfully."
	} else {
		G.Msg.Fail = "Some Error Occurred,Payment Status Update Failed. Please Try Again Later."
	}

	UserID := c.Param("userId")
	if UserID != "orders" {
		c.Redirect(http.StatusFound, "/user-orders/"+UserID)
	} else {
		c.Redirect(http.StatusFound, "/orders")
	}
}


func MakePaymentDone(c *gin.Context) {
	_, success := M.IsAuthAdminUser(c, G.FStore)
	if !success {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var order Mod.Order
	order.ID = uint(id)
	order = R.Order(order)
	if R.UpdateOrder(order, map[string]interface{}{"payment_status": 1, "paid_amount": order.GrandTotal}) {
		G.Msg.Success = "Payment Status Updated Successfully."
	} else {
		G.Msg.Fail = "Some Error Occurred,Payment Status Update Failed. Please Try Again Later."
	}

	UserID := c.Param("userId")
	if UserID != "orders" {
		c.Redirect(http.StatusFound, "/user-orders/"+UserID)
	} else {
		c.Redirect(http.StatusFound, "/orders")
	}
}
