package main

import (
	"database/sql"
	"fmt"
	//_ "github.com/go-sql-driver/mysql"
)

/*func sd() {
	var i = 0

	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_crud")

	create, err := db.Query(`CREATE TABLE users (
										id int(10) not null AUTO_INCREMENT,
										full_name varchar(40) not null,
										email varchar(40) UNIQUE,
										PRIMARY KEY(id)
									);`)

	if err != nil {
		panic(err.Error())
	}

	if create != nil {
		fmt.Println("Table Created Successfully.\n")
	}



	results, err := db.Query("SELECT * FROM users where full_name like ?", "%")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("ID", " FULL_NAME", "            EMAIL")
	var user = make([]User, 0)
	for i = 0; results.Next(); i++ {
		var data User
		err = results.Scan(&data.ID, &data.Name, &data.Email)
		if err != nil {
			panic(err.Error())
		}
		user = append(user, data)
		fmt.Print(user[i].ID, "   ", user[i].Name.String, "      ")
		if user[i].Email.String == "" {
			fmt.Println("NULL")
		} else {
			fmt.Println(user[i].Email.String)
		}
	}

	if len(user) <= 0 {
		fmt.Println("No Data Found.\n")
	}
	fmt.Println("Total:", len(user), "\n")

	update, err := db.Query("UPDATE users SET email='dipankarsaha1059@gmail.com' WHERE id=2")
	if err != nil {
		panic(err.Error())
	}

	if update != nil {
		fmt.Println("Data Updated Successfully.\n")
	}

	results, err = db.Query("SELECT * FROM users where id=2")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("ID", " FULL_NAME", "       EMAIL")
	//var user= make([]User, 0)
	for i = 0; results.Next(); i++ {
		var data User
		err = results.Scan(&data.ID, &data.Name, &data.Email)
		if err != nil {
			panic(err.Error())
		}
		//user = append(user, data)
		fmt.Print(data.ID, "   ", data.Name.String, "        ")
		if data.Email.String == "" {
			fmt.Println("NULL")
		} else {
			fmt.Println(data.Email.String)
		}
	}
	//fmt.Println("Total:", len(user))
	if i <= 0 {
		fmt.Println("No Data Found.\n")
	}

	defer db.Close()
	defer create.Close()
	defer results.Close()
	defer update.Close()
}*/

func DBConnect() (*sql.DB, error){
	db, _ := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_crud")
	return db, nil
}

func Insert(user User) bool {
	db, _ := DBConnect()
	var insert *sql.Rows
	var err error

	if user.Email.String == "" {
		insert, err = db.Query("INSERT INTO users(full_name, email) VALUES(?, NULL);", user.Name.String)
	} else {
		insert, err = db.Query("INSERT INTO users(full_name, email) VALUES(?, ?);", user.Name.String, user.Email.String)
	}

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	fmt.Println("Data Inserterd Successfully.\n")
	defer db.Close()
	defer insert.Close()
	return true
}

func Read() ([]User, bool) {
	db, _ := DBConnect()
	var users = make([]User, 0)
	var results *sql.Rows
	var err error
	var i int

	results, err = db.Query("SELECT * FROM users;")
	if err != nil {
		fmt.Println(err.Error())
		return users, false
	}

	for i = 0; results.Next(); i++ {
		//important part
		var data User
		//
		err = results.Scan(&data.ID, &data.Name, &data.Email)
		if err != nil {
			fmt.Println(err.Error())
		}
		users = append(users, data)
	}
	fmt.Println("Total Rows Found:", len(users))

	defer db.Close()
	defer results.Close()
	return users, true
}
