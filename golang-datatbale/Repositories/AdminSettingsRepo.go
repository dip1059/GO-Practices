package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)


func AllAdminSettings(admSets []Mod.AdminSetting, where ... interface{}) []Mod.AdminSetting{
	db := Cfg.DBConnect()
	db.Find(&admSets, where...)
	defer db.Close()
	return admSets
}


func UpdateAdminSetting(admSet Mod.AdminSetting,values interface{}) bool{
	if admSet.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Model(&admSet).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}

func SaveAdminSetting(admSet Mod.AdminSetting) bool{
	if admSet.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Save(&admSet).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func DeleteAdminSetting(admSet Mod.AdminSetting) bool{
	if admSet.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Unscoped().Delete(&admSet).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func AddBank(bank Mod.Bank) bool{
	db := Cfg.DBConnect()
	db.Create(&bank)
	if bank.ID !=0 {
		return true
	}
	defer db.Close()
	return false
}


func AllBank(banks []Mod.Bank, where ... interface{}) []Mod.Bank{
	db := Cfg.DBConnect()
	db.Find(&banks, where...)
	defer db.Close()
	return banks
}


func Bank(bank Mod.Bank) Mod.Bank{
	db := Cfg.DBConnect()
	db.Find(&bank)
	defer db.Close()
	return bank
}

func BankExists(bank Mod.Bank, where ... interface{}) bool {
	db := Cfg.DBConnect()
	if db.Find(&bank, where...).RecordNotFound() {
		defer db.Close()
		return false
	}
	defer db.Close()
	return true
}


func UpdateBank(bank Mod.Bank) bool {
	if bank.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Save(&bank).RowsAffected == 1{
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func DeleteBank(bank Mod.Bank) bool{
	if bank.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Delete(&bank).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func UpdatePayMethod(payMethod Mod.PayMethod) bool {
	if payMethod.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Save(&payMethod).RowsAffected == 1{
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func DeletePayMethod(payMethod Mod.PayMethod) bool{
	if payMethod.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Delete(&payMethod).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}