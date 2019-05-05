package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

var (
	Msg Message
)

type Users struct {
	Users []User
	Count int
	Msg Message
}

type User struct {
	ID int
	Name sql.NullString
	Email sql.NullString
}

type Message struct {
	Success string
	Fail string
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("view/*.html")
	r.GET("/", Welcome)
	r.GET("/add-user", AddUserGet)
	r.POST("/add-user", AddUserPost)
	r.GET("/all-user", AllUser)
	r.GET("/view-user/:id", ViewUser)

	r.Run(":2000")
}

func Welcome(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome.html", nil)
}

func AddUserGet(c *gin.Context)  {
	c.HTML(http.StatusOK, "add-user.html", Msg)
	Msg.Success = ""
	Msg.Fail = ""
}

func AddUserPost(c *gin.Context) {
	user := User{}

	user.Name.String = c.PostForm("full_name")
	user.Email.String = c.PostForm("email")

	success := Insert(user)
	if success {
		Msg.Success = "Successfully Added."
		c.Redirect(http.StatusFound, "/add-user")
	} else {
		Msg.Fail = "Some error occurred, please try again."
		c.Redirect(http.StatusInternalServerError, "/add-user")
	}
}

func AllUser(c *gin.Context) {
	var users Users
	var success bool
	users.Users, success = Read()
	if success {
		users.Count = len(users.Users)
		c.HTML(http.StatusOK, "all-user.html", users)
		users.Msg.Fail = ""
	} else {
		users.Msg.Fail = "Some error occurred, please try again."
		c.HTML(http.StatusOK, "all-user.html", users)
	}
}

func ViewUser(c *gin.Context) {
	fmt.Println("Yo")
}