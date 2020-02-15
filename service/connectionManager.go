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
		LocalAddr() net.Addr
		RemoteAddr() net.Addr
	
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
	user string
	wsC *websocket.Conn
	redisC *redis.Client
}

func NewConn(wsc *websocket.Conn,redisOptions *redis.Options, username string) *Conn {
	var c Conn
	c.redisC = redis.NewClient(redisOptions/*&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}*/)
	c.wsC = wsc
	c.user = username
	return c
}

func (c *conn) ID() string {
	return c.user
}

func(c *conn) Close() error {
	_ := c.redisC.Close()
	return c.wsC.Close()
}

func (c *conn) LocalAddr() net.Addr {
	return c.wsC.UnderlyingConn().LocalAddr()
}

func (c *conn) RemoteAddr()  net.Addr {
	return c.wsC.UnderlyingConn().RemoteAddr()
}

// usage 
// for msg := range ch {
// 	fmt.Println("send: ", msg.Channel, msg.Payload)
// 	err = ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
// }
func (c *conn) Subscribe(room string)  <-chan *Message {
	pubsub := c.redisC.Subscribe(room)
	_, err = pubsub.Receive()
	if err != nil {
		return nil
	}
	return pubsub.Channel()
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