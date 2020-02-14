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

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"strconv"

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
	defer ws.Close()
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		// get couter and return +=1
		key := "counter"
		myval, _ := client.Get(key).Result()
		intv, _ := strconv.Atoi(string(myval))
		intstr := strconv.Itoa(intv)
		err = ws.WriteMessage(mt, []byte(intstr))

		log.Printf("send: %s", intstr)

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
