package engine

import (
	"fmt"

	"github.com/tonyOreglia/glee/pkg/bitboard"
	"github.com/tonyOreglia/glee/pkg/evaluate"
	"github.com/tonyOreglia/glee/pkg/generate"
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/moves"
	"github.com/tonyOreglia/glee/pkg/position"
)

var ht = hashtables.Lookup

type searchParams struct {
	depth           int
	ply             int
	pos             **position.Position
	engineMove      *moves.Move
	perft           *int
	singlePlyPerft  *int
	evaluationScore int
	root            bool
}

func minMax(p searchParams) int {
	p.root = p.ply == p.depth
	if p.ply == 0 {
		return evaluate.EvaluatePosition(*p.pos, p.singlePlyPerft)
	}
	p.evaluationScore = 30000
	if (*p.pos).IsWhitesTurn() {
		p.evaluationScore = -30000
	}
	mvList := generate.GenerateMoves(*p.pos).GetMovesList()
	for _, move := range mvList {
		if !makeValidMove(move, p.pos) {
			continue
		}
		evaluateMove(move, p)
	}
	return p.evaluationScore
}

// evaluateMove checks if pseudo legal move is valid, and if so
// increments ply and calls minMax to continue the search
func evaluateMove(move moves.Move, p searchParams) error {
	p.ply = p.ply - 1
	temp := minMax(p)
	if temp > p.evaluationScore {
		p.evaluationScore = temp
		if p.root {
			p.engineMove = move.CopyMove()
		}
	}
	if p.root {
		move.Print()
		fmt.Println(*p.singlePlyPerft)
		*p.perft += *p.singlePlyPerft
		*p.singlePlyPerft = 0
	}
	*p.pos = (*p.pos).UnMakeMove()
	return nil
}

func makeValidMove(move moves.Move, pos **position.Position) bool {
	if (*pos).IsCastlingMove(move) {
		(*pos).Move(move)
		if !castlingMoveIsValid(move, pos) {
			*pos = (*pos).UnMakeMove()
			return false
		}
	} else {
		(*pos).Move(move)
	}
	legalMoves := generate.GenerateMoves(*pos)
	if (*pos).IsAttacked((*pos).InactiveSideKingBb(), legalMoves.AttackedSqsBb()) {
		*pos = (*pos).UnMakeMove()
		return false
	}
	return true
}

func castlingMoveIsValid(move moves.Move, pos **position.Position) bool {
	kingPosition := bitboard.NewBitboardFromIndex(move.Origin())
	legalMoves := generate.GenerateMoves(*pos)
	castlingSlidingSqBb := bitboard.NewBitboard(uint64(0))
	castlingSlidingSqBb.SetBit(int(ht.LookupCastlingSlidingSqByDest[uint64(move.Destination())]))
	castlingSlidingSqBb.Combine(kingPosition)
	if (*pos).IsAttacked(*castlingSlidingSqBb, legalMoves.AttackedSqsBb()) {
		return false
	}
	return true
}
