package websocket

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/tonyOreglia/glee/pkg/engine"
	"github.com/tonyOreglia/glee/pkg/moves"
	"github.com/tonyOreglia/glee/pkg/position"
)

func handleMove(mv string, p *position.Position, mvs *moves.Moves) bool {
	lookupPromo := map[string]int{
		// Queen = 2 Bishops = 3 Knights = 4 Rooks = 5
		"Q": 2,
		"B": 3,
		"N": 4,
		"R": 5,
	}
	promotionPiece := 0
	if len(mv) != 4 && len(mv) != 5 {
		badInput(mv)
		return false
	}
	if len(mv) == 5 {
		promotionPiece = lookupPromo[string(mv[4])]
	}
	origin, err := moves.ConvertAlgebriacToIndex(mv[0:2])
	if err != nil {
		badInput(mv)
		return false
	}
	dest, err := moves.ConvertAlgebriacToIndex(mv[2:4])
	if err != nil {
		badInput(mv)
		return false
	}
	move, found := mvs.FindMove(origin, dest, promotionPiece)
	if !found {
		badInput(mv)
		return false
	}
	if !engine.MakeValidMove(move, &p) {
		badInput(mv)
		return false
	}
	return true
}

func setboard(p *position.Position) {
	in := bufio.NewReader(os.Stdin)
	fen, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fen = fen[0 : len(fen)-1]
	p, err = position.NewPositionFen(string(fen))
	if err != nil {
		badInput(string(fen))
	}
	p.Print()
}

func search(p *position.Position, mvs *moves.Moves) (*position.Position, *moves.Move) {
	perft := 0
	singlePlyPerft := 0
	params := engine.SearchParams{
		Depth:          5,
		Ply:            5,
		Pos:            &p,
		Perft:          &perft,
		SinglePlyPerft: &singlePlyPerft,
		EngineMove:     &moves.Move{},
	}
	if p.IsWhitesTurn() {
		engine.AlphaBetaMax(-10000, 10000, 5, params)
	} else {
		engine.AlphaBetaMin(-10000, 10000, 5, params)
	}
	return p, params.EngineMove
}

func badInput(c string) {
	fmt.Printf("\ninput correct ?: %s\n\n", c)
}

func undo(p *position.Position) {
	newPos := p.UnMakeMove()
	if newPos != nil {
		p = newPos
		p.Print()
		return
	}
	badInput("no previous move to undo")
}
