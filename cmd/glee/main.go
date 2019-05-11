package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/tonyOreglia/glee/pkg/websocket"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	server := websocket.NewWebsocketServer()
	log.Info("starting websocket server")
	server.Start()
}
