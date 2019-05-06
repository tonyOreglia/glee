package websocket

import (
	"net/http"

	"github.com/namsral/flag"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	upgrader websocket.Upgrader
	addr     *string
}

func NewWebsocketServer() *WebsocketServer {
	w := new(WebsocketServer)
	w.addr = flag.String("addr", "localhost:8081", "http websocket service address")
	flag.Parse()
	w.upgrader = websocket.Upgrader{} // use default options
	http.HandleFunc("/uci", w.uciHandler)
	return w
}

func (w *WebsocketServer) uciHandler(rw http.ResponseWriter, r *http.Request) {
	log.Error("upgrading to websocket connection")
	w.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := w.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	go w.UCI(rw, r, conn)
}

func (w *WebsocketServer) Start() {
	log.Info("starting websocket server")
	log.Fatal(http.ListenAndServe(*w.addr, nil))
}

func Write(conn *websocket.Conn, msg string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write:", err)
	}
}
