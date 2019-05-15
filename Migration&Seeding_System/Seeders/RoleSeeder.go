package main

import (
	"database/sql"
	"log"
)

func RoleSeeder() {
	if checkIfSeederExists("RoleSeeder") {
		return
	}
	db, _ := DBConnect()
	var results *sql.Rows

	results, err := db.Query("INSERT INTO roles(role, status) VALUES('Admin', 1), ('User', 1);")
	if err != nil {
		log.Println("RoleSeeder Log1", err.Error())
	}
	defer postProcessing("RoleSeeder", "roles", err)

	defer db.Close()
	defer results.Close()
	return
}
