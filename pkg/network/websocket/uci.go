package websocket

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tonyOreglia/glee/pkg/generate"
	"github.com/tonyOreglia/glee/pkg/moves"
	"github.com/tonyOreglia/glee/pkg/position"
)

// UCI interacts with a UCI compatible chess UI
func (w *WebsocketServer) UCI(rw http.ResponseWriter, r *http.Request) {

	var err error
	log.Info("UCI websocket opened")
	w.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	w.conn, err = w.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	pos := position.StartingPosition()
	var move *moves.Move
	w.Write("GLEE-GoLang chEss Engine")
	w.Write("tony.oreglia@gmail.com")
	w.Write("id name Glee 0.0.1")
	w.Write("id author Tony Oreglia")
	w.Write("uciok")

	// for true {
	// 	_, message, err := w.conn.ReadMessage()
	// 	if err != nil {
	// 		log.Println("read error:", err)
	// 		break
	// 	}
	// 	log.Printf("recv: %s", message)
	// 	channel <- string(message)
	// }

	for true {
		_, commands, err := w.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		// severReadChan := make(chan string)
		// go w.StartReader(severReadChan) // returns values to p.Chan
		// commands := <-severReadChan // reads in a full line
		commandTokens := strings.Split(string(commands), " ")
		command := commandTokens[0]
		switch string(command) {
		case "debug":
			w.Write("not yet implemented")
		case "isready":
			w.Write("readyok")
		case "setoption":
			w.Write("not yet implemented")
		case "register":
			w.Write("not yet implemented")
		case "later":
			w.Write("not yet implemented")
		case "name":
			w.Write("not yet implemented")
		case "code":
			w.Write("not yet implemented")
		case "ucinewgame":
			pos = position.StartingPosition()
		case "position":
			pos = setPositionUCI(pos, commandTokens)
			pos.Print()
		case "go":
			pos, move = search(pos, generate.GenerateMoves(pos))
			w.Write(fmt.Sprintf("bestmove %s\n", move.String()))
		case "searchmoves":
			w.Write("not yet implemented")
		case "ponder":
			w.Write("not yet implemented")
		case "wtime":
			w.Write("not yet implemented")
		case "btime":
			w.Write("not yet implemented")
		case "winc":
			w.Write("not yet implemented")
		case "binc":
			w.Write("not yet implemented")
		case "movestogo":
			w.Write("not yet implemented")
		case "depth":
			w.Write("not yet implemented")
		case "nodes":
			w.Write("not yet implemented")
		case "mate":
			w.Write("not yet implemented")
		case "movetime":
			w.Write("not yet implemented")
		case "infinite":
			w.Write("not yet implemented")
		case "stop":
			w.Write("not yet implemented")
		case "ponderhit":
			w.Write("not yet implemented")
		case "quit":
			os.Exit(0)
		default:
			w.Write(fmt.Sprintf("Not yet implemented: %s", command))
		}
	}
}

func setPositionUCI(p *position.Position, posCommandTokens []string) *position.Position {
	// in := bufio.NewReader(os.Stdin)
	// UCIPosition, err := in.ReadString('\n')
	// if err != nil {
	// 	log.Fatal("unable to read position string", err.Error())
	// }
	log.Infof("Setting position: %s", strings.Join(posCommandTokens, ","))
	var err error
	UCIPositionTokens := posCommandTokens[1:]
	// UCIPositionTokens := strings.Split(UCIPosition, " ")

	for i, posToken := range UCIPositionTokens {
		log.Infof("handling token %s", posToken)
		if i == 0 {
			if posToken == "startpos" {
				p = position.StartingPosition()
			} else {
				p, err = position.NewPositionFen(posToken)
				if err != nil {
					badInput(strings.Join(UCIPositionTokens, " "))
					return p
				}
			}
			continue
		}
		if posToken == "moves" {
			continue
		}
		handleMove(posToken, p, generate.GenerateMoves(p))
	}
	return p
}
