package Helpers

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"html/template"
	"log"
	mrand "math/rand"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"
)



func RandomString(n int) string {
	var Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	var lenLetters = len(Letters)
	mrand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = Letters[mrand.Intn(lenLetters)]
	}
	return string(b)
}


func RandomNumber(n int) string {
	var Letters = "0123456789"
	var lenLetters = len(Letters)
	mrand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = Letters[mrand.Intn(lenLetters)]
	}
	return string(b)
}


func ParseTemplate(fileName string, templateData interface{}) (string, error) {
	var str string
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return str, err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, templateData); err != nil {
		return str, err
	}
	str = buffer.String()
	return str, nil
}


func ValidateFile(file multipart.File) (int,bool){

	buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err := file.Read(buff)

	if err != nil {
		log.Println(err.Error())
	}

	filetype := http.DetectContentType(buff)
	log.Println(filetype)

	switch filetype {
	case "image/jpeg", "image/jpg":
		return 1, true

	/*case "image/gif":
	return filetype, true*/

	case "image/png":
		return 1, true

	case "application/pdf": // not image, but application !
		return 2, true
	default:
		return 0, false
	}
}


func NullStringProcess(data sql.NullString) sql.NullString{
	if data.String != "" {
		data.Valid = true
	} else {
		data.Valid = false
	}
	return data
}

func NullFloatProcess(data sql.NullFloat64) sql.NullFloat64{
	if data.Float64 != 0.0 {
		data.Valid = true
	} else {
		data.Valid = false
	}
	return data
}


func SetCookie(hashKey interface{}, blockKey interface{}, value string, name string, age int, c *gin.Context) bool{
	var sc *securecookie.SecureCookie
	if blockKey != nil {
		sc = securecookie.New([]byte(hashKey.(string)), []byte(blockKey.(string)))
	} else {
		sc = securecookie.New([]byte(hashKey.(string)), nil)
	}
	encoded, err := sc.Encode(name, value)
	if err != nil {
		log.Println(name,"Cookie Set Error.",err.Error())
		return false
	}

	cookie := http.Cookie{
		Name:     name,
		Value:    encoded,
		MaxAge:   age,
		Path: "/",
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)
	return true
}


func GetCookie(hashKey interface{}, blockKey interface{}, name string, c *gin.Context) string{
	var sc *securecookie.SecureCookie
	if blockKey != nil {
		sc = securecookie.New([]byte(hashKey.(string)), []byte(blockKey.(string)))
	} else {
		sc = securecookie.New([]byte(hashKey.(string)), nil)
	}

	cookie, err := c.Request.Cookie(name)
	if err != nil {
		log.Println(name,err.Error())
		return ""

	} else {
		var value string
		err = sc.Decode(name, cookie.Value, &value)
		if  err != nil {
			log.Println(err.Error())
			return ""
		} else {
			return value
		}
	}
}


func MakeUrl(str string) string {
	str = LowerCase(str)
	str = strings.Replace(str, " ", "-", -1)
	reg, err := regexp.Compile("[^a-zA-Z0-9-]+")
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	str = reg.ReplaceAllString(str, "")
	return str
}

func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func LowerCase(str string) string {
	var r []rune
	for _, v := range str {
		r = append(r, unicode.ToLower(v))
	}
	return string(r)
}

func UpperCase(str string) string {
	var r []rune
	for _, v := range str {
		r = append(r, unicode.ToUpper(v))
	}
	return string(r)
}

func Recover() {
	if r := recover(); r != nil {
		var ok bool
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("pkg: %v", r)
		}
		log.Println(err.Error())
	}
}