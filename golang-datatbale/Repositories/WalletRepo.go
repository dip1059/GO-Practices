package Repositories

import (
	"github.com/jinzhu/gorm"
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

func AddWallet(db *gorm.DB, wallet Mod.Wallet) Mod.Wallet{
	if db == nil {
		db = Cfg.DBConnect()
		db.Create(&wallet)
		defer db.Close()
		return wallet
	} else {
		db.Create(&wallet)
		return wallet
	}
}

func OnlyWallets(wallets []Mod.Wallet, where ... interface{}) []Mod.Wallet {
	db := Cfg.DBConnect()
	db.Order("updated_at desc").Find(&wallets, where...)
	for i, _ := range wallets {
		db.First(&wallets[i].Users, "id=?", wallets[i].UserID)
	}
	defer db.Close()
	return wallets
}

func OnlyWallet(wallet Mod.Wallet, where ... interface{}) Mod.Wallet{
	db := Cfg.DBConnect()
	db.First(&wallet, where...)
	defer db.Close()
	return wallet
}

func WalletWithOthers(wallet Mod.Wallet, where ... interface{}) Mod.Wallet{
	db := Cfg.DBConnect()
	db.First(&wallet, where...).Related(&wallet.WalletHistories)
	db.First(&wallet.Users, "id=?", wallet.UserID)
	defer db.Close()
	return wallet
}


func SaveWallet(db *gorm.DB, wallet Mod.Wallet) bool {
	if wallet.ID == 0 {
		return false
	}
	if db == nil {
		db = Cfg.DBConnect()
		if db.Save(&wallet).RowsAffected == 1{
			defer db.Close()
			return true
		} else {
			defer db.Close()
			return false
		}
	} else {
		if db.Save(&wallet).RowsAffected == 1{
			return true
		} else {
			return false
		}
	}

}


func DeleteWallet(wallet Mod.Wallet) bool{
	if wallet.ID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if db.Delete(&wallet).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}
