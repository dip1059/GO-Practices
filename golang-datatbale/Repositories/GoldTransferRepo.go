package Repositories

import (
	Cfg "gold-store/Config"
	G "gold-store/Globals"
	Mod "gold-store/Models"
)

func AddGoldTransfer(transfer Mod.GoldTransfer) bool {
	db := Cfg.DBConnect()
	db = db.Begin()

	err := db.Create(&transfer).Error
	if err != nil {
		db.Rollback()
		defer db.Close()
		return false
	}
	if transfer.DeliveryType == 2 {
		transfer.ShipAddress.GoldTransferID = transfer.ID
		if !AddTransferShippingAddress(db, transfer.ShipAddress) {
			db.Rollback()
			defer db.Close()
			return false
		}
	}
	db.Commit()
	defer db.Close()
	return true
}

func GoldTransfers(transfers []Mod.GoldTransfer, where ...interface{}) []Mod.GoldTransfer {
	db := Cfg.DBConnect()
	db.Order("updated_at desc").Find(&transfers, where...)
	for i, _ := range transfers {
		db.Raw(`select * from users u, wallets w where u.id=w.user_id and w.id=?`,transfers[i].SenderWalletID).Scan(&transfers[i].SenderUser)
		db.Raw(`select * from users u, wallets w where u.id=w.user_id and w.id=?`, transfers[i].ReceiverWalletID).Scan(&transfers[i].ReceiverUser)
	}
	defer db.Close()
	return transfers
}

func GoldTransfer(transfer Mod.GoldTransfer, where ...interface{}) Mod.GoldTransfer {
	db := Cfg.DBConnect()
	db.First(&transfer, where...)
	db.Raw(`select * from users u, wallets w where u.id=w.user_id and w.id=?`,transfer.SenderWalletID).Scan(&transfer.SenderUser)
	db.Raw(`select * from users u, wallets w where u.id=w.user_id and w.id=?`, transfer.ReceiverWalletID).Scan(&transfer.ReceiverUser)
	if transfer.DeliveryType == 2 {
		db.First(&transfer.ShipAddress, "gold_transfer_id=?", transfer.ID)
		transfer.ShipAddress.Country = G.Country[transfer.ShipAddress.Country]
	}
	defer db.Close()
	return transfer
}


func UpdateGoldTransfer(gt Mod.GoldTransfer,values interface{}) bool {
	if gt.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Model(&gt).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func SaveGoldTransfer(transfer Mod.GoldTransfer) bool {
	if transfer.ID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if db.Save(&transfer).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}

func DeleteGoldTransfer(transfer Mod.GoldTransfer) bool {
	if transfer.ID == 0 {
		return false
	}

	db := Cfg.DBConnect()
	if db.Delete(&transfer).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}
