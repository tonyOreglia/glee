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
	w.upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	w.conn, err = w.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	log.Info("websocket conection established")
	pos := position.StartingPosition()
	var move *moves.Move
	for true {
		_, commands, err := w.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		commandTokens := strings.Split(string(commands), " ")
		command := commandTokens[0]
		log.Infof("Got command %s", string(commands))
		switch string(command) {
		case "uci":
			log.Info("executing uci response")
			w.Write("GLEE-GoLang chEss Engine")
			w.Write("tony.oreglia@gmail.com")
			w.Write("id name GLEE (GoLang chEss Engine) 0.0.1")
			w.Write("id author Tony Oreglia")
			w.Write("uciok")
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
			log.Info("setting engine position")
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
	log.Infof("Setting position: %s", strings.Join(posCommandTokens, " "))
	var err error
	UCIPositionTokens := posCommandTokens[1:]
	for i, posToken := range UCIPositionTokens {
		log.Infof("handling position token %s", posToken)
		if i == 0 {
			if posToken == "startpos" {
				p = position.StartingPosition()
			} else {
				if len(UCIPositionTokens) < 6 {
					log.Errorf("invalid position: %s", strings.Join(posCommandTokens, " "))
					return nil
				}
				fenString := strings.Join(UCIPositionTokens[0:6], " ")
				p, err = position.NewPositionFen(fenString)
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
