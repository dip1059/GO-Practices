package main

import (
	Cfg "gold-store/Config"
	Mig "gold-store/Database/Migrations"
	Seed "gold-store/Database/Seeders"
	"gold-store/Routes"
	_ "github.com/chekun/golaravelsession"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	Mig.Migrate()
	Seed.Seed()
	Cfg.Config()
	Routes.Routes()
}
