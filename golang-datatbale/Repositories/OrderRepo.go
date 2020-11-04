package Repositories

import (
	"github.com/jinzhu/gorm"
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

func OrderDatable(fields interface{},sort interface{},offset interface{}, limit interface{}, whereQuery interface{}, args ...interface{}) ([]Mod.Order, int){
	db := Cfg.DBConnect()
	var count int
	db.Model(&Mod.Order{}).Where(whereQuery, args...).Count(&count)

	var orders []Mod.Order
	db.Order(sort).Select(fields).Offset(offset).Limit(limit).Where(whereQuery, args...).Find(&orders)
	for i, _ := range orders {
		db.Select("first_name, last_name").Find(&orders[i].User, "id=?", orders[i].UserID)
		db.Select("method").Find(&orders[i].PayMethod, "id=?", orders[i].PayMethodID)
	}
	defer db.Close()
	return orders, count
}


func Orders(orders []Mod.Order,sort string, where ...interface{}) []Mod.Order{
	db := Cfg.DBConnect()
	db.Order(sort).Find(&orders, where ...)
	for i, _ := range orders {
		db.Find(&orders[i].User, "id=?", orders[i].UserID)
		db.Find(&orders[i].PayMethod, "id=?", orders[i].PayMethodID)
		db.Find(&orders[i].OrderDetails, "order_id=?", orders[i].ID)
	}
	defer db.Close()
	return orders
}


func OnlyOrder(order Mod.Order, where ...interface{}) Mod.Order{
	db := Cfg.DBConnect()
	db.Find(&order, where ...)
	defer db.Close()
	return order
}


func Order(order Mod.Order, where ...interface{}) Mod.Order{
	db := Cfg.DBConnect()
	db.First(&order, where ...)
	db.First(&order.User, "id=?", order.UserID)
	db.First(&order.PayMethod, "id=?", order.PayMethodID)
	db.Find(&order.OrderDetails, "order_id=?", order.ID)
	db.First(&order.OrderShippingAddress, "order_id=?", order.ID)
	order.CouponOrder = CouponOrder(order.CouponOrder, "order_id=?", order.ID)
	if order.CouponOrder.Coupon.ID != 0 && order.CouponOrder.Coupon.Type == 1 {
		order.IsDiscountCoupon = true
	} else if order.CouponOrder.Coupon.ID != 0 && order.CouponOrder.Coupon.Type == 2 {
		order.IsFree = true
	}

	defer db.Close()
	return order
}


func UpdateOrder(order Mod.Order,values interface{}) bool {
	if order.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Model(&order).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func DeleteOrder(order Mod.Order) bool{
	if order.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Delete(&order).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


/*func OrderDetails(orderDetails []Mod.OrderDetail) []Mod.OrderDetail {
	db := Cfg.DBConnect()
	for i, _ := range orderDetails {
		db.Find(&orderDetails[i]).Related(&orderDetails[i].Product)
	}

	defer db.Close()
	return orderDetails
}*/

func AddOrderDoc(orderDoc Mod.OrderDoc) bool{
	db := Cfg.DBConnect()
	db.Create(&orderDoc)
	if orderDoc.ID !=0 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}



type CountStruct struct {
	AwaitingPayment int
	AwaitingShipment int
}

func CountDiffTypeOrders(userID uint) CountStruct {
	db := Cfg.DBConnect()
	var counts CountStruct
	db.Raw("select count(id) awaiting_payment from orders where user_id = ? and payment_status = 0", userID).Scan(&counts)
	db.Raw("select count(id) awaiting_shipment from orders where user_id = ? and order_status = 0", userID).Scan(&counts)
	return counts
}





//from ebuy

func AddOrder(db *gorm.DB,order Mod.Order) (Mod.Order, bool){
	db.Create(&order)
	//db.Model(&order).Related(&order.PayMethod)
	if order.ID !=0 {
		return order, true
	}
	return order, false
}

func OnlyOrders(orders []Mod.Order,sort string, limit int, offset int, where ...interface{}) []Mod.Order{
	db := Cfg.DBConnect()
	db.Order(sort).Offset(offset).Limit(limit).Find(&orders, where ...)
	defer db.Close()
	return orders
}

/*func Orders(orders []Mod.Order,sort string, where ...interface{}) []Mod.Order{
	db := Cfg.DBConnect()
	db.Order(sort).Find(&orders, where ...)
	for i, _ := range orders {
		db.Find(&orders[i].User, "id=?", orders[i].UserID)
		db.Find(&orders[i].PayMethod, "id=?", orders[i].PayMethodID)
		db.Find(&orders[i]).Related(&orders[i].OrderDetails)
		db.Find(&orders[i]).Related(&orders[i].OrderDocs)
		db.Find(&orders[i]).Related(&orders[i].ShippingAddress)
		db.Find(&orders[i]).Related(&orders[i].Invoice)
	}
	defer db.Close()
	return orders
}*/

type DateWiseOrder struct {
	Date string
	Orders []Mod.Order
}

func DateWiseOrders(offset int, limit int, where ... interface{}) []DateWiseOrder {
	db := Cfg.DBConnect()
	var dates []struct {
		Date string
	}
	db.Raw("select distinct `date` from (select id, date_format(created_at, '%d %M %Y') `date` from orders where "+where[0].(string)+" order by id desc limit ? offset ?) tbl",limit, offset).Scan(&dates)
	var dateWiseOrders = make([]DateWiseOrder, len(dates))
	for i, _ := range dates {
		dateWiseOrders[i].Date = dates[i].Date
		db.Raw("select * from orders where date_format(created_at, '%d %M %Y') = ? and "+where[0].(string)+" order by id desc", dates[i].Date).Scan(&dateWiseOrders[i].Orders)
		for j, _ := range dateWiseOrders[i].Orders {
			db.Find(&dateWiseOrders[i].Orders[j].PayMethod, "id=?", dateWiseOrders[i].Orders[j].PayMethodID)
			db.Find(&dateWiseOrders[i].Orders[j]).Related(&dateWiseOrders[i].Orders[j].OrderDetails)
			db.Find(&dateWiseOrders[i].Orders[j]).Related(&dateWiseOrders[i].Orders[j].OrderShippingAddress)
		}
	}

	defer db.Close()
	return dateWiseOrders
}


/*func Order(order Mod.Order, where ...interface{}) Mod.Order{
	db := Cfg.DBConnect()
	db.Find(&order, where ...)
	db.Find(&order.User, "id=?", order.UserID)
	db.Find(&order.PayMethod, "id=?", order.PayMethodID)
	db.Find(&order).Related(&order.OrderDetails)
	db.Find(&order).Related(&order.OrderDocs)
	db.Find(&order).Related(&order.ShippingAddress)
	db.Find(&order).Related(&order.Invoice)

	defer db.Close()
	defer db.Close()
	return order
}*/


func SaveUpdatedOrder(db *gorm.DB, order Mod.Order) bool {
	if order.ID == 0 {
		return false
	}
	if db.Save(&order).RowsAffected == 1 {
		return true
	}
	return false
}

/*func UpdateOrder(order Mod.Order,values interface{}) bool {
	if order.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Model(&order).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}*/


/*func DeleteOrder(order Mod.Order) bool{
	db := Cfg.DBConnect()
	if db.Delete(&order).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func AddOrderDoc(orderDoc Mod.OrderDoc) bool{
	db := Cfg.DBConnect()
	db.Create(&orderDoc)
	if orderDoc.ID !=0 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}*/




