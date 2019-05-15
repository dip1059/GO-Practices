package main

import (
	"log"
)

func CreateUsersTable() {
	if checkIfMigrationExists("Create_Users_Table") {
		return
	}
	db, _ := DBConnect()
	results, err := db.Query(`create table users (
											id int(10) NOT null AUTO_INCREMENT,
											full_name varchar(191) not null,
											email varchar(191) not null unique,
											password varchar(191) not null,
											activestatus tinyint(4) not null default 0,
											role_id int(10) not null,
											email_verification varchar(191) null,
											remember_token varchar(191) null,
											deleted_at timestamp null,
											created_at timestamp not null DEFAULT CURRENT_TIMESTAMP,
											updated_at timestamp not null DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
											primary key(id),
											foreign key(role_id) references roles(id)
									);`)
	if err != nil {
		log.Println("Create_Users_Table.go Log1", err.Error())
	}
	defer postProcessing("Create_Users_Table", "Create_Table", "users","", err)

	defer db.Close()
	defer results.Close()
	return
}
