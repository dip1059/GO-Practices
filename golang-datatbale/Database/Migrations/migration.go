package Migrtaions

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
	"github.com/jinzhu/gorm"
)


func Migrate() {
	db := Cfg.DBConnect()
	AlterColumns(db)
	db.AutoMigrate(&Mod.AdminSetting{})
	db.AutoMigrate(&Mod.Role{})
	db.AutoMigrate(&Mod.User{})
	db.AutoMigrate(&Mod.PasswordReset{})
	db.AutoMigrate(&Mod.Karat{})
	db.AutoMigrate(&Mod.Product{})
	db.AutoMigrate(&Mod.Wishlist{})
	db.AutoMigrate(&Mod.Bank{})
	db.AutoMigrate(&Mod.PayMethod{})
	db.AutoMigrate(&Mod.Order{})
	db.AutoMigrate(&Mod.OrderDetail{})
	db.AutoMigrate(&Mod.StripeTransaction{})
	db.AutoMigrate(&Mod.OrderDoc{})
	db.AutoMigrate(&Mod.Page{})
	db.AutoMigrate(&Mod.Menu{})
	db.AutoMigrate(&Mod.Newsletter{})
	db.AutoMigrate(&Mod.WebsiteSetting{})
	db.AutoMigrate(&Mod.ContactMessage{})
	db.AutoMigrate(&Mod.Invoice{})
	db.AutoMigrate(&Mod.Coupon{})
	db.AutoMigrate(&Mod.CouponOrder{})
	db.AutoMigrate(&Mod.HomeContent{})
	db.AutoMigrate(&Mod.Brand{})
	db.AutoMigrate(&Mod.EmailChangeHistory{})
	db.AutoMigrate(&Mod.MollieOrder{})
	db.AutoMigrate(&Mod.MollieOrderDetail{})
	db.AutoMigrate(&Mod.ShippingAddress{})
	db.AutoMigrate(&Mod.OrderShippingAddress{})
	db.AutoMigrate(&Mod.GoldCertificate{})
	db.AutoMigrate(&Mod.Wallet{})
	db.AutoMigrate(&Mod.WalletHistory{})
	db.AutoMigrate(&Mod.GoldTransfer{})
	db.AutoMigrate(&Mod.TransferShippingAddress{})
	db.AutoMigrate(&Mod.ProcessedTransfer{})

	AddForeignKeys(db)
	defer db.Close()

	/*db2 := Cfg.DBConnect2()
	db2.AutoMigrate(&Mod.User{})
	db2.AutoMigrate(&Mod.PasswordReset{})
	db2.AutoMigrate(&Mod.Session{})
	db2.AutoMigrate(&Mod.Tutorial{})
	defer db2.Close()*/
}


func AddForeignKeys(db *gorm.DB) {
	//db.Model(&User{}).AddUniqueIndex("idx_user_name_age", "name", "age")

	db.Model(&Mod.User{}).AddForeignKey("role_id", "roles(id)", "RESTRICT", "RESTRICT")
	db.Model(&Mod.Wishlist{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&Mod.Wishlist{}).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	db.Model(&Mod.Order{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	db.Model(&Mod.OrderDetail{}).AddForeignKey("order_id", "orders(id)", "RESTRICT", "RESTRICT")
	db.Model(&Mod.StripeTransaction{}).AddForeignKey("order_id", "orders(id)", "RESTRICT", "RESTRICT")
	db.Model(&Mod.OrderDoc{}).AddForeignKey("order_id", "orders(id)", "RESTRICT", "RESTRICT")

	db.Model(&Mod.CouponOrder{}).AddForeignKey("coupon_id", "coupons(id)", "RESTRICT", "RESTRICT")
	db.Model(&Mod.CouponOrder{}).AddForeignKey("order_id", "orders(id)", "RESTRICT", "RESTRICT")
}


func AlterColumns(db *gorm.DB) {
	/*db.Exec("ALTER TABLE orders CHANGE COLUMN delivery_status order_status tinyint(4) not null default 0;")
	db.Model(&Mod.PayMethod{}).ModifyColumn("description", "text")
	db.Model(&Mod.AdminSetting{}).ModifyColumn("value", "text")*/
	/*db.Exec("ALTER TABLE products MODIFY COLUMN max varchar(40) not null default '0';")
	db.Model(&Mod.OrderDetail{}).ModifyColumn("quantity", "double")*/
	//db.Exec("ALTER TABLE orders CHANGE COLUMN bank_ref_code unique_code varchar(255) not null default '0';")

}

