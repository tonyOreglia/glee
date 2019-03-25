package websocket

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tonyOreglia/glee/pkg/generate"
	"github.com/tonyOreglia/glee/pkg/moves"
	"github.com/tonyOreglia/glee/pkg/position"
)

// UCI interacts with a UCI compatible chess UI
func (w *WebsocketServer) UCI(rw http.ResponseWriter, r *http.Request, conn *websocket.Conn) {
	defer conn.Close()
	log.Info("websocket conection established")
	pos := position.StartingPosition()
	var move *moves.Move
	for true {
		_, commands, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		commandTokens := strings.Split(string(commands), " ")
		command := commandTokens[0]
		switch string(command) {
		case "uci":
			log.Info("executing uci response")
			Write(conn, "GLEE-GoLang chEss Engine")
			Write(conn, "tony.oreglia@gmail.com")
			Write(conn, "id name GLEE (GoLang chEss Engine) 0.0.1")
			Write(conn, "id author Tony Oreglia")
			Write(conn, "uciok")
		case "debug":
			Write(conn, "not yet implemented")
		case "isready":
			Write(conn, "readyok")
		case "setoption":
			Write(conn, "not yet implemented")
		case "register":
			Write(conn, "not yet implemented")
		case "later":
			Write(conn, "not yet implemented")
		case "name":
			Write(conn, "not yet implemented")
		case "code":
			Write(conn, "not yet implemented")
		case "ucinewgame":
			pos = position.StartingPosition()
		case "position":
			log.Info("setting engine position")
			pos = setPositionUCI(pos, commandTokens)
			pos.Print()
		case "go":
			log.Info("calculating best move")
			pos, move = search(pos, generate.GenerateMoves(pos))
			log.Infof("found best move %s", move.String())
			Write(conn, fmt.Sprintf("bestmove %s\n", move.String()))
		case "searchmoves":
			Write(conn, "not yet implemented")
		case "ponder":
			Write(conn, "not yet implemented")
		case "wtime":
			Write(conn, "not yet implemented")
		case "btime":
			Write(conn, "not yet implemented")
		case "winc":
			Write(conn, "not yet implemented")
		case "binc":
			Write(conn, "not yet implemented")
		case "movestogo":
			Write(conn, "not yet implemented")
		case "depth":
			Write(conn, "not yet implemented")
		case "nodes":
			Write(conn, "not yet implemented")
		case "mate":
			Write(conn, "not yet implemented")
		case "movetime":
			Write(conn, "not yet implemented")
		case "infinite":
			Write(conn, "not yet implemented")
		case "stop":
			Write(conn, "not yet implemented")
		case "ponderhit":
			Write(conn, "not yet implemented")
		case "quit":
			os.Exit(0)
		default:
			Write(conn, fmt.Sprintf("Not yet implemented: %s", command))
		}
	}
}

func setPositionUCI(p *position.Position, posCommandTokens []string) *position.Position {
	log.Infof("Setting position: %s", strings.Join(posCommandTokens, " "))
	var err error
	UCIPositionTokens := posCommandTokens[1:]

	if UCIPositionTokens[0] == "startpos" {
		p = position.StartingPosition()
		UCIPositionTokens = UCIPositionTokens[1:]
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
		UCIPositionTokens = UCIPositionTokens[6:]
	}

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
