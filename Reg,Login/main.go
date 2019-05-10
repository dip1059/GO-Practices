package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	//"github.com/gorilla/sessions"
	"net/http"
	"strconv"
)

var (
	Msg Message
	store = sessions.NewCookieStore([]byte("secret"))
)

type User struct {
	ID int
	Name string
	Email string
	Password string
}

type Message struct {
	Success string
	Fail string
}

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("view/*.html")

	r.GET("/", Welcome)
	r.GET("/register", RegisterGet)
	r.POST("/register", RegisterPost)
	r.GET("/login", LoginGet)
	r.POST("/login", LoginPost)
	r.GET("/home", Home)
	r.GET("/logout", Logout)

	r.Run(":2000")
}

func Welcome(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome.html", nil)
}

func RegisterGet(c *gin.Context)  {
	c.HTML(http.StatusOK, "register.html", Msg)
	Msg.Success = ""
	Msg.Fail = ""
}

func RegisterPost(c *gin.Context) {
	user := User{}
	user.Name = c.PostForm("full_name")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")

	success := Insert(user)
	if success {
		Msg.Success = "Successfully Added."
		c.Redirect(http.StatusFound, "/register")
	} else {
		Msg.Fail = "Some error occurred, please try again."
		c.Redirect(http.StatusInternalServerError, "/register")
	}
}

func LoginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", Msg)
	Msg.Success = ""
	Msg.Fail =""
}

func LoginPost(c *gin.Context) {
	user := User{}
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")
	user, success := Read(user)
	if success {
		session, _ := store.Get(c.Request, "session-cookie")
		session.Values["userId"] = user.ID
		session.Values["userName"] = user.Name
		session.Values["userEmail"] = user.Email
		session.Save(c.Request, c.Writer)
		c.Redirect(http.StatusFound, "/home")
	} else {
		Msg.Fail = "Wrong Credentials."
		c.Redirect(http.StatusFound, "/login")
	}

}

func Home( c *gin.Context) {
	session, _ := store.Get(c.Request, "session-cookie")
	userName := session.Values["userName"]
	userId := session.Values["userId"]
	userEmail := session.Values["userEmail"]
	if userName!=nil {
		//user := user.(User)
		c.String(http.StatusOK, "Welcome "+userName.(string)+" "+strconv.Itoa(userId.(int))+" "+userEmail.(string))
	} else {
		Msg.Fail = "You are not Logged in."
		c.Redirect(http.StatusFound, "/login")
	}
}

func Logout(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-cookie")
	session.Values["userName"] = nil
	session.Values["userId"] = nil
	session.Values["userEmail"] = nil
	session.Save(c.Request, c.Writer)
	c.Redirect(http.StatusFound, "/login")
}
