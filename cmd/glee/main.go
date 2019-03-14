package main

import (
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/network/websocket"
)

var ht = hashtables.Lookup

func main() {
	server := websocket.NewWebsocketServer()
	server.Start()
}
