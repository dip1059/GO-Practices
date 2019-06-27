package main

import (
	"fmt"
	"github.com/bykovme/gotrans"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
)

func main() {
	err := gotrans.InitLocales("Langs")
	if err != nil {
		log.Println(err.Error())
	}

	r := gin.Default()

	r.LoadHTMLGlob("*.html")

	r.GET("/", Home)
	r.GET("/html", Html)

	r.Run(":2000")
}

func Home(c *gin.Context) {
	//lang := "en"
	lang := "bd"
	w := c.Writer

	//lang := gotrans.DetectLanguage(r.Header.Get("Accept-Language"))

	fmt.Fprintf(w, "<html><head><title> %s </title></head><body>", gotrans.Tr(lang, "Hello World"))
	fmt.Fprintf(w, "<h2> %s </h2>", gotrans.Tr(lang, "Hello World"))
	githubLink := "https://github.com/bykovme/gotrans"
	link := fmt.Sprintf(`<a href="%s">%s</a>`, githubLink, githubLink)
	fmt.Fprintf(w, gotrans.Tr(lang, "Find more information about the project on the website"), link)
	fmt.Fprint(w, "</body></html>")
}


func Html(c *gin.Context) {
	//lang := "en"
	lang := "bd"
	title := gotrans.Tr(lang, "Hello World")
	head := gotrans.Tr(lang, "Hello World")
	githubLink := "https://github.com/bykovme/gotrans"
	link := fmt.Sprintf(`<a href="%s">%s</a>`, githubLink, githubLink)
	var body template.HTML
	body = template.HTML(fmt.Sprintf(gotrans.Tr(lang, "Find more information about the project on the website"),link))
	c.HTML(http.StatusOK, "home.html", map[string]interface{}{
		"Title":title, "Head":head, "Body":body})
}