package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)


func AddNewsletter(nl Mod.Newsletter) (Mod.Newsletter, bool){
	db := Cfg.DBConnect()
	err := db.Create(&nl).Error
	if err == nil {
		return nl, true
	}
	defer db.Close()
	return nl, false
}


func Newsletters(nls []Mod.Newsletter, where ... interface{}) []Mod.Newsletter{
	db := Cfg.DBConnect()
	db.Find(&nls, where...)
	defer db.Close()
	return nls
}

func Newsletter(nl Mod.Newsletter, where ... interface{}) Mod.Newsletter{
	db := Cfg.DBConnect()
	db.Find(&nl, where...)
	defer db.Close()
	return nl
}

func UpdateNewsletter(nl Mod.Newsletter,values interface{}, where ...interface{}) bool{
	db := Cfg.DBConnect()
	if db.Model(&nl).Where(where).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}

func DeleteNewsletter(nl Mod.Newsletter) bool{
	db := Cfg.DBConnect()
	if db.Delete(&nl).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}