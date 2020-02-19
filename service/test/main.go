package main

import (
	"fmt"
	"strconv"

	eb "muses.service/service/eventbus"
)

func add(e *eb.Event) {
	a, _ := strconv.Atoi(e.Get("a"))
	b, _ := strconv.Atoi(e.Get("b"))

	fmt.Println(a + b)
}

func main() {

	eb.HandleFunc("/add", add)

	eb.Serve()
}
