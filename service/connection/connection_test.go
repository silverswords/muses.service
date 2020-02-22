package connection_test

import (
	"testing"

	conn "muses.service/service/connection"
)

func TestManager(t *testing.T) {
	connectionManager := conn.NewConnectionManager()
	go connectionManager.Run()

	t.Log("run ConnectionManager")
}
