package Repositories

import (
	"github.com/joho/godotenv"
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

type Result struct {
	TotalSales float64
}

type ThisYearSale struct {
	Month int
	Total float64
}

type CountrySale struct {
	Value float64	`json:"value"`
	Name string		`json:"name"`
}

func TotalUsers() int {
	db2 := Cfg.DBConnect()
	var users []Mod.User
	var count int
	db2.Find(&users).Where("role_id <> 1").Count(&count)
	defer db2.Close()
	return count
}

func TotalOrders() int {
	db := Cfg.DBConnect()
	var orders []Mod.Order
	var count int
	db.Unscoped().Find(&orders).Count(&count)
	defer db.Close()
	return count
}


func TotalProducts() int {
	db := Cfg.DBConnect()
	var product []Mod.Product
	var count int
	db.Find(&product).Where("").Count(&count)
	defer db.Close()
	return count
}


func TotalSales() float64 {
	db := Cfg.DBConnect()
	var result Result
	db.Table("orders").Select("convert(sum(grand_total)-1, signed) as total_sales").Where("order_status=1 or order_status=3").Scan(&result)
	defer db.Close()
	return result.TotalSales
}


func NewUsers() []Mod.User{
	db2 := Cfg.DBConnect()
	var users []Mod.User
	db2.Order("created_at desc").Where("role_id <> 1").Limit(5).Find(&users)
	defer db2.Close()
	return users
}

func ThisYearSales() []ThisYearSale{
	db := Cfg.DBConnect()
	var data []ThisYearSale
	db.Raw("select convert(month, signed) `month`, total from (SELECT DATE_FORMAT(updated_at, '%m') `month`, convert(sum(grand_total)-1, signed) total FROM `orders` where order_status=1 or order_status=3 and updated_at like CONCAT('%',DATE_FORMAT(now(), '%Y'),'%') group BY DATE_FORMAT(updated_at, '%m')) tbl order by month").Scan(&data)
	defer db.Close()
	return data
}

func CountrySales() []CountrySale {
	db := Cfg.DBConnect()
	db2 := db
	db = Cfg.DBConnect()
	var data []CountrySale
	godotenv.Load()
	db.Raw("select convert(sum(o.grand_total)-1, signed) value, ui.country_code `name`  from users ui, orders o where o.user_id = ui.id and ui.country_code is not null and (o.order_status = 1 or o.order_status = 3) and o.deleted_at is NULL GROUP by ui.country_code").Scan(&data)
	defer db.Close()
	defer db2.Close()
	return data
}
