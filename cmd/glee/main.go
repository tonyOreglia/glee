package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/network/websocket"
)

var ht = hashtables.Lookup

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	server := websocket.NewWebsocketServer()
	server.Start()
}
