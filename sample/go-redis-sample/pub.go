/**
 * @author abser
 * @email [abser@foxmail.com]
 * @create date 2020-02-14 22:37:58
 * @modify date 2020-02-14 22:37:58
 * @desc [description]
 */
package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func main() {

	var client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	go func() {
		for {
			client.Publish("room1", "hello")
			time.Sleep(time.Second)
		}
	}()

	for {
		pubsub := client.Subscribe("room1")
		_, err := pubsub.Receive()
		if err != nil {
			return
		}
		ch := pubsub.Channel()
		for msg := range ch {
			fmt.Println(msg.Channel, msg.Payload, "\r\n")
		}
	}
}
