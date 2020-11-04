package Seeders

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

var wallets []Mod.Wallet

func WalletSeeder() {
	db := Cfg.DBConnect()

	wallet1()
	for i, _ := range wallets {
		db.Where(&Mod.Wallet{UserID: wallets[i].UserID}).FirstOrCreate(&wallets[i])
	}
	defer db.Close()
}

func wallet1() {
	wallet := Mod.Wallet{
		UserID: 1,
	}
	wallets = append(wallets, wallet)
}
