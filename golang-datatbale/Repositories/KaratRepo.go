package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

func AddKarat(karat Mod.Karat) bool{
	db := Cfg.DBConnect()
	db.Create(&karat)
	if karat.ID !=0 {
		return true
	}
	defer db.Close()
	return false
}

func Karats(karats []Mod.Karat, where ... interface{}) []Mod.Karat {
	db := Cfg.DBConnect()
	db.Find(&karats, where...)
	defer db.Close()
	return karats
}


func KaratWithProducts(karat Mod.Karat) Mod.Karat{
	db := Cfg.DBConnect()
	db.Find(&karat).Related(&karat.Products)
	defer db.Close()
	return karat
}


func SaveKarat(karat Mod.Karat) bool {
	if karat.ID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if db.Save(&karat).RowsAffected == 1{
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func DeleteKarat(karat Mod.Karat) bool{
	if karat.ID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if db.Delete(&karat).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}
