package Repositories

import (
	M "PracticeGoland/Models"
	"database/sql"
	"log"
)

func DBConnect() (*sql.DB, error) {
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_crud")
	return db, nil
}

func Register(user M.User) (M.User, bool) {
	db, _ := DBConnect()
	var results *sql.Rows
	var err error

	_, err = db.Query("INSERT INTO users(full_name, email, password, email_verification) VALUES(?, ?, ?, ?);", user.Name, user.Email, user.Password, user.EmailVf)
	if err != nil {
		log.Println("AuthRepo.go Log1", err.Error())
		return user, false
	}
	results, err = db.Query("SELECT * FROM users WHERE email=?;", user.Email)
	if results.Next() {
		err = results.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.ActiveStauts, &user.EmailVf)
		if err != nil {
			log.Println("AuthRepo.go Log2", err.Error())
		}
		return user, true
	} else {
		return user, false
	}

	log.Println("AuthRepo.go Log3 Data Inserterd Successfully.\n")
	defer db.Close()
	defer results.Close()
	return user, true
}

func Read(user M.User) (M.User, bool) {
	db, _ := DBConnect()
	var results *sql.Rows
	var err error

	results, err = db.Query("SELECT * FROM users WHERE email=?;", user.Email)
	if err != nil {
		log.Println("AuthRepo.go Log4", err.Error())
		return user, false
	}
	if results.Next() {
		err = results.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.ActiveStauts, &user.EmailVf)
		if err != nil {
			log.Println("AuthRepo.go Log5", err.Error())
		}
		return user, true
	} else {
		return user, false
	}

	defer db.Close()
	defer results.Close()
	return user, true
}

