package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)


func AddPage(page Mod.Page) bool{
	db := Cfg.DBConnect()
	err := db.Create(&page).Error
	if err == nil {
		return true
	}
	defer db.Close()
	return false
}


func Pages(pages []Mod.Page, where ... interface{}) []Mod.Page{
	db := Cfg.DBConnect()
	db.Find(&pages, where...)
	defer db.Close()
	return pages
}

func Page(page Mod.Page, where ... interface{}) Mod.Page{
	db := Cfg.DBConnect()
	db.Find(&page, where...)
	defer db.Close()
	return page
}

func UpdatePage(page Mod.Page,values interface{}, where ...interface{}) bool{
	db := Cfg.DBConnect()
	if db.Model(&page).Where(where).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func SavePage(page Mod.Page) bool{
	if page.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Save(&page).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func DeletePage(page Mod.Page) bool{
	if page.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	db = db.Begin()
	if db.Unscoped().Delete(&page).RowsAffected == 1 {
		db.Commit()
		defer db.Close()
		return true
	} else {
		db.Rollback()
		defer db.Close()
		return false
	}
}