package main

import (
	"log"
)

func AddPhoneColumnOnUsersTable() {
	if checkIfMigrationExists("Add_Phone_Column_On_Users_Table") {
		return
	}
	db, _ := DBConnect()
	results, err := db.Query(`alter table users add phone varchar(30) null unique after email;`)
	if err != nil {
		log.Println("Add_Phone_Column_On_Users_Table.go Log1", err.Error())
	}
	defer postProcessing("Add_Phone_Column_On_Users_Table", "Add_Column", "users","phone", err)

	defer db.Close()
	defer results.Close()
	return
}
