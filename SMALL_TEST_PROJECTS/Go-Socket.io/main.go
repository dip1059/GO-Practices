package main

import (
	"github.com/gin-gonic/gin"
	//socket "github.com/googollee/go-socket.io"
	socket "github.com/nkovacs/go-socket.io"
	"log"
	"net/http"
)

var (
	socketServer, err = socket.NewServer(nil)
    sock socket.Socket
)

func main() {

	if err != nil {
		log.Println(err.Error())
	}
	
	r := gin.Default()
	r.LoadHTMLGlob("*.html")


	//wshandler()
	wshandler2()

	//go socketServer.Serve()
	//defer socketServer.Close()

	r.GET("/", Welcome)
	r.GET("/2", Welcome2)
	r.GET("/broadcast", Broadcast)
	r.GET("/websock/*any", gin.WrapH(socketServer))
	r.POST("/websock/*any", gin.WrapH(socketServer))

	err = r.Run(":3000")
	if err != nil {
		log.Println(err.Error())
	}
}

func Welcome(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome.html", nil)
}

func Welcome2(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome2.html", nil)
}

func Broadcast(c *gin.Context) {
	room := c.Query("room")
	event := c.Query("event")
	//socketServer.BroadcastToRoom("", room, event, "Hello all, it's from room: "+room+" and  event: "+event)
	socketServer.BroadcastTo(room,  event, "Hello all, it's from room: "+room+" and  event: "+event)
	/*err := sock.BroadcastTo(room, event, "Hello all, it's from room: "+room+" and  event: "+event)
	if err != nil {
		log.Println(err.Error())
	}*/

}

func wshandler2() {

	err := socketServer.On("connection", func(s socket.Socket) error {

		socketServer.SetAllowRequest(func(request *http.Request) error {
			log.Println(request.Cookie("lang"))
			log.Println(request.Cookies())
			return nil
		})

		//sock = s
		err = s.Join("chatroom")
		if err != nil {
			log.Println(err.Error())
		}
		err = s.Emit("connection", "Hello good people")
		if err != nil {
			log.Println(err.Error())
		}
		return nil
	})
	if err != nil {
		log.Println(err.Error())
	}



	err = socketServer.On("notice1", func(s socket.Socket, msg string) {
		log.Println("notice1:", msg)
		err = s.Emit("notice1", "You sent: "+msg)
		if err != nil {
			log.Println(err.Error())
		}
	})
	if err != nil {
		log.Println(err.Error())
	}



	err = socketServer.On("notice2", func(s socket.Socket, msg string) string {
		log.Println("notice2:", msg)
		err = s.Emit("notice2", "You sent: "+msg)
		if err != nil {
			log.Println(err.Error())
		}
		/*err := s.BroadcastTo("chatroom", "notice2", "Hello all, it's from room: chatroom and  event: ")
		if err != nil {
			log.Println(err.Error())
		}*/
		return "recv " + msg
	})
	if err != nil {
		log.Println(err.Error())
	}



}

/*func wshandler() {
	socketServer.OnConnect("/", func(s socket.Conn) error {
		s.SetContext("")

		log.Println("connected:", s.ID())
		s.Join("chatroom")
		s.Emit("connection", "How r u")
		return nil
	})

	socketServer.OnEvent("/", "notice1", func(s socket.Conn, msg string) {
		log.Println("notice1:", msg)
		s.Emit("notice1", "You sent: "+msg)
	})
	socketServer.OnEvent("/", "notice2", func(s socket.Conn, msg string) string {
		s.SetContext(msg)
		log.Println("notice2:", msg)
		s.Emit("notice2", "You sent: "+msg)
		return "recv " + msg
	})
	socketServer.OnEvent("/", "bye", func(s socket.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	socketServer.OnError("/", func(s socket.Conn, e error) {
		log.Println("meet error:", e)
	})
	socketServer.OnDisconnect("/", func(s socket.Conn, reason string) {
		log.Println("closed", reason)
	})
}*/