package connection_test

import (
	"flag"
	"testing"

	conn "muses.service/service/connection"
)

var addr = flag.String("addr", ":8080", "http service address")

func TestManager(t *testing.T) {
	connectionManager := conn.NewConnectionManager()
	go connectionManager.Run()

	t.Log("run ConnectionManager")
}
