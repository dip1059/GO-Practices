package main

import (
	"fmt"
	"log"
)

func createSeederTable() bool {
	if checkIfTableExists("seeders") {
		return true
	}
	db, _ := DBConnect()
	results, err := db.Query(`create table seeders (
											id int(10) NOT null AUTO_INCREMENT,
											seeder varchar(191) not null unique,
											created_at timestamp not null DEFAULT CURRENT_TIMESTAMP,
											updated_at timestamp not null DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
											primary key(id)
									);`)
	if err != nil {
		log.Println("seeder.go log1", err.Error())
		return false
	}
	defer db.Close()
	defer results.Close()
	return true
}

func checkIfTableExists(tableName string) bool {
	db, _ := DBConnect()
	results, err := db.Query(`select TABLE_NAME from information_schema.tables 
									where table_name=? and TABLE_TYPE='BASE TABLE' 
									and TABLE_SCHEMA=?`, tableName, DbName)
	if err != nil {
		log.Println("seeder.go log2", err.Error())
		return false
	}
	if results.Next() {
		return true
	}
	defer db.Close()
	defer results.Close()
	return false
}

func checkIfSeederExists(seederName string) bool {
	db, _ := DBConnect()
	results, err := db.Query("select id from seeders where seeder=?", seederName)
	if err != nil {
		log.Println("seeder.go log3 error during checking on "+seederName+" seeder\n", err.Error())
		return true
	}
	if results.Next() {
		return true
	} else {
		return false
	}

	defer db.Close()
	defer results.Close()
	return false
}

func postProcessing(seederName, tableName string, err error) {
	recover()
	if err != nil {
		fail++
		fmt.Println("'" + seederName + "' Seeding Failed.")
	} else {
		if insertSeeder(seederName) {
			success++
			fmt.Println("'" + seederName + "' Seeded Successfully.")
		} else {
			fmt.Println("Delete Rows From Table '" + tableName + "' Manually From The DB, Or")
			fmt.Println("Insert The Seeder '" + seederName + "' Into Seeder Table Manually, If")
			fmt.Println("All Seeding Data Are Inserted Into Table '" + tableName + "' Successfully.")
			fail++
		}
	}
}

func insertSeeder(seederName string) bool {
	db, _ := DBConnect()
	results, err := db.Query("Insert into seeders(seeder) values(?);", seederName)
	if err != nil {
		log.Println("seeder.go log4 error during insert "+seederName+"\n", err.Error())
		return false
	}
	defer db.Close()
	defer results.Close()
	return true
}
