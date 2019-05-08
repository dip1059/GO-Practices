package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	Msg Message
	users Users
	mUsers = make(map[int] User)
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
	Msg Message
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
	r.GET("/api/all-user", ApiAllUser)
	r.GET("/use-api-all-user", UseApi)
	r.GET("/view-user/:id", ViewUser)
	r.GET("/edit-user/:id", EditUser)
	r.POST("/update-user", UpdateUser)
	r.GET("/delete-user/:id", DeleteUser)

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
	//var users Users
	var success bool
	users.Users, success = Read()
	if success {
		users.Count = len(users.Users)
		for _, user := range users.Users {
			mUsers[user.ID] = user
			//fmt.Println(mUsers[user.ID])
		}
		c.HTML(http.StatusOK, "all-user.html", users)

		users.Msg.Success = ""
		users.Msg.Fail = ""
	} else {
		users.Msg.Fail = "Some error occurred, please try again."
		c.HTML(http.StatusOK, "all-user.html", users)
	}
}

func ApiAllUser(c *gin.Context) {
	var success bool
	users.Users, success = Read()
	if success {
		users.Count = len(users.Users)
		for _, user := range users.Users {
			mUsers[user.ID] = user
			//fmt.Println(mUsers[user.ID])
		}
		c.JSON(http.StatusOK, users)

		users.Msg.Success = ""
		users.Msg.Fail = ""
	} else {
		users.Msg.Fail = "Some error occurred, please try again."
		c.HTML(http.StatusOK, "all-user.html", users)
	}
}

func ViewUser(c *gin.Context) {
	id,_ := strconv.Atoi(c.Param("id"))
	user := mUsers[id]
	c.HTML(http.StatusOK, "view-user.html", user)
}

func EditUser(c *gin.Context) {
	id,_ := strconv.Atoi(c.Param("id"))
	user := mUsers[id]
	c.HTML(http.StatusOK, "edit-user.html", user)
}

func UpdateUser(c *gin.Context) {
	var user User
	user.ID, _ = strconv.Atoi(c.PostForm("id"))
	user.Name.String = c.PostForm("full_name")
	user.Email.String = c.PostForm("email")
	success := Update(user)
	if success {
		users.Msg.Success = "Successfully Updated."
		c.Redirect(http.StatusFound, "/all-user")
	} else {
		users.Msg.Fail = "Some error occurred, please try again."
		c.Redirect(http.StatusInternalServerError, "/all-user")
	}
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	success := Delete(id)
	if success {
		users.Msg.Success = "Successfully Deleted."
		c.Redirect(http.StatusFound, "/all-user")
	} else {
		users.Msg.Fail = "Some error occurred, please try again."
		c.Redirect(http.StatusInternalServerError, "/all-user")
	}
}

func UseApi(c *gin.Context) {
	resp, err := http.Get("http://localhost:2000/api/all-user")
	if err !=nil {
		fmt.Println(err.Error())
	}
	Bytes, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(Bytes, users)
	fmt.Println(users.Users[0].Name.String)
}
