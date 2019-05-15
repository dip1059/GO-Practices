package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func UserSeeder() {
	if !checkIfSeederExists("UserSeeder") {
		db, _ := DBConnect()
		var results *sql.Rows
		cost := bcrypt.DefaultCost
		hash, err := bcrypt.GenerateFromPassword([]byte("09876"), cost)
		if err != nil {
			log.Println("UserSeeder Log1", err.Error())
		}
		results, err = db.Query(`INSERT INTO users(full_name, email, password, role_id, activestatus) 
									VALUES('Mr. Admin', 'admin@xyz.com', ?, 1, 1);`, string(hash))
		if err != nil {
			log.Println("UserSeeder Log2", err.Error())
		}
		defer postProcessing("UserSeeder", "users", err)

		defer db.Close()
		defer results.Close()
		return
	}
}
