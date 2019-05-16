package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var success, fail int
const DbName string = "go_crud"

func DBConnect() (*sql.DB, error) {
	db, _ := sql.Open("mysql", "root:razu@tcp(127.0.0.1:3306)/"+DbName+"?parseTime=true")
	return db, nil
}

func main() {
	if !createMigrationTable() {
		fmt.Println("Migration Failed. Internal Server Error.")
		return
	}
	CreateRolesTable()
	CreateUsersTable()
	CreatePasswordResetsTable()
	AddPhoneColumnOnUsersTable()

	if success+fail == 0{
		fmt.Println("No New Migration, Already Up-to-date.")
	} else {
		fmt.Printf("\nTotal New Migration = %d\nSuccessfully done = %d\nFailed = %d\n", success+fail, success, fail)
	}
}


