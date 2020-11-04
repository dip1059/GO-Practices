package Repositories

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)


func AddInvoice(invoice Mod.Invoice) Mod.Invoice{
	db := Cfg.DBConnect()

	if db.Find(&invoice, "order_id=?", invoice.OrderID).RecordNotFound() {
		err := db.Create(&invoice).Error
		if err == nil {
			return invoice
		}
	}

	defer db.Close()
	return invoice
}


func Invoices(invoices []Mod.Invoice, where ... interface{}) []Mod.Invoice{
	db := Cfg.DBConnect()
	db.Find(&invoices, where...)
	defer db.Close()
	return invoices
}

func Invoice(invoice Mod.Invoice, where ... interface{}) Mod.Invoice{
	db := Cfg.DBConnect()
	db.Find(&invoice, where...)
	defer db.Close()
	return invoice
}

func UpdateInvoice(invoice Mod.Invoice,values interface{}, where ...interface{}) bool{
	db := Cfg.DBConnect()
	if db.Model(&invoice).Where(where).Updates(values).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func SaveInvoice(invoice Mod.Invoice) bool{
	db := Cfg.DBConnect()
	if db.Save(&invoice).RowsAffected == 1 {
		defer db.Close()
		return true
	}
	defer db.Close()
	return false
}


func DeleteInvoice(invoice Mod.Invoice) bool{
	db := Cfg.DBConnect()
	if db.Unscoped().Delete(&invoice).RowsAffected == 1 {
		defer db.Close()
		return true
	} else {
		defer db.Close()
		return false
	}
}