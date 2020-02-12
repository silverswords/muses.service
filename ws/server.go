package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type roomManager interface {
	Join(room string, connection Conn)
	Leave(room string, connection Conn)

	// may send msg or webRTC request
	Send(room, event string, args ...interface{})

	// list members
	Members(room string) int
	Rooms(connection Conn) []string
}

type client interface {
	// session id
	ID() string
	Close() error

	Send(msg string, v ...interface{})
	// receive msg or webRTC request
	Emit(msg string, v ...interface{})

	Join(room string)
	Leave(room string)
}

type server struct {
	rooms map[string]map[string]client
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ping(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func main() {
	bindAddress := "localhost:2303"
	r := gin.Default()
	r.GET("/ping", ping)
	r.Run(bindAddress)
}
