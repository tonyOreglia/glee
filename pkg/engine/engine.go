package engine

import (
	"fmt"
	"log"

	"github.com/tonyOreglia/glee/pkg/bitboard"
	"github.com/tonyOreglia/glee/pkg/evaluate"
	"github.com/tonyOreglia/glee/pkg/generate"
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/moves"
	"github.com/tonyOreglia/glee/pkg/position"
)

var ht = hashtables.Lookup

type searchParams struct {
	depth          int
	ply            int
	pos            **position.Position
	engineMove     **moves.Move
	perft          *int
	singlePlyPerft *int
}

func minMax(p searchParams) int {
	var value, tempValue int
	root := p.ply == p.depth
	if root {
		fmt.Println("")
		fmt.Print("Perft Divide: ")
		fmt.Println((*p.pos).GetFenString())
	}
	generate.GenerateMoves(*p.pos)
	if p.ply == 0 {
		return evaluate.EvaluatePosition(*p.pos, p.singlePlyPerft)
	}
	if (*p.pos).IsWhitesTurn() {
		value = -30000
		generateMoves(&value, &tempValue, root, p)
		return value
	}
	value = 30000
	generateMoves(&value, &tempValue, root, p)
	return value
}

func generateMoves(value *int, tempValue *int, root bool, p searchParams) error {
	mvList := generate.GenerateMoves(*p.pos).GetMovesList()
	for _, move := range mvList {
		evaluateMove(move, value, tempValue, root, p)
	}
	return nil
}

func evaluateMove(move moves.Move, value *int, tempValue *int, root bool, p searchParams) error {
	if (*p.pos).IsCastlingMove(move) {
		if !castlingMoveIsValid(move, value, tempValue, root, p) {
			return nil
		}
	}
	(*p.pos).Move(move)
	legalMoves := generate.GenerateMoves(*p.pos)
	// mg.GenerateMoves()
	if (*p.pos).IsAttacked((*p.pos).InactiveSideKingBb(), legalMoves.AttackedSqsBb()) {
		*p.pos = (*p.pos).UnMakeMove()
		return nil
	}
	p.ply = p.ply - 1
	*tempValue = minMax(p)
	if *tempValue > *value {
		value = tempValue
		if root {
			*p.engineMove = move.CopyMove()
		}
	}
	if root {
		move.Print()
		fmt.Println(*p.singlePlyPerft)
		*p.perft += *p.singlePlyPerft
		*p.singlePlyPerft = 0
	}
	*p.pos = (*p.pos).UnMakeMove()
	return nil
}

func castlingMoveIsValid(move moves.Move, value *int, tempValue *int, root bool, p searchParams) bool {
	tempPos := *p.pos
	*p.pos = (*p.pos).UnMakeMove()
	legalMoves := generate.GenerateMoves(*p.pos)
	// mg.GenerateMoves()
	castlingSlidingSqBb, err := bitboard.NewBitboard(uint64(0))
	castlingSlidingSqBb.SetBit(int(ht.LookupCastlingSlidingSqByDest[uint64(move.Destination())]))
	if err != nil {
		log.Fatal(err.Error())
	}
	if (*p.pos).IsAttacked(*castlingSlidingSqBb, legalMoves.AttackedSqsBb()) {
		*p.pos = tempPos
		return false
	}
	*p.pos = tempPos
	return true
}
