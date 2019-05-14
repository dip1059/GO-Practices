package Repositories

import (
	G "PracticeGoland/Globals"
	M "PracticeGoland/Models"
	"database/sql"
	"log"
)

func DBConnect() (*sql.DB, error) {
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/"+G.DbName+"?parseTime=true")
	return db, nil
}

func ReadWithEmail(user M.User) (M.User, bool) {
	db, _ := DBConnect()
	var results *sql.Rows
	var err error
	results, err = db.Query("SELECT * FROM users WHERE email=?;", user.Email)
	if results.Next() {
		err = results.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Password, &user.ActiveStauts, &user.Role, &user.EmailVf, &user.RememberToken, &user.Deleted_At, &user.Created_At, &user.Updated_At)
		if err != nil {
			log.Println("AuthRepo.go Log1", err.Error())
		}
		return user, true
	} else {
		return user, false
	}

	defer db.Close()
	defer results.Close()
	return user, true
}

func Register(user M.User) (M.User, bool) {
	db, _ := DBConnect()
	var results *sql.Rows
	var err error
	var success bool

	_, err = db.Query("INSERT INTO users(full_name, email, password, email_verification) VALUES(?, ?, ?, ?);", user.Name, user.Email, user.Password, user.EmailVf)
	if err != nil {
		log.Println("AuthRepo.go Log2", err.Error())
		return user, false
	}
	user, success = ReadWithEmail(user)
	if success {
		return user, true
	} else {
		return user, false
	}

	log.Println("AuthRepo.go Log3 Data Inserterd Successfully.\n")
	defer db.Close()
	defer results.Close()
	return user, true
}

func Login(user M.User) (M.User, bool) {
	var success bool
	user, success = ReadWithEmail(user)
	if success {
		return user, true
	} else {
		return user, false
	}
}

func ActivateAccount(user M.User) (M.User, bool) {
	db, _ := DBConnect()
	var results *sql.Rows
	var success bool

	results, err := db.Query("SELECT * FROM users WHERE email=? and email_verification=?;", user.Email, user.EmailVf.String)

	if results.Next() {
		results, err = db.Query("UPDATE users SET activestatus=1, email_verification=NULL WHERE email=? and email_verification=?;", user.Email, user.EmailVf.String)

		if err != nil {
			log.Println("AuthRepo.go Log4", err.Error())
			return user, false
		}

		user, success = ReadWithEmail(user)
		if success {
			return user, true
		} else {
			return user, false
		}
	} else {
		return user, false
	}

	defer db.Close()
	defer results.Close()
	return user, true
}
