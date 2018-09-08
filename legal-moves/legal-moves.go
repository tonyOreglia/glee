// Package legalmoves calculates all legal chess moves
// from a single given chess position, which is represented
// using the position package.
package legalmoves

import (
	"github.com/tonyoreglia/glee/position"
)

type LegalMoves struct {
	// slice with capacity of 100 moves, starting len of 0
	moves [][2]int
	pos   position.Position
}

func NewLegalMoves(pos position.Position) *LegalMoves {
	moves := &LegalMoves{}
	moves.moves = make([][2]int, 0, 100)
	moves.pos = pos
	return moves
}
