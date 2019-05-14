package Services

import (
	"gopkg.in/gomail.v2"
	"log"
)

func SendEmail(from, to, subject, htmlString string) bool {
	mail := gomail.NewMessage()
	mail.SetHeader("From", from)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", htmlString)
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "razibbiswas777@gmail.com", "random.09876")
	if err := dialer.DialAndSend(mail); err != nil {
		log.Println("EmailService.go Log1", err.Error())
		return false
	}
	return true
}