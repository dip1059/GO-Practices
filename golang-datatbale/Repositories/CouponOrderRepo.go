package Repositories

import (
	"github.com/jinzhu/gorm"
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

/*func AddCouponOrder(couponOrder Mod.CouponOrder) (Mod.CouponOrder, bool) {
	db := Cfg.DBConnect()
	db.Create(&couponOrder)
	if couponOrder.ID !=0 {
		couponOrder.Coupon.ID = couponOrder.CouponID
		db.First(&couponOrder.Coupon)
		return couponOrder, true
	}
	defer db.Close()
	return couponOrder, false
}*/

func CouponOrders(couponOrders []Mod.CouponOrder, where ... interface{}) []Mod.CouponOrder {
	db := Cfg.DBConnect()
	db.Find(&couponOrders, where...)
	defer db.Close()
	return couponOrders
}

func CouponOrder(couponOrder Mod.CouponOrder, where ... interface{}) Mod.CouponOrder {
	db := Cfg.DBConnect()
	db.First(&couponOrder, where...)
	if couponOrder.ID == 0 {
		defer db.Close()
		return couponOrder
	}
	couponOrder.Coupon.ID = couponOrder.CouponID
	db.First(&couponOrder.Coupon)
	defer db.Close()
	return couponOrder
}

func SaveCouponOrder(couponOrder Mod.CouponOrder) bool {
	db := Cfg.DBConnect()
	if db.Save(&couponOrder).RowsAffected == 1{
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func DeleteCouponOrder(couponOrder Mod.CouponOrder) bool{
	db := Cfg.DBConnect()
	if db.Delete(&couponOrder).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func CouponUsageCount(userID uint, couponID uint) int {
	db := Cfg.DBConnect()
	var count int
	db.Find(&[]Mod.CouponOrder{}).Where( "user_id=? and coupon_id=?", userID, couponID).Count(&count)
	defer db.Close()
	return count
}




//from ebuy

func AddCouponOrder(db *gorm.DB, couponOrder Mod.CouponOrder) bool{
	db.Create(&couponOrder)
	if couponOrder.ID !=0 {
		return true
	}
	return false
}

/*func CouponOrders(couponOrders []Mod.CouponOrder, where ... interface{}) []Mod.CouponOrder {
	db := Cfg.DBConnect()
	db.Find(&couponOrders, where...)
	defer db.Close()
	return couponOrders
}

func CouponOrder(couponOrder Mod.CouponOrder, where ... interface{}) Mod.CouponOrder {
	db := Cfg.DBConnect()
	db.First(&couponOrder, where...)
	if couponOrder.ID == 0 {
		defer db.Close()
		return couponOrder
	}
	couponOrder.Coupon.ID = couponOrder.CouponID
	db.First(&couponOrder.Coupon)
	defer db.Close()
	return couponOrder
}

func SaveCouponOrder(couponOrder Mod.CouponOrder) bool {
	db := Cfg.DBConnect()
	if db.Save(&couponOrder).RowsAffected == 1{
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func DeleteCouponOrder(couponOrder Mod.CouponOrder) bool{
	db := Cfg.DBConnect()
	if db.Delete(&couponOrder).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func CouponUsageCount(userID uint, couponID uint) int {
	db := Cfg.DBConnect()
	var count int
	db.Find(&[]Mod.CouponOrder{}).Where( "user_id=? and coupon_id=?", userID, couponID).Count(&count)
	defer db.Close()
	return count
}*/