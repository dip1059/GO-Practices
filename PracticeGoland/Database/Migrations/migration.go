package main

import (
	"fmt"
	"log"
)

func createMigrationTable() bool {
	if checkIfTableExists("migrations") {
		return true
	}
	db, _ := DBConnect()
	results, err := db.Query(`create table migrations (
											id int(10) NOT null AUTO_INCREMENT,
											migration varchar(191) not null unique,
											created_at timestamp not null DEFAULT CURRENT_TIMESTAMP,
											updated_at timestamp not null DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
											primary key(id)
									);`)
	if err != nil {
		log.Println("migration.go log1", err.Error())
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
		log.Println("migration.go log2", err.Error())
		return false
	}
	if results.Next() {
		return true
	}
	defer db.Close()
	defer results.Close()
	return false
}

func checkIfMigrationExists(migrationName string) bool {
	db, _ := DBConnect()
	results, err := db.Query("select id from migrations where migration=?", migrationName)
	if err != nil {
		log.Println("migration.go log3 error during checking on "+migrationName+" migration\n", err.Error())
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

func postProcessing(migrationName, migrationType, tableName, columnName string, err error) {
	recover()
	if err != nil {
		fail++
		fmt.Println("'" + migrationName + "' Migration Failed.")
	} else {
		if insertMigration(migrationName) {
			success++
			fmt.Println("'" + migrationName + "' Migrated Successfully.")
		} else {
			if migrationType == "Create_Table" {
				dropTable(tableName)
			} else if migrationType == "Add_Column" {
				dropColumn(tableName, columnName)
			}
			fail++
		}
	}
}

func insertMigration(migrationName string) bool {
	db, _ := DBConnect()
	results, err := db.Query("Insert into migrations(migration) values(?);", migrationName)
	if err != nil {
		log.Println("migration.go log4 error during insert "+migrationName+"\n", err.Error())
		return false
	}
	defer db.Close()
	defer results.Close()
	return true
}

func dropTable(tableName string) {
	db, _ := DBConnect()
	results, err := db.Query("Drop table " + tableName + ";")
	if err != nil {
		log.Println("migration.go log6", err.Error())
		fmt.Println("Rollback Failed. Drop Table '" + tableName + "' Manually From The DB.")
		return
	}
	defer db.Close()
	defer results.Close()
	return
}

func dropColumn(tableName string, columnName string) {
	db, _ := DBConnect()
	results, err := db.Query("alter table " + tableName + " drop column " + columnName + ";")
	if err != nil {
		log.Println("migration.go log7", err.Error())
		fmt.Println("Rollback Failed. Drop Column '" + columnName + "' From Table '" + tableName + "' Manually From The DB.")
		return
	}
	defer db.Close()
	defer results.Close()
	return
}
