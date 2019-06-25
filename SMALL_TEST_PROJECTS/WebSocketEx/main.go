package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("*.html")

	r.GET("/", Welcome)
	r.GET("/welcome2", Welcome2)
	r.GET("/websock", WS)

	r.Run(":2000")
}

func Welcome(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome.html", nil)
}


func Welcome2(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome2.html", nil)
}


func WS(c *gin.Context) {
	wshandler(c.Request, c.Writer)
}


//
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(r *http.Request, w http.ResponseWriter) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to set websocket upgrade: %+v", err.Error())
		return
	}

	i := 0
	for {
		i++
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		//log.Println(i)
		conn.WriteMessage(t, msg)
		if i == 5 {
			conn.Close()
		}
	}
}
//