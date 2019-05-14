package main

import (
	Cont "PracticeGoland/Controllers"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("View/*.html")
	r.GET("/", Cont.Welcome)
	r.GET("/register", Cont.RegisterGet)
	r.POST("/register", Cont.RegisterPost)
	r.GET("resend-email-verification", Cont.ResendEmailVf)
	r.GET("/login", Cont.LoginGet)
	r.POST("/login", Cont.LoginPost)
	r.GET("/home", Cont.Home)
	r.GET("/dashboard", Cont.Dashboard)
	r.GET("/logout", Cont.Logout)

	r.Run(":2000")
}

