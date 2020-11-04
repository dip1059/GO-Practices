package Config

import (
	G "gold-store/Globals"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
)

/*func init() {
	godotenv.Load()
	G.DBEnv = G.DB_ENV{
		Host:os.Getenv("DB_HOST"),
		Port:os.Getenv("DB_PORT"),
		Dialect:os.Getenv("DB_DIALECT"),
		Username:os.Getenv("DB_USERNAME"),
		Password:os.Getenv("DB_PASSWORD"),
		DBname:os.Getenv("DB_NAME"),
	}
}*/

func DBConnect() *gorm.DB{
	godotenv.Load()
	G.DBEnv = G.DB_ENV{
		Host:os.Getenv("DB_HOST"),
		Port:os.Getenv("DB_PORT"),
		Dialect:os.Getenv("DB_DIALECT"),
		Username:os.Getenv("DB_USERNAME"),
		Password:os.Getenv("DB_PASSWORD"),
		DBname:os.Getenv("DB_NAME"),
	}

	db, err := gorm.Open(G.DBEnv.Dialect, G.DBEnv.Username+":"+G.DBEnv.Password+"@tcp("+G.DBEnv.Host+":"+G.DBEnv.Port+")/"+G.DBEnv.DBname+"?parseTime=true")
	if err !=nil {
		log.Println("log", err.Error())
		return db
	}
	db.SetLogger(log.New(gin.DefaultWriter, "\r\n", 0))
	return db
}

/*func DBConnect2() *gorm.DB{
	godotenv.Load()
	DBEnv := G.DB_ENV{
		Host:os.Getenv("DB_HOST2"),
		Port:os.Getenv("DB_PORT2"),
		Dialect:os.Getenv("DB_DIALECT2"),
		Username:os.Getenv("DB_USERNAME2"),
		Password:os.Getenv("DB_PASSWORD2"),
		DBname:os.Getenv("DB_NAME2"),
	}

	var err error
	db, err := gorm.Open(DBEnv.Dialect, DBEnv.Username+":"+DBEnv.Password+"@tcp("+DBEnv.Host+":"+DBEnv.Port+")/"+DBEnv.DBname+"?parseTime=true")
	if err !=nil {
		log.Println("log", err.Error())
		return db
	}
	db.SetLogger(log.New(gin.DefaultWriter, "\r\n", 0))
	return db
}*/