package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)


func AddContactMessage(cm Mod.ContactMessage) (Mod.ContactMessage, bool){
	db := Cfg.DBConnect()
	err := db.Create(&cm).Error
	if err == nil {
		return cm, true
	}
	defer db.Close()
	return cm, false
}

func ContactMessages(cms []Mod.ContactMessage, where ... interface{}) []Mod.ContactMessage{
	db := Cfg.DBConnect()
	db.Find(&cms, where...)
	defer db.Close()
	return cms
}

func ContactMessage(cm Mod.ContactMessage, where ... interface{}) Mod.ContactMessage{
	db := Cfg.DBConnect()
	db.Find(&cm, where...)
	defer db.Close()
	return cm
}

func UpdateContactMessage(cm Mod.ContactMessage,values interface{}, where ...interface{}) bool{
	db := Cfg.DBConnect()
	if db.Model(&cm).Where(where).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}

func DeleteContactMessage(cm Mod.ContactMessage) bool{
	db := Cfg.DBConnect()
	if db.Delete(&cm).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}