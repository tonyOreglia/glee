package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	ws "github.com/tonyOreglia/glee/pkg/network/websocket"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	wsServer := new(ws.WebsocketServer)
	wsServer.Upgrader = websocket.Upgrader{} // use default options
	wsServer.UCI(w, r)
}
