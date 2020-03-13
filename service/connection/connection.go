package connection

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"muses.service/middleware"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client - connection
type Client struct {
	id      string
	manager *Manager
	conn    *websocket.Conn
	send    chan []byte
}

// Manager - manager connections
type Manager struct {
	// Registered clients.
	Connections map[string]*Client

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewConnectionManager -
func NewConnectionManager() *Manager {
	return &Manager{
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		Connections: make(map[string]*Client),
	}
}

// Run -
func (manager *Manager) Run() {
	for {
		select {
		case client := <-manager.register:
			manager.Connections[client.id] = client
		case client := <-manager.unregister:
			if _, ok := manager.Connections[client.id]; ok {
				delete(manager.Connections, client.id)
				close(client.send)
			}
		}
	}
}

// UpGraderWs - upgrade ws
func (manager *Manager) UpGraderWs(ctx *gin.Context) {
	var upgrader = websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
	}

	hToken := ctx.GetHeader("Authorization")
	if len(hToken) < len("Bearer ") {
		ctx.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": "header Authorization has not Bearer token"})
	}
	token := strings.TrimSpace(hToken[len("Bearer ") : len(hToken)-1])

	usrID, _ := middleware.JwtParseUser(token)
	fmt.Println(usrID)

	client := &Client{id: usrID, manager: manager, conn: conn, send: make(chan []byte, 256)}
	client.manager.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
}

// WritePump - receive msg
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SubMsg -
func (manager *Manager) SubMsg(users []string, msg string) {
	for _, v := range users {
		client, ok := manager.Connections[v]
		if ok {
			client.send <- []byte(msg)
		}
	}
}
