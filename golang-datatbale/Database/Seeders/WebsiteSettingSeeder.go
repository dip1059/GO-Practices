package Seeders

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

var wbsts = make([]Mod.WebsiteSetting,0)

func WebsiteSettingSeeder() {
	db := Cfg.DBConnect()

	wbst1()
	wbst2()
	wbst3()
	wbst4()
	wbst5()
	for i,_ := range wbsts {
		db.FirstOrCreate(&wbsts[i],&Mod.WebsiteSetting{
			ContentName:wbsts[i].ContentName},
		)
	}
	defer db.Close()
}

func wbst1() {
	var wbst = Mod.WebsiteSetting{
		ContentName: "Header Logo",
		Content: "/Public/User/images/logo.png",
		Status: 1,
	}
	wbsts = append(wbsts, wbst)
}

func wbst2() {
	var wbst = Mod.WebsiteSetting{
		ContentName: "Footer Logo",
		Content: "/Public/User/images/footer-logo.png",
		Status: 1,
	}
	wbsts = append(wbsts, wbst)
}


func wbst3() {
	var wbst = Mod.WebsiteSetting{
		ContentName: "Footer Text",
		Content: "Footer Text",
		Status: 1,
	}
	wbsts = append(wbsts, wbst)
}

func wbst4() {
	var wbst = Mod.WebsiteSetting{
		ContentName: "Copyright Text",
		Content: "Copyright 2019",
		Status: 1,
	}
	wbsts = append(wbsts, wbst)
}

func wbst5() {
	var wbst = Mod.WebsiteSetting{
		ContentName: "Login-Signup Header Logo",
		Content: "/Public/User/images/logo-light.png",
		Status: 1,
	}
	wbsts = append(wbsts, wbst)
}