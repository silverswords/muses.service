package main

import (
	"fmt"

	"muses.service/apis"
	"muses.service/service/eventbus"
	"muses.service/service/room"
)

func main() {
	eb := eventbus.New()
	roomChannel := eb.Attach("roomManager")
	connectionChannel := eb.Attach("connectionManager")

	r := apis.InitRouter(eb)

	roomManager := room.Init()
	go func() {
		fmt.Println("== start listing roomChannel ==")
		for {
			event, err := roomChannel.Receive()
			if err != nil {
				fmt.Println(err)
				return
			}

			switch {
			case "createRoom" == event.Type():
				err := roomManager.InitRoom(event.Payload())
				if err != nil {
					fmt.Println(err)
				}
			}

			fmt.Println(event)
		}
	}()

	go func() {
		fmt.Println("== start listing roomChannel ==")
		for {
			event, err := connectionChannel.Receive()
			if err != nil {
				fmt.Println(err)
				return
			}

			switch {
			case "createRoom" == event.Type():
				roomManager.InitRoom(event.Payload())
			}
		}
	}()

	r.Run(":8080")
}
