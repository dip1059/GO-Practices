package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
	"github.com/jinzhu/gorm"
)

func AddOrderShippingAddress(db *gorm.DB, orderShippingAddress Mod.OrderShippingAddress) bool{
	db.Create(&orderShippingAddress)
	if orderShippingAddress.ID !=0 {
		return true
	}
	return false
}

func OrderShippingAddresses(orderShippingAddresses []Mod.OrderShippingAddress, where ... interface{}) []Mod.OrderShippingAddress {
	db := Cfg.DBConnect()
	db.Find(&orderShippingAddresses, where...)
	defer db.Close()
	return orderShippingAddresses
}


func OrderShippingAddress(orderShippingAddress Mod.OrderShippingAddress) Mod.OrderShippingAddress{
	db := Cfg.DBConnect()
	db.Find(&orderShippingAddress)
	defer db.Close()
	return orderShippingAddress
}


func UpdateOrderShippingAddress(orderShippingAddress Mod.OrderShippingAddress) bool {
	if orderShippingAddress.ID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if db.Save(&orderShippingAddress).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func DeleteOrderShippingAddress(orderShippingAddress Mod.OrderShippingAddress) bool{
	if orderShippingAddress.ID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if db.Unscoped().Delete(&orderShippingAddress).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}
