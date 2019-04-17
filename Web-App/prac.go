package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Data struct {
	Name string
	Sex  string
	Age  int
}

func main() {
	r := gin.Default()

	r.GET("/", DisplayString)
	r.GET("/json/:name/:sex/:age", JsonData)

	r.LoadHTMLGlob("html/*.html")
	r.GET("/page", LoadPage)
	r.GET("/page2", LoadPage2)
	r.GET("/form", LoadForm)
	r.POST("/welcome", Welcome)

	r.Run(":2000")
}

func DisplayString(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

func JsonData(c *gin.Context) {
	name := c.Param("name")
	sex := c.Param("sex")
	age, _ := strconv.Atoi(c.Param("age"))
	var data = []*Data{
		{
			Name: name,
			Sex:  sex,
			Age:  age,
		},
		{
			Name: name,
			Sex:  sex,
			Age:  age + 2,
		},
	}
	c.JSON(http.StatusOK, &data)
}

func LoadPage(c *gin.Context) {
	c.HTML(http.StatusOK, "page.html", nil)
}

func LoadPage2(c *gin.Context) {
	data := Data{}
	data.Name = "Dipankar Saha"
	data.Age = 25
	c.HTML(http.StatusOK, "page2.html", data)
}

func LoadForm(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", nil)
}

func Welcome(c *gin.Context) {

	name := c.PostForm("name")
	//type casting string2int
	age, err := strconv.Atoi(c.PostForm("age"))
	//
	if err != nil {
		c.String(http.StatusOK, "Only Integer value.")
		return
	}

	data := Data{}
	data.Name = name
	data.Age = age
	c.HTML(http.StatusOK, "page2.html", data)
}
