package main

import (
	"fmt"
	"time"

	eventbus "muses.service/service/eventbus"
)

func add(a string) {
	fmt.Println("ds", time.Now())
	fmt.Println(a)
}

func main() {

	evb, err := eventbus.New()
	if err != nil {
		fmt.Println("Couldn't read env url to create evb")
	}
	go evb.Register("add", add)
	evb.Publish("add", "hello")
	fmt.Println(time.Now())

	time.Sleep(time.Second * 6)
	evb.Publish("add", "/exit")
	time.Sleep(time.Second * 1)
}
