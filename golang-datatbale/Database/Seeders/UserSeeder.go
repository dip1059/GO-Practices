package Seeders

import (
	"golang.org/x/crypto/bcrypt"
	Cfg "gold-store/Config"
	Mod "gold-store/Models"
)

var users []Mod.User

func UserSeeder() {
	db := Cfg.DBConnect()

	user1()
	for i, _ := range users {
		db.Where(&Mod.User{Email: users[i].Email}).FirstOrCreate(&users[i])
	}
	defer db.Close()
}

func user1() {
	hash, _ := bcrypt.GenerateFromPassword([]byte("*123456#"), 10)
	user := Mod.User{
		FirstName:    "Mr.",
		LastName:     "Admin",
		Email:        "admin@xyz.com",
		Phone:        "000000000",
		Password:     string(hash),
		ActiveStatus: 1,
		RoleID:       1,
		CountryCode:  "XXX",
		City:         "XXX",
		Address:      "XXX",
		ZipCode:      1234,
	}
	users = append(users, user)
}
