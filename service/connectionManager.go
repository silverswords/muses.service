package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"

	"github.com/go-redis/redis"
	"github.com/go-redis/redis/v7"
)

type Conn struct {
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

func (c *Conn) isUserOnline(username string) bool {
	userkey := "online." + username
	set, err := c.redisC.SetNX(userkey, username, 10*time.Second).Result()
	if err != nil {
		fmt.Println("Error on Client SetNX", err)
		return false
	}
	return set
}
