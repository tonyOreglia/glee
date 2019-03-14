package websocket

import (
	"flag"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type WebsocketServer struct {
	upgrader websocket.Upgrader
	addr     *string
	conn     *websocket.Conn
}

func NewWebsocketServer() *WebsocketServer {
	w := new(WebsocketServer)
	w.addr = flag.String("addr", "localhost:8081", "http service address")
	flag.Parse()
	w.upgrader = websocket.Upgrader{} // use default options
	http.HandleFunc("/uci", w.UCI)
	return w
}

func (w *WebsocketServer) Start() {
	http.ListenAndServe(*w.addr, nil)
}

func (w *WebsocketServer) CloseConnection() {
	w.conn.Close()
}

func (w *WebsocketServer) StartReader(channel chan string) {
	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		log.Printf("recv: %s", message)
		channel <- string(message)
	}
}

func (w *WebsocketServer) Write(msg string) {
	err := w.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write:", err)
	}
}

// func (w *WebsocketServer) uci(rw http.ResponseWriter, r *http.Request) {
// 	var err error
// 	log.Info("UCI websocket opened")
// 	w.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
// 	w.conn, err = w.upgrader.Upgrade(rw, r, nil)
// 	if err != nil {
// 		log.Print("upgrade:", err)
// 		return
// 	}
// }
