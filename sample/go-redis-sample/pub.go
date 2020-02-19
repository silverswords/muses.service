/**
 * @author abser
 * @email [abser@foxmail.com]
 * @create date 2020-02-14 22:37:58
 * @modify date 2020-02-14 22:37:58
 * @desc [description]
 */
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v7"
)

type foo struct {
	Boo int `json:"boo"`
	Hei int `json:"hei"`
}

func main() {
	var client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// var hello = "hello"
	// hello := map[string]int{"hello": 1}
	// hello := 23.23
	// msg, _ := json.Marshal(&foo{Boo: 1, Hei: 2})
	msgH := make(http.Header)
	msgH.Add("apple", "yellow")
	msg, _ := json.Marshal(msgH)
	fmt.Println(msg)
	go func() {
		for {
			t1 := time.Now()
			client.Publish("/room1", msg)
			fmt.Println("pub", time.Now().Sub(t1))

			time.Sleep(time.Second)
		}
	}()

	pubsub := client.Subscribe("/room1")
	_, err := pubsub.Receive()
	if err != nil {
		return
	}
	ch := pubsub.Channel()
	for msg := range ch {
		t2 := time.Now()
		// fmt.Println(msg.Channel, msg.Payload)
		recvH := make(http.Header)
		json.Unmarshal([]byte(msg.Payload), &recvH)
		fmt.Println(recvH.Get("apple"))
		fmt.Println("sub", time.Now().Sub(t2))
	}

}
