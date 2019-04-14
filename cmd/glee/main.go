package main

import (
	_ "github.com/qodrorid/godaemon"
	log "github.com/sirupsen/logrus"
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/network/websocket"
)

var ht = hashtables.Lookup

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	server := websocket.NewWebsocketServer()
	log.Info("starting websocket server")
	server.Start()
}
