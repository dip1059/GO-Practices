package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)


func AllMethod(payMethods []Mod.PayMethod, where ... interface{}) []Mod.PayMethod{
	db := Cfg.DBConnect()
	db.Find(&payMethods, where...)
	defer db.Close()
	return payMethods
}

func PayMethod(payMethod Mod.PayMethod, where ... interface{}) Mod.PayMethod{
	db := Cfg.DBConnect()
	db.First(&payMethod, where...)
	defer db.Close()
	return payMethod
}

func PayMethodExists(payMethod Mod.PayMethod, where ... interface{}) bool {
	db := Cfg.DBConnect()
	if db.Find(&payMethod, where...).RecordNotFound() {
		defer db.Close()
		return false
	}
	defer db.Close()
	return true
}