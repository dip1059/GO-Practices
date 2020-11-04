package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

func AddCoupon(coupon Mod.Coupon) bool{
	db := Cfg.DBConnect()
	db.Create(&coupon)
	if coupon.ID !=0 {
		return true
	}
	defer db.Close()
	return false
}

func Coupons(coupons []Mod.Coupon, where ... interface{}) []Mod.Coupon {
	db := Cfg.DBConnect()
	db.Order("id desc").Find(&coupons, where...)
	defer db.Close()
	return coupons
}

func CouponExists(code string) bool {
	db := Cfg.DBConnect()
	if db.First(&Mod.Coupon{}, "code=?", code).RecordNotFound() {
		defer db.Close()
		return false
	}
	defer db.Close()
	return true
}

func Coupon(coupon Mod.Coupon, where ... interface{}) Mod.Coupon {
	db := Cfg.DBConnect()
	db.First(&coupon, where...)
	defer db.Close()
	return coupon
}

func SaveCoupon(coupon Mod.Coupon) bool {
	db := Cfg.DBConnect()
	if db.Save(&coupon).RowsAffected == 1{
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func DeleteCoupon(coupon Mod.Coupon) bool{
	db := Cfg.DBConnect()
	if db.Delete(&coupon).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}
