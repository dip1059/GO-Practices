package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
	"math"
)


func UsersWithOthers(users []Mod.User, where ... interface{}) []Mod.User{
	db2 := Cfg.DBConnect()
	db2.Order("created_at desc").Find(&users, where...)
	defer db2.Close()
	return users
}


func UpdateUser(user Mod.User,values interface{}) bool{
	if user.ID == 0 {
		return false
	}
	db2 := Cfg.DBConnect()
	if db2.Model(&user).Updates(values).RowsAffected == 1 {
		defer db2.Close()
		return true
	}
	defer db2.Close()
	return false
}


func DeleteUser(user Mod.User) bool{
	if user.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Delete(&user).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func PayMethods(payMethods []Mod.PayMethod, where ...interface{}) []Mod.PayMethod{
	db := Cfg.DBConnect()
	db.Find(&payMethods, where...)
	defer db.Close()
	return payMethods
}


/*func AddOrder(order Mod.Order) (Mod.Order, bool){
	db := Cfg.DBConnect()
	db.Create(&order)
	db.Model(&order).Related(&order.PayMethod)
	if order.ID !=0 {
		return order, true
	}
	defer db.Close()
	return order, false
}*/


/*func AddOrderDetails(orderDetails []Mod.OrderDetail) ([]Mod.OrderDetail, bool){
	db := Cfg.DBConnect()
	for i, _ := range orderDetails {
		db.Create(&orderDetails[i])
		if orderDetails[i].ID == 0 {
			return orderDetails, false
		}
	}
	defer db.Close()
	return orderDetails, true
}*/


func AddWish(wish Mod.Wishlist) bool{
	db := Cfg.DBConnect()
	db.Create(&wish)
	if wish.ID !=0 {
		return true
	}
	defer db.Close()
	return false
}

func CountWishlist(user Mod.User) int{
	db := Cfg.DBConnect()
	var wishlist []Mod.Wishlist
	var count int
	db.Find(&wishlist).Where("user_id=?", user.ID).Count(&count)
	defer db.Close()
	return count
}


func Wishlist(wishlist []Mod.Wishlist, user Mod.User) []Mod.Wishlist {
	db := Cfg.DBConnect()
	var wishes []Mod.Wishlist
	db.Find(&wishlist, "user_id=?", user.ID)
	for i, _ := range wishlist {
		db.Find(&wishlist[i].Product, "id=? and status=?", wishlist[i].ProductID,1)
		if wishlist[i].Product.Status == 1 {
			db.Find(&wishlist[i].Product.Karat, "id=?", wishlist[i].Product.KaratID)
			wishlist[i].Product.DiscountPrice = math.Ceil((wishlist[i].Product.Price-wishlist[i].Product.Price*wishlist[i].Product.Discount/100)*100) / 100
			wishes = append(wishes, wishlist[i])
		}
	}
	defer db.Close()
	return wishes
}


func RemoveFromWishlist(wish Mod.Wishlist, where ... interface{}) bool {
	if wish.ID == 0 && where == nil {
		return false
	}
	db := Cfg.DBConnect()
	db = db.Begin()
	if db.Unscoped().Delete(&wish, where ...).RowsAffected == 1 {
		db.Commit()
		defer db.Close()
		return true
	}
	db.Rollback()
	defer db.Close()
	return false
}

func SaveUserChanges(user Mod.User) bool{
	if user.ID == 0 {
		return false
	}
	db2 := Cfg.DBConnect()
	if db2.Save(&user).RowsAffected == 1{
		defer db2.Close()
		return true
	} else {
		defer db2.Close()
		return false
	}
}

func CreateEmailChangeHistory(emailChangeHistort Mod.EmailChangeHistory) bool{
	db := Cfg.DBConnect()
	err := db.Create(&emailChangeHistort).Error
	if err != nil {
		defer db.Close()
		return false
	} else {
		defer db.Close()
		return true
	}
}