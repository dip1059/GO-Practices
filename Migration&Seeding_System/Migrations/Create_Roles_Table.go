package main

import (
	"log"
)

func CreateRolesTable() {
	if checkIfMigrationExists("Create_Roles_Table") {
		return
	}
	db, _ := DBConnect()
	results, err := db.Query(`create table roles (
											id int(10) NOT null AUTO_INCREMENT,
											role varchar(191) not null unique,
											status tinyint(4) not null default 0,
											deleted_at timestamp null ,
											created_at timestamp not null DEFAULT CURRENT_TIMESTAMP,
											updated_at timestamp not null DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
											primary key(id)
									);`)
	if err != nil {
		log.Println("Create_Roles_Table.go Log1", err.Error())
	}
	defer postProcessing("Create_Roles_Table", "Create_Table", "roles", "", err)

	defer db.Close()
	defer results.Close()
	return
}
