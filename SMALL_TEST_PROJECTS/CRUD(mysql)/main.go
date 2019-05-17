package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID    int
	Name  sql.NullString
	Email sql.NullString
}

func main() {
	var i = 0

	db, _ := sql.Open("mysql", "root:razu@tcp(127.0.0.1:3306)/go_crud")

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

	insert, err := db.Query("INSERT INTO go_crud.users(full_name, email) SELECT full_name, email FROM medicare2.users;")

	if err != nil {
		panic(err.Error())
	}

	if insert != nil {
		fmt.Println("Data Inserterd Successfully.\n")
	}

	results, err := db.Query("SELECT * FROM users where full_name like ?", "%")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("ID", " FULL_NAME", "       EMAIL")
	var user = make([]User, 0)
	for i = 0; results.Next(); i++ {
		var data User
		err = results.Scan(&data.ID, &data.Name, &data.Email)
		if err != nil {
			panic(err.Error())
		}
		user = append(user, data)
		fmt.Print(user[i].ID, "   ", user[i].Name.String, "        ")
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
	defer insert.Close()
	defer results.Close()
	defer update.Close()
}
