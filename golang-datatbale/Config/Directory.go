package Config

import (
	"log"
	"os"
)

func DirectoryConfig() {
	CreateDirectories()
}

func CreateDirectories() {

	err := os.MkdirAll("./Storage/Session", 0777)
	if err != nil {
		log.Println(err.Error())
	}

	err = os.MkdirAll("./Storage/Images", 0777)
	if err != nil {
		log.Println(err.Error())
	}

	err = os.MkdirAll("./Storage/Logs", 0777)
	if err != nil {
		log.Println(err.Error())
	}

	//err = os.MkdirAll("./Storage/PanicLogs", 0777)
	//if err != nil {
	//	log.Println(err.Error())
	//}

	err = os.MkdirAll("./Storage/Temp", 0777)
	if err != nil {
		log.Println(err.Error())
	}

	err = os.MkdirAll("./Storage/DailyLogs", 0777)
	if err != nil {
		log.Println(err.Error())
	}
}