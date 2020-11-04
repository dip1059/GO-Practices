package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

func AddShippingAddress(shippingAddress Mod.ShippingAddress) bool{
	db := Cfg.DBConnect()
	db.Create(&shippingAddress)
	if shippingAddress.ID !=0 {
		return true
	}
	defer db.Close()
	return false
}

func ShippingAddresses(shippingAddresses []Mod.ShippingAddress, where ... interface{}) []Mod.ShippingAddress {
	db := Cfg.DBConnect()
	db.Order("status desc").Find(&shippingAddresses, where...)
	defer db.Close()
	return shippingAddresses
}


func ShippingAddress(shippingAddress Mod.ShippingAddress, where ... interface{}) Mod.ShippingAddress{
	db := Cfg.DBConnect()
	db.First(&shippingAddress, where...)
	defer db.Close()
	return shippingAddress
}


func SaveShippingAddress(shippingAddress Mod.ShippingAddress, isDefault bool) bool {
	if shippingAddress.UserID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if isDefault {
		db.Model(&Mod.ShippingAddress{}).Where("user_id=?", shippingAddress.UserID).Updates(map[string]interface{}{
			"status":0})
	}

	if db.Save(&shippingAddress).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func DeleteShippingAddress(shippingAddress Mod.ShippingAddress, where ... interface{}) bool{
	if shippingAddress.ID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if db.Unscoped().Delete(&shippingAddress, where ... ).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}
