package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)


func AddRedsysTransaction(redsysTransaction Mod.RedsysTransaction) bool{
	db := Cfg.DBConnect()
	err := db.Create(&redsysTransaction).Error
	if err == nil {
		return true
	}
	defer db.Close()
	return false
}


func AddRedsysOrderDetails(redsysOrderDetails []Mod.RedsysOrderDetail) ([]Mod.RedsysOrderDetail, bool){
	db := Cfg.DBConnect()
	for i, _ := range redsysOrderDetails {
		db.Create(&redsysOrderDetails[i])
		if redsysOrderDetails[i].ID == 0 {
			return redsysOrderDetails, false
		}
	}
	defer db.Close()
	return redsysOrderDetails, true
}


func RedsysTransactions(redsysTransactions []Mod.RedsysTransaction, where ... interface{}) []Mod.RedsysTransaction{
	db := Cfg.DBConnect()
	db.Find(&redsysTransactions, where...)
	for i, _ := range redsysTransactions {
		db.Find(&redsysTransactions[i].RedsysOrderDetails, "redsys_order=?", redsysTransactions[i].RedsysOrder)
	}
	defer db.Close()
	return redsysTransactions
}

func RedsysTransaction(redsysTransaction Mod.RedsysTransaction, where ... interface{}) Mod.RedsysTransaction{
	db := Cfg.DBConnect()
	db.Find(&redsysTransaction, where...)
	db.Find(&redsysTransaction.RedsysOrderDetails, "redsys_order=?", redsysTransaction.RedsysOrder)

	defer db.Close()
	return redsysTransaction
}

func UpdateRedsysTransaction(redsysTransaction Mod.RedsysTransaction,values interface{}, where ...interface{}) bool{
	db := Cfg.DBConnect()
	if db.Model(&redsysTransaction).Where(where).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func SaveRedsysTransaction(redsysTransaction Mod.RedsysTransaction) bool{
	if redsysTransaction.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Save(&redsysTransaction).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func DeleteRedsysTransaction(redsysTransaction Mod.RedsysTransaction) bool{
	if redsysTransaction.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	db = db.Begin()
	if db.Unscoped().Delete(&redsysTransaction).RowsAffected == 1 {
		db.Commit()
		defer db.Close()
		return true
	} else {
		db.Rollback()
		defer db.Close()
		return false
	}
}