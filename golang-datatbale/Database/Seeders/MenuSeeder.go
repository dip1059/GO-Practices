package Seeders

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
	"database/sql"
)

var menus = make([]Mod.Menu,0)

func MenuSeeder() {
	db := Cfg.DBConnect()

	menu1()
	menu2()
	menu3()
	menu4()
	menu5()
	menu6()
	menu7()
	for i,_ := range menus {
		db.FirstOrCreate(&menus[i])
	}
	defer db.Close()
}

func menu1() {
	var menu = Mod.Menu{
		Title: "Footer Left Menu Header",
		Type: 0,
		Position: sql.NullString{
			String: "Footer Left",
			Valid:true,
		},
		Status: 1,
	}
	menu.ID = 1
	menus = append(menus, menu)
}

func menu2() {
	var menu = Mod.Menu{
		Title: "Footer Middle Menu Header",
		Type: 0,
		Position: sql.NullString{
			String: "Footer Middle",
			Valid:true,
		},
		Status: 1,
	}
	menu.ID = 2
	menus = append(menus, menu)
}


func menu3() {
	var menu = Mod.Menu{
		Title: "Footer Right Menu Header",
		Type: 0,
		Position: sql.NullString{
			String: "Footer Right",
			Valid:true,
		},
		Status: 1,
	}
	menu.ID = 3
	menus = append(menus, menu)
}

func menu4() {
	var menu = Mod.Menu{
		Title: "Facebook",
		Type: 1,
		Status: 1,
	}
	menu.ID = 4
	menus = append(menus, menu)
}

func menu5() {
	var menu = Mod.Menu{
		Title: "YouTube",
		Type: 1,
		Status: 1,
	}
	menu.ID = 5
	menus = append(menus, menu)
}

func menu6() {
	var menu = Mod.Menu{
		Title: "Twitter",
		Type: 1,
		Status: 1,
	}
	menu.ID = 6
	menus = append(menus, menu)
}

func menu7() {
	var menu = Mod.Menu{
		Title: "Telegram",
		Type: 1,
		Status: 1,
	}
	menu.ID = 7
	menus = append(menus, menu)
}