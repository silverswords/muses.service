/**
 * @author abser
 * @email [abser@foxmail.com]
 * @create date 2020-02-14 21:20:08
 * @modify date 2020-02-14 21:20:08
 * @desc [description]
 */

package main

import (
	"flag"
	"fmt"
	"time"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8082", "http service address")
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func echo(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade ws:", err)
		return
	}

	_, err = client.SetNX("userkey", "username", 10*time.Second).Result()
	if err != nil {
		fmt.Println("Error on Client SetNX", err)
		return
	}

	// tickerChan := time.NewTicker(time.Second * 60).C

	var subclient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	val, err := subclient.SAdd("users", "username").Result()
	if err != nil {
		fmt.Println("Error on add user: ", err)
	}
	fmt.Println(val)

	pubsub := subclient.Subscribe("room1")
	_, err = pubsub.Receive()
	if err != nil {
		return
	}
	ch := pubsub.Channel()

	// handle sub messages
	go func() {
		for msg := range ch {
			// t2 := time.Now()
			fmt.Println("send: ", msg.Channel, msg.Payload)
			err = ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
			// fmt.Println("sub", time.Now().Sub(t2))
		}
	}()

	defer ws.Close()
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		subclient.Publish("room1", message)
		// get couter and return +=1
		// key := "counter"
		// myval, _ := client.Get(key).Result()
		// intv, _ := strconv.Atoi(string(myval))
		// dosomething on intv
		// intstr := strconv.Itoa(intv)

		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	r := gin.Default()

	r.GET("/echo", echo)
	// http.HandleFunc("/", home)
	log.Fatal(r.Run(*addr))
}
