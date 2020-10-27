package main

import (
	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"
	"github.com/tonyOreglia/glee/pkg/commandline"
	"github.com/tonyOreglia/glee/pkg/websocket"
)

func main() {
	var serve bool
	flag.BoolVar(&serve, "serve", false, "run as a webhook server (defaults to false which runs an interactive command line mode)")
	flag.Parse()

	if serve {
		log.SetFormatter(&log.JSONFormatter{})
		server := websocket.NewWebsocketServer()
		log.Info("starting websocket server")
		server.Start()
	}
	commandline.CLI()

}
