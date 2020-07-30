package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

type Req struct {
	Username string	`form:"username" json:"username"`
	Password string	`form:"password" json:"password"`
	Country string	`form:"country" json:"country"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := os.Setenv("JWT_KEY", "ImagineLennonJohn")
	if err != nil {
		log.Println(err.Error())
	}

	r := gin.Default()

	api := r.Group("/api")
	api.POST("/login", Login)
	api.GET("/profile", Profile)

	_ = r.Run(":2000")

}

func Login(c *gin.Context) {
	/*body, _ := ioutil.ReadAll(c.Request.Body)

	var r Req
	err := json.Unmarshal(body, &r)
	if err !=nil {
		fmt.Println(err.Error())
	}*/
	var r Req
	err := c.ShouldBind(&r)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(r)
	token, _ := GenerateJWT(r.Username)
	log.Println(token)
	c.JSON(http.StatusOK, gin.H{
		"token":token,
	})
}


func Profile(c *gin.Context) {
	if !IsTokenValid(c) {
		c.JSON(http.StatusOK, gin.H{
			"success":false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":true,
		"message":"Welcome",
	})
}


func GenerateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = username
	claims["iat"] = time.Now().Unix()
	//claims["exp"] = time.Now().Add(time.Minute * 5).Unix()
	key := os.Getenv("JWT_KEY")
	tokenString, err := token.SignedString([]byte(key))

	if err != nil {
		log.Println("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}


func IsTokenValid(c *gin.Context) bool {
	if c.Request.Header["Token"] != nil {

		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(c.Request.Header["Token"][0], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing mehod error")
			}
			return []byte(os.Getenv("JWT_KEY")), nil
		})
		if err != nil {
			log.Println(err.Error())
			return false
		}

		log.Println(claims["client"])

		if token.Valid {
			return true
		}
	} else {
		return false
	}
	return false
}