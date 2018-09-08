// Package legalmoves calculates all legal chess moves
// from a single given chess position, which is represented
// using the position package.
package legalmoves

import (
	"github.com/tonyoreglia/glee/bitboard"
	"github.com/tonyoreglia/glee/hashtables"
	"github.com/tonyoreglia/glee/position"
)

// LegalMoves stores the legal moves from a given position
type LegalMoves struct {
	moves [][2]int
	pos   *position.Position
	ht    *hashtables.HashTables
}

func NewLegalMoves(pos *position.Position, ht *hashtables.HashTables) *LegalMoves {
	moves := &LegalMoves{}
	moves.moves = make([][2]int, 0, 100)
	moves.pos = pos
	moves.ht = ht
	return moves
}

func (mvs *LegalMoves) generateBishopMoves() {
	bishopBbCopy := mvs.pos.GetActiveSidesBitboards().Bishops
	for bishopBbCopy.Value() != 0 {
		bishopPosition := bishopBbCopy.Lsb()
		bishopBbCopy.RemoveBit(bishopPosition)
		validMovesBb := mvs.generateValidDiagonalSlidingMovesBb(bishopPosition)
		validMovesBb &= ^mvs.pos.ActiveSideOccupiedSqsBb()
		mvs.addValidMovesToArray(bishopPosition, validMovesBb)
	}
}

// AddValidMovesToArray save subset of valid moves from current position
func (mvs *LegalMoves) addValidMovesToArray(index int, validMovesBb uint64) {

}

func (mvs *LegalMoves) generateValidDiagonalSlidingMovesBb(index int) uint64 {
	return mvs.generateValidDirectionalMovesBb(index, mvs.ht.NorthEastArrayBbHash) |
		mvs.generateValidDirectionalMovesBb(index, mvs.ht.NorthWestArrayBbHash) |
		mvs.generateValidDirectionalMovesBb(index, mvs.ht.SouthEastArrayBbHash) |
		mvs.generateValidDirectionalMovesBb(index, mvs.ht.SouthEastArrayBbHash)
}

func (mvs *LegalMoves) generateValidDirectionalMovesBb(index int, directionalHash [64]uint64) uint64 {
	occupiedSqsOverlapsNEastArrayBb, _ := bitboard.NewBitboard(mvs.pos.AllOccupiedSqsBb() & directionalHash[index])
	if occupiedSqsOverlapsNEastArrayBb.Value() != 0 {
		msb := occupiedSqsOverlapsNEastArrayBb.Msb()
		return directionalHash[index] ^ directionalHash[msb]
	}
	return directionalHash[index]
}
