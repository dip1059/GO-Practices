package Seeders

import (
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

var roles =make([]Mod.Role,0)

func RoleSeeder() {
	db := Cfg.DBConnect()

	role1()
	role2()
	for i,_ := range roles {
		db.FirstOrCreate(&roles[i],&Mod.Role{Name:roles[i].Name})
	}
	defer db.Close()
}

func role1() {
	var role = Mod.Role{
		Name: "Admin",
		Status: 1,
	}
	roles = append(roles, role)
}

func role2() {
	var role = Mod.Role{
		Name: "Customer",
		Status: 1,
	}
	roles = append(roles, role)
}
