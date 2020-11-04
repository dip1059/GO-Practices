package Seeders

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
	"database/sql"
)

var adminSettings =make([]Mod.AdminSetting,0)

func AdminSettingSeeder() {
	db := Cfg.DBConnect()

	adminSetting1()
	adminSetting2()
	adminSetting3()
	adminSetting4()
	adminSetting5()
	adminSetting6()
	adminSetting7()
	adminSetting8()
	adminSetting9()
	adminSetting10()
	adminSetting11()
	adminSetting12()
	adminSetting13()
	adminSetting14()
	adminSetting15()
	adminSetting16()


	for i,_ := range adminSettings {
		db.FirstOrCreate(&adminSettings[i],&Mod.AdminSetting{
			Slug:adminSettings[i].Slug,
		})
	}
	defer db.Close()
}

func adminSetting1() {
	var adminSetting = Mod.AdminSetting{
		Slug: "Company_Address",
		Value:sql.NullString{
			String: "451 Wall Street, Lisbon, Portugal",
			Valid: true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting2() {
	var adminSetting = Mod.AdminSetting{
		Slug: "Company_Phone",
		Value:sql.NullString{
			String: "(064) 332-1233",
			Valid: true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting3() {
	var adminSetting = Mod.AdminSetting{
		Slug: "Company_Email",
		Value:sql.NullString{
			String: "contact@trustgold999.com",
			Valid: true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting4() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Stripe_Secret_Key",
		Value:sql.NullString{
			String:"sk_test_5Kb3EwJX6Il6MqvjCRd3xrZV00Xz8H7IZc",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting5() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Stripe_Publishable_Key",
		Value:sql.NullString{
			String:"pk_test_ljRwJ1qHESBElROWXtQf4MEl00NHsDvTaf",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting6() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Mollie_Api_Key",
		Value:sql.NullString{
			String:"test_d9QynJzQ2pbjNE9RWH9DEsx4caUkDm",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting7() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Order_Fees_1_In_Percent",
		Value: sql.NullString{
			String:"0.0",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting8() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Order_Fees_1_Fixed_In_Euro",
		Value: sql.NullString{
			String:"0.0",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting9() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Order_Fees_2_In_Percent",
		Value: sql.NullString{
			String:"0.0",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting10() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Order_Fees_2_Fixed_In_Euro",
		Value: sql.NullString{
			String:"0.0",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting11() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Order_Fees_3_In_Percent",
		Value: sql.NullString{
			String:"0.0",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting12() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Order_Fees_3_Fixed_In_Euro",
		Value: sql.NullString{
			String:"0.0",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting13() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Transfer_Sender_Fees_In_Percent",
		Value: sql.NullString{
			String:"0.0",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting14() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Transfer_Sender_Fees_Fixed_In_Milligram",
		Value: sql.NullString{
			String:"0.0",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting15() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Transfer_Receiver_Fees_In_Percent",
		Value: sql.NullString{
			String:"0.0",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}

func adminSetting16() {
	var adminSetting = Mod.AdminSetting{
		Slug:"Transfer_Receiver_Fees_Fixed_In_Milligram",
		Value: sql.NullString{
			String:"0.0",
			Valid:true,
		},
	}
	adminSettings = append(adminSettings, adminSetting)
}