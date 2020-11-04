package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)


func AddLusopay(lusopay Mod.Lusopay) bool{
	db := Cfg.DBConnect()
	err := db.Create(&lusopay).Error
	if err == nil {
		return true
	}
	defer db.Close()
	return false
}


func Lusopays(lusopays []Mod.Lusopay, where ... interface{}) []Mod.Lusopay{
	db := Cfg.DBConnect()
	db.Find(&lusopays, where...)
	defer db.Close()
	return lusopays
}

func Lusopay(lusopay Mod.Lusopay, where ... interface{}) Mod.Lusopay{
	db := Cfg.DBConnect()
	db.Find(&lusopay, where...)
	defer db.Close()
	return lusopay
}

func UpdateLusopay(lusopay Mod.Lusopay,values interface{}, where ...interface{}) bool{
	if lusopay.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Model(&lusopay).Where(where).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func SaveLusopay(lusopay Mod.Lusopay) bool{
	if lusopay.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Save(&lusopay).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func DeleteLusopay(lusopay Mod.Lusopay) bool{
	if lusopay.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	db = db.Begin()
	if db.Unscoped().Delete(&lusopay).RowsAffected == 1 {
		db.Commit()
		defer db.Close()
		return true
	} else {
		db.Rollback()
		defer db.Close()
		return false
	}
}