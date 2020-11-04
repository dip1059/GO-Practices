package Services

import (
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"strconv"
)

type Email_Env struct {
	Host, Port, Username, Password, FromName, FromAdd string
}

var (
	emailEnv Email_Env
)


func init() {
	godotenv.Load()
	emailEnv = Email_Env{
		Host: os.Getenv("MAIL_HOST"),
		Port: os.Getenv("MAIL_PORT"),
		Username: os.Getenv("MAIL_USERNAME"),
		Password: os.Getenv("MAIL_PASSWORD"),
		FromName: os.Getenv("MAIL_FROM_NAME"),
		FromAdd: os.Getenv("MAIL_FROM_ADDRESS"),
	}
}


func SendEmail(to []string, subject, htmlString, filename string) bool {
	mail := gomail.NewMessage()
	mail.SetHeader("From", emailEnv.FromName+"<"+emailEnv.FromAdd+">")
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", htmlString)
	if filename != "" {
		mail.Attach(filename)
	}

	port,_ := strconv.Atoi(emailEnv.Port)
	dialer := gomail.NewDialer(emailEnv.Host, port, emailEnv.Username, emailEnv.Password)

	log.Println(emailEnv)

	for i, _ := range to {
		mail.SetHeader("To", to[i])
		if err := dialer.DialAndSend(mail); err != nil {
			log.Println("Email Sending failed to", to[i])
			log.Println(err.Error())
			return false
		}
	}
	return true
}