package Seeders

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
	"database/sql"
)

var payMethods = make([]Mod.PayMethod,0)

func PayMethodSeeder() {
	db := Cfg.DBConnect()

	payMethod1()
	payMethod2()

	for i,_ := range payMethods {
		db.Unscoped().FirstOrCreate(&payMethods[i],&Mod.PayMethod{Method:payMethods[i].Method})
	}
	defer db.Close()
}

func payMethod1() {
	var payMethod = Mod.PayMethod{
		Method: "Stripe",
		ImgUrl:sql.NullString{
			String:"/Public/User/images/stripe.png",
			Valid:true,
		},
		Status: 1,
	}
	payMethods = append(payMethods, payMethod)
}

func payMethod2() {
	var payMethod = Mod.PayMethod{
		Method: "Mollie",
		ImgUrl:sql.NullString {
			String:"/Public/User/images/mollie.png",
			Valid:true,
		},
		Status: 1,
	}
	payMethods = append(payMethods, payMethod)
}