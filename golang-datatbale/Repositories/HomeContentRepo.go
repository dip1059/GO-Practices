package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

func AddHomeContent(homeContent Mod.HomeContent) bool {
	db := Cfg.DBConnect()
	db.Create(&homeContent)
	if homeContent.ID !=0 {
		return true
	}
	defer db.Close()
	return false
}

func HomeContents(homeContents []Mod.HomeContent, where ... interface{}) []Mod.HomeContent {
	db := Cfg.DBConnect()
	db.Find(&homeContents, where...)
	defer db.Close()
	return homeContents
}


func HomeContent(homeContent Mod.HomeContent, where ... interface{}) Mod.HomeContent{
	db := Cfg.DBConnect()
	db.Find(&homeContent, where...)
	defer db.Close()
	return homeContent
}


func UpdateHomeContent(homeContent Mod.HomeContent) bool {
	if homeContent.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Save(&homeContent).RowsAffected == 1{
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func UpdatesHomeContent(homeContent Mod.HomeContent,values interface{}, whereQuery interface{}, whereArgs ... interface{}) bool{
	if homeContent.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Model(&homeContent).Where(whereQuery, whereArgs ...).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func DeleteHomeContent(homeContent Mod.HomeContent) bool{
	if homeContent.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Delete(&homeContent).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}
