package main

import (
	"database/sql"
	"fmt"
	//_ "github.com/go-sql-driver/mysql"
)

func DBConnect() (*sql.DB, error) {
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_crud")
	return db, nil
}

func Insert(user User) bool {
	db, _ := DBConnect()
	var insert *sql.Rows
	var err error

	insert, err = db.Query("INSERT INTO users(full_name, email, password) VALUES(?, ?, ?);", user.Name, user.Email, user.Password)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println("\nData Inserterd Successfully.\n")
	defer db.Close()
	defer insert.Close()
	return true
}

func Read(user User) (User, bool) {
	db, _ := DBConnect()
	var results *sql.Rows
	var err error
	var data User

	results, err = db.Query("SELECT * FROM users WHERE email=? and password=?;", user.Email, user.Password)
	if err != nil {
		fmt.Println(err.Error())
		return data, false
	}
	if results.Next() {
		err = results.Scan(&data.ID, &data.Name, &data.Email, &data.Password)
		if err != nil {
			fmt.Println(err.Error())
		}
		return data, true
	} else {
		return data, false
	}

	defer db.Close()
	defer results.Close()
	return data, true
}
