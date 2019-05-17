package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"os"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}


type DBEnv struct {
	host string
	port string
	user string
	pass string
	dialect string
	db string
}

func main() {
	godotenv.Load()

	dbEnv := &DBEnv{
		host: os.Getenv("host"),
		port: os.Getenv("port"),
		user: os.Getenv("user"),
		pass: os.Getenv("pass"),
		dialect: os.Getenv("dialect"),
		db: os.Getenv("db"),
	}
	db, err := gorm.Open(dbEnv.dialect, dbEnv.user+":"+dbEnv.pass+"@tcp("+dbEnv.host+":"+dbEnv.port+")/"+dbEnv.db+"?parseTime=true")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()

	db.NewRecord(&Product{Code:"101", Price:2000})
	//db.Create()
}
