package Repositories

import (
	"github.com/jinzhu/gorm"
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
	"math"
)

func AddProduct(product Mod.Product) bool{
	db := Cfg.DBConnect()
	db.Create(&product)
	if product.ID !=0 {
		return true
	}
	defer db.Close()
	return false
}


func Product(product Mod.Product, userID uint, where ... interface{}) Mod.Product {
	db := Cfg.DBConnect()
	db.First(&product, where ...)
	if product.ID > 0 {
		db.Find(&product.Karat, "id=?", product.KaratID)
		product.DiscountPrice = math.Round((product.Price-product.Price*product.Discount/100)*100) / 100
		product.IsInWishlist = IsInWishlist(db, userID, product.ID)
	}
	defer db.Close()
	return product
}


func ProductsWithOthers(products []Mod.Product, userID uint, where ... interface{}) []Mod.Product {
	db := Cfg.DBConnect()
	db.Order("id desc").Find(&products, where ...)
	for i, _ := range products {
		db.Find(&products[i].Karat, "id=?", products[i].KaratID)
		products[i].DiscountPrice = math.Round((products[i].Price - products[i].Price * products[i].Discount / 100) * 100) / 100
		products[i].IsInWishlist = IsInWishlist(db, userID, products[i].ID)
	}
	defer db.Close()
	return products
}


func IsInWishlist(db *gorm.DB, userID uint, proID uint) bool {
	var wish Mod.Wishlist
	if db.Find(&wish, "user_id=? and product_id=?", userID, proID).RecordNotFound() {
		return false
	}
	return true
}


func UpdateProduct(product Mod.Product) bool {
	if product.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Save(&product).RowsAffected == 1{
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}


func DeleteProduct(product Mod.Product) bool{
	if product.ID == 0 {
		return false
	}
	db := Cfg.DBConnect()
	if db.Delete(&product).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}
