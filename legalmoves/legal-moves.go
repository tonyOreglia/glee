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
	bishopBbCopy, _ := bitboard.NewBitboard(mvs.pos.GetActiveSidesBitboards().Bishops.Value())
	for bishopBbCopy.Value() != 0 {
		bishopPosition := bishopBbCopy.Lsb()
		bishopBbCopy.RemoveBit(bishopPosition)
		validMovesBb := mvs.generateValidDiagonalSlidingMovesBb(bishopPosition).Value()
		validMovesBb &= ^mvs.pos.ActiveSideOccupiedSqsBb()
		mvs.addValidMovesToArray(bishopPosition, validMovesBb)
	}
}

// AddValidMovesToArray save subset of valid moves from current position
func (mvs *LegalMoves) addValidMovesToArray(index int, validMovesBb uint64) {
	mvs.moves = append(mvs.moves, [2]int{index})
}

func (mvs *LegalMoves) generateValidDiagonalSlidingMovesBb(index int) *bitboard.Bitboard {
	occSqsBb := mvs.pos.AllOccupiedSqsBb()
	validDiagonalMoves, _ := bitboard.NewBitboard(
		generateValidDirectionalMovesBb(index, mvs.ht.NorthEastArrayBbHash, occSqsBb).Value() |
			generateValidDirectionalMovesBb(index, mvs.ht.NorthWestArrayBbHash, occSqsBb).Value() |
			generateValidDirectionalMovesBb(index, mvs.ht.SouthEastArrayBbHash, occSqsBb).Value() |
			generateValidDirectionalMovesBb(index, mvs.ht.SouthEastArrayBbHash, occSqsBb).Value())
	return validDiagonalMoves
}

func generateValidDirectionalMovesBb(index int, directionalHash [64]uint64, occupiedSqsBb uint64) *bitboard.Bitboard {
	var validDirectionalMoves *bitboard.Bitboard
	occupiedSqsOverlapsNEastArrayBb, _ := bitboard.NewBitboard(occupiedSqsBb & directionalHash[index])
	if occupiedSqsOverlapsNEastArrayBb.Value() != 0 {
		msb := occupiedSqsOverlapsNEastArrayBb.Msb()
		validDirectionalMoves, _ = bitboard.NewBitboard(directionalHash[index] ^ directionalHash[msb])
		return validDirectionalMoves
	}
	validDirectionalMoves, _ = bitboard.NewBitboard(directionalHash[index])
	return validDirectionalMoves
}
