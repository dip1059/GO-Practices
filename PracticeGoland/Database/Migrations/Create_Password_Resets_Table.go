package main

import (
	"log"
)

func CreatePasswordResetsTable() {
	if checkIfMigrationExists("Create_Password_Resets_Table") {
		return
	}
	//dropTableIfExists("password_resets")
	db, _ := DBConnect()
	results, err := db.Query(`create table password_resets (
											id int(10) NOT null AUTO_INCREMENT,
											email varchar(191) not null,
											token varchar(191) null,
											status tinyint(4) default 0,
											deleted_at timestamp null,
											created_at timestamp not null DEFAULT CURRENT_TIMESTAMP,
											updated_at timestamp not null DEFAULT CURRENT_TIMESTAMP on UPDATE CURRENT_TIMESTAMP,
											primary key(id),
											foreign key(email) references users(email)
									);`)
	if err != nil {
		log.Println("Create_Password_Resets_Table.go Log1", err.Error())
	}
	defer postProcessing("Create_Password_Resets_Table","Create_Table", "password_resets","", err)

	defer db.Close()
	defer results.Close()
	return
}

/*func dropTableIfExists(tableName string) {
	db, _ := DBConnect()
	results, err := db.Query(`select TABLE_NAME from information_schema.tables
									where table_name=? and TABLE_TYPE='BASE TABLE'
									and TABLE_SCHEMA=?`, tableName, DbName)
	if err != nil {
		log.Println("Create_Password_Resets_Table.go log2", err.Error())
		return
	}
	if results.Next() {
		dropTable(tableName)
	}
	defer db.Close()
	defer results.Close()
	return
}*/
