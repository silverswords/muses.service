package main

import (
	"muses.service/apis"
)

func main() {
	r := apis.InitRouter()

	r.Run(":8080")
}
