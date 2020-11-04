package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
	"github.com/jinzhu/gorm"
)

func AddTransferShippingAddress(db *gorm.DB, transferShippingAddress Mod.TransferShippingAddress) bool{
	db.Create(&transferShippingAddress)
	if transferShippingAddress.ID !=0 {
		return true
	}
	return false
}

func TransferShippingAddresses(transferShippingAddresses []Mod.TransferShippingAddress, where ... interface{}) []Mod.TransferShippingAddress {
	db := Cfg.DBConnect()
	db.Find(&transferShippingAddresses, where...)
	defer db.Close()
	return transferShippingAddresses
}


func TransferShippingAddress(transferShippingAddress Mod.TransferShippingAddress) Mod.TransferShippingAddress{
	db := Cfg.DBConnect()
	db.Find(&transferShippingAddress)
	defer db.Close()
	return transferShippingAddress
}


func UpdateTransferShippingAddress(transferShippingAddress Mod.TransferShippingAddress) bool {
	if transferShippingAddress.ID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if db.Save(&transferShippingAddress).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func DeleteTransferShippingAddress(transferShippingAddress Mod.TransferShippingAddress) bool{
	if transferShippingAddress.ID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if db.Unscoped().Delete(&transferShippingAddress).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}
