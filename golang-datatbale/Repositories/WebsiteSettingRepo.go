package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)


func WebsiteSettings(wbsts []Mod.WebsiteSetting, where ... interface{}) []Mod.WebsiteSetting{
	db := Cfg.DBConnect()
	db.Find(&wbsts, where...)
	defer db.Close()
	return wbsts
}

func WebsiteSetting(wbst Mod.WebsiteSetting, where ... interface{}) Mod.WebsiteSetting{
	db := Cfg.DBConnect()
	db.Find(&wbst, where...)
	defer db.Close()
	return wbst
}

func UpdateWebsiteSetting(wbst Mod.WebsiteSetting,values interface{}, where ...interface{}) bool{
	db := Cfg.DBConnect()
	if db.Model(&wbst).Where(where).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}

func SaveWebsiteSetting(wbst Mod.WebsiteSetting) bool{
	if wbst.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Save(&wbst).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}

func DeleteWebsiteSetting(wbst Mod.WebsiteSetting) bool{
	db := Cfg.DBConnect()
	if db.Delete(&wbst).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}