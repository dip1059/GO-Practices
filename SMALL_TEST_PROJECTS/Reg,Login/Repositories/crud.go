package Repositories

import (
	M "PracticeGoland/Models"
	"database/sql"
	"fmt"
)

func DBConnect() (*sql.DB, error) {
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_crud")
	return db, nil
}

func Register(user M.User) bool {
	db, _ := DBConnect()
	var results *sql.Rows
	var err error

	results, err = db.Query("INSERT INTO users(full_name, email, password) VALUES(?, ?, ?);", user.Name, user.Email, user.Password)
	if err != nil {
		fmt.Println("crud.go Line:21", err.Error())
		return false
	}

	fmt.Println("crud.go Line:25 Data Inserterd Successfully.\n")
	defer db.Close()
	defer results.Close()
	return true
}

func Login(user M.User) (M.User, bool) {
	db, _ := DBConnect()
	var results *sql.Rows
	var err error

	results, err = db.Query("SELECT * FROM users WHERE email=?;", user.Email)
	if err != nil {
		fmt.Println("crud.go Line:38", err.Error())
		return user, false
	}
	if results.Next() {
		err = results.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			fmt.Println("crud.go Line:44", err.Error())
		}
		return user, true
	} else {
		return user, false
	}

	defer db.Close()
	defer results.Close()
	return user, true
}

