package Config

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"time"
)

func WriteLogFile() {
	f, err := os.OpenFile("Storage/Logs/route-logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err.Error())
	}
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	//log.SetOutput(gin.DefaultWriter)
}


func WritePanicLogFile() {
	//var err error
	/*if runtime.GOOS == "windows" {
		fileName = "log.txt"
	} else {
		str = time.Now().String()+"_log.txt"
		fileName = strings.Replace(str, " ", "_", -1)
	}*/
	/*os.Stderr*/f, err := os.OpenFile("Storage/Logs/panic-logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err.Error())
	}
	//syscall.Dup2(int(os.Stderr.Fd()), 2)

	gin.DefaultErrorWriter = io.MultiWriter(/*os.Stderr*/f, os.Stdout)
	/*log.SetOutput(gin.DefaultErrorWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)*/
}


func DailyLogFile(){
	count := 1
	var now, tomor time.Time
	var err error
	var nowT, tom, fileName string
	var f *os.File
	var duration time.Duration

	for {
		if count == 1 {
			nowT = time.Now().Format("2006-Jan-02 15:04:05")
			now, err = time.Parse("2006-Jan-02 15:04:05", nowT)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			tom = now.AddDate(0, 0, 1).Format("2006-Jan-02") + " 00:00:01"
			tomor, err = time.Parse("2006-Jan-02 15:04:05", tom)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			duration = tomor.Sub(now)
			log.Println(duration)

		} else if count > 1 {
			now = time.Now()
			log.Println("next file will be created at:",now.Add(time.Hour * 24).Format("2006-Jan-02 15:04:05"))
		}

		fileName = now.Format("2006-Jan-02")
		f, err = os.OpenFile("Storage/DailyLogs/"+fileName+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println(err.Error())
		}
		//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		//log.SetOutput(gin.DefaultWriter)
		log.SetOutput(io.MultiWriter(f, os.Stdout))
		log.Println("New file "+fileName+".log ------ Count:", count)

		if count == 1 {
			time.Sleep(duration)
		} else if count > 1 {
			time.Sleep(time.Hour * 24)
		}
		count++
	}
}
