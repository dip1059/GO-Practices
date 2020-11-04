package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
	"github.com/jinzhu/gorm"
)

func AddOrderDetail(db *gorm.DB, orderDetail Mod.OrderDetail) bool{
	db.Create(&orderDetail)
	if orderDetail.ID == 0 {
		db.Rollback()
		
		return false
	}

	return true
}

func AddOrderDetails(orderDetails []Mod.OrderDetail) ([]Mod.OrderDetail, bool){
	db := Cfg.DBConnect()
	for i, _ := range orderDetails {
		db.Create(&orderDetails[i])
		if orderDetails[i].ID == 0 {
			return orderDetails, false
		}
	}
	defer db.Close()
	return orderDetails, true
}


func OrderDetails(orderDetails []Mod.OrderDetail) []Mod.OrderDetail {
	db := Cfg.DBConnect()
	for i, _ := range orderDetails {
		db.Find(&orderDetails[i]).Related(&orderDetails[i].Product)
	}

	defer db.Close()
	return orderDetails
}


func OrderDetail(orderDetail Mod.OrderDetail) Mod.OrderDetail{
	db := Cfg.DBConnect()
	db.Find(&orderDetail)
	defer db.Close()
	return orderDetail
}


func UpdateOrderDetail(orderDetail Mod.OrderDetail) bool {
	db := Cfg.DBConnect()
	if db.Save(&orderDetail).RowsAffected == 1{
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func DeleteOrderDetail(orderDetail Mod.OrderDetail) bool{
	db := Cfg.DBConnect()
	if db.Delete(&orderDetail).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}

func HasReviewPermission(proID int, userID uint) bool {
	db := Cfg.DBConnect()
	var ord Mod.Order
	db.Raw("select id, user_id from orders where id in (select distinct order_id from order_details where product_id=?) and user_id=? and delivery_status=? limit 1",proID, userID,1).Scan(&ord)
	if ord.ID == 0 {
		defer db.Close()
		return false
	}

	defer db.Close()
	return true
}
