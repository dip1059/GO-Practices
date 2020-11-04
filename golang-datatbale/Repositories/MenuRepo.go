package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)


func AddMenu(menu Mod.Menu) bool{
	db := Cfg.DBConnect()
	err := db.Create(&menu).Error
	if err == nil {
		return true
	}
	defer db.Close()
	return false
}


func Menus(menus []Mod.Menu, where ... interface{}) []Mod.Menu{
	db := Cfg.DBConnect()
	db.Find(&menus, where...)
	for i, _ := range menus {
		if menus[i].PageUrl != "#" {
			db.Find(&menus[i].Page, "url=?", menus[i].PageUrl)
		}
	}
	defer db.Close()
	return menus
}

func Menu(menu Mod.Menu, where ... interface{}) Mod.Menu{
	db := Cfg.DBConnect()
	db.Find(&menu, where...)
	if menu.PageUrl != "#" {
		db.Find(&menu.Page, "url=?", menu.PageUrl)
	}
	defer db.Close()
	return menu
}

func UpdateMenu(menu Mod.Menu,values interface{}, where ...interface{}) bool{
	if menu.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Model(&menu).Where(where).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}

func SaveMenu(menu Mod.Menu) bool{
	if menu.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Save(&menu).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}

func DeleteMenu(menu Mod.Menu) bool{
	if menu.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Unscoped().Delete(&menu).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}