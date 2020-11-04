package Repositories

import (
	"github.com/jinzhu/gorm"
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

func AddWalletHistory(db *gorm.DB, walletHistory Mod.WalletHistory) bool{
	if db == nil {
		db = Cfg.DBConnect()
		db.Create(&walletHistory)
		if walletHistory.ID != 0 {
			defer db.Close()
			return true
		}
		defer db.Close()
		return false
	} else {
		db.Create(&walletHistory)
		if walletHistory.ID != 0 {
			return true
		}
		return false
	}
}

func WalletHistories(walletHistories []Mod.WalletHistory, where ... interface{}) []Mod.WalletHistory {
	db := Cfg.DBConnect()
	db.Find(&walletHistories, where...)
	defer db.Close()
	return walletHistories
}


func WalletHistoryWithProducts(walletHistory Mod.WalletHistory) Mod.WalletHistory{
	db := Cfg.DBConnect()
	db.Find(&walletHistory)
	defer db.Close()
	return walletHistory
}


func SaveWalletHistory(walletHistory Mod.WalletHistory) bool {
	if walletHistory.ID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if db.Save(&walletHistory).RowsAffected == 1{
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func DeleteWalletHistory(walletHistory Mod.WalletHistory) bool{
	if walletHistory.ID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if db.Delete(&walletHistory).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}
