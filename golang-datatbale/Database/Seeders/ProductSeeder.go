package Seeders

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

var products =make([]Mod.Product,0)

func ProductSeeder() {
	db := Cfg.DBConnect()

	//product1()
	for i,_ := range products {
		db.FirstOrCreate(&products[i],&Mod.Product{Type:products[i].Type,Title:products[i].Title})
	}
	defer db.Close()
}

/*func product1() {
	var product = Mod.Product{
		Title: "Free Bids",
		ImgUrl:sql.NullString{
			String:"/Public/User/images/free-bids.jpg",
			Valid:true,
		},
		GrmAmount:0.0,
		Type:3,
		Price:0.0,
		Status: 1,
	}
	products = append(products, product)
}*/

