package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Data struct {
	Result int
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("html/*.html")
	r.GET("/", LoadForm)
	r.POST("/fact", GetFact)
	r.Run(":3003")
	//ans := Fact(10)
	//fmt.Println(ans)
}

func LoadForm(c *gin.Context) {
	c.HTML(http.StatusOK, "form.html", nil)
	return
}

func GetFact(c *gin.Context) {
	number, err := strconv.Atoi(c.PostForm("number"))
	if err != nil {
		c.String(http.StatusOK, "Number Please")
		return
	}else {
		data := Data{}
		data.Result= Fact(number)
		ans := data.Result
		fmt.Print(data.Result)
		res := strconv.Itoa(ans)
		c.String(http.StatusOK, res)
	}
}

func Fact(num int) int {
	res := 1
	for i := 2; i <= num; i++ {
		res *= i
	}
	return res
}