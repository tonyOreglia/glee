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
		validMovesBb := generateValidDiagonalSlidingMovesBb(bishopPosition, mvs.pos.AllOccupiedSqsBb(), mvs.ht)
		validMovesBb.Set(validMovesBb.Value() & ^mvs.pos.ActiveSideOccupiedSqsBb())
		mvs.addValidMovesToArray(bishopPosition, validMovesBb)
	}
}

// AddValidMovesToArray save subset of valid moves from current position
func (mvs *LegalMoves) addValidMovesToArray(index int, validMovesBb *bitboard.Bitboard) {
	var validMove int
	for validMovesBb.Value() != 0 {
		validMove = validMovesBb.Lsb()
		validMovesBb.RemoveBit(validMove)
		mvs.moves = append(mvs.moves, [2]int{index, validMove})
	}
}

func generateValidDiagonalSlidingMovesBb(index int, occSqsBb uint64, ht *hashtables.HashTables) *bitboard.Bitboard {
	validDiagonalMoves, _ := bitboard.NewBitboard(
		generateValidDirectionalMovesBb(index, ht.NorthEastArrayBbHash, occSqsBb).Value() |
			generateValidDirectionalMovesBb(index, ht.NorthWestArrayBbHash, occSqsBb).Value() |
			generateValidDirectionalMovesBb(index, ht.SouthEastArrayBbHash, occSqsBb).Value() |
			generateValidDirectionalMovesBb(index, ht.SouthEastArrayBbHash, occSqsBb).Value())
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
