package Repositories

import (
	"github.com/jinzhu/gorm"
	Cfg "gold-store/Config"
	//Cfg "gold-store/Config"
	Mod "gold-store/Models"
)


/*func AddStripeTransaction(strpTx Mod.StripeTransaction) bool{
	db := Cfg.DBConnect()
	db.Create(&strpTx)
	if strpTx.ID == 0 {
		defer db.Close()
		return false
	}
	defer db.Close()
	return true
}

func AddBTTx(btTx Mod.BrainTreeTransaction) bool{
	db := Cfg.DBConnect()
	db.Create(&btTx)
	if btTx.ID == 0 {
		defer db.Close()
		return false
	}
	defer db.Close()
	return true
}*/




//from ebuy

func AddStripeTransaction(db *gorm.DB, strpTx Mod.StripeTransaction) bool{
	db.Create(&strpTx)
	if strpTx.ID == 0 {
		return false
	}
	return true
}
//

func AddMollieOrder(db *gorm.DB, mo Mod.MollieOrder) bool{
	db.Create(&mo)
	if mo.ID == 0 {
		return false
	}
	return true
}

func UpdateMollieOrder(order Mod.MollieOrder,values interface{}, whereQuery string, args ... interface{}) bool {
	if order.ID == 0 && whereQuery == "" && args == nil {
		return false
	}
	db := Cfg.DBConnect()
	if db.Model(&order).Where(whereQuery,args ...).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}
