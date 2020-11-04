package Repositories

import (
	"github.com/jinzhu/gorm"
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

func AddGoldCertificate(db *gorm.DB, goldCertificate Mod.GoldCertificate) bool {
	err := db.Create(&goldCertificate).Error
	if err != nil {
		db.Rollback()
		return false
	}
	return true
}

func ExecuteCertificate(db *gorm.DB, query string) bool {
	err := db.Exec(query).Error
	if err != nil {
		db.Rollback()
		return false
	}
	return true
}

func AddProcessedTransfer(db *gorm.DB, tr Mod.ProcessedTransfer) bool {
	err := db.Create(&tr).Error
	if err != nil {
		db.Rollback()
		return false
	}
	return true
}

/*func GenerateCertificate(query string) bool {
	db := Cfg.DBConnect()
	db = db.Begin()
	err := db.Exec(query).Error;
	if err != nil {
		db.Rollback()
		return false
	}
	db.Commit()
	return true
}*/

type CountCertificate struct {
	 Total int
}

func CountWalletCertificate(userID uint) int {
	db := Cfg.DBConnect()
	var cc CountCertificate
	db.Raw(`select count(id) total from gold_certificates where user_id=? and delivery_type=1`, userID).Scan(&cc)
	defer db.Close()
	return cc.Total
}

func GoldCertificates(goldCertificates []Mod.GoldCertificate, where ...interface{}) []Mod.GoldCertificate {
	db := Cfg.DBConnect()
	db.Find(&goldCertificates, where...)
	defer db.Close()
	return goldCertificates
}

func GoldCertificate(goldCertificate Mod.GoldCertificate, where ...interface{}) Mod.GoldCertificate {
	db := Cfg.DBConnect()
	db.Find(&goldCertificate, where...)
	defer db.Close()
	return goldCertificate
}

func SaveGoldCertificate(goldCertificate Mod.GoldCertificate) bool {
	if goldCertificate.ID == "" {
		return false
	}

	db := Cfg.DBConnect()
	if db.Save(&goldCertificate).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}

func DeleteGoldCertificate(goldCertificate Mod.GoldCertificate) bool {
	if goldCertificate.ID == "" {
		return false
	}

	db := Cfg.DBConnect()
	if db.Delete(&goldCertificate).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}
