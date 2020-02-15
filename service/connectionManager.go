package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"

	"github.com/go-redis/redis/v7"
)

type Conn interface{
		// ID returns session id
		ID() string
		Close() error
		URL() url.URL
		LocalAddr() net.Addr
		RemoteAddr() net.Addr
		RemoteHeader() http.Header
	
		// Context of this connection. You can save one context for one
		// connection, and share it between all handlers. The handlers
		// is called in one goroutine, so no need to lock context if it
		// only be accessed in one connection.
		Context() interface{}
		SetContext(v interface{})
		Namespace() string
		Emit(msg string, v ...interface{})
	
		// Broadcast server side apis
		Join(room string)
		Leave(room string)
		LeaveAll()
		Rooms() []string
}

type conn struct {
	userid uint64
	wsC *websocket.Conn
	redisC *redis.Client
}

func NewConn(wsc *websocket.Conn,redisOptions *redis.Options) *Conn {
	var c Conn
	c.redisC = redis.NewClient(redisOptions/*&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}*/)
	c.wsC = wsc
	return c
}

func (c *conn) isUserOnline(username string) bool {
	userkey := "online." + username
	set, err := c.redisC.SetNX(userkey, username, 10*time.Second).Result()
	if err != nil {
		fmt.Println("Error on Client SetNX: ", err)
		return false
	}
	// set == 0 already online
	return set == 0
}

func (c *conn) AddUser(username string) int64{
	val, err := c.redisC.SAdd("users",username).Result()
	if err != nil {
		fmt.Println("Error on add user: ", err)
	}
	return val
}