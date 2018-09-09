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

// NewLegalMoves exposes functionality to generate legal moves from a specific position
func NewLegalMoves(pos *position.Position, ht *hashtables.HashTables) *LegalMoves {
	moves := &LegalMoves{}
	moves.moves = make([][2]int, 0, 100)
	moves.pos = pos
	moves.ht = ht
	return moves
}

func (mvs *LegalMoves) generateLegalMovesForSinglePiece(
	pieceLocationsBb uint64, genValidMovesFn func(int, uint64, *hashtables.HashTables) *bitboard.Bitboard) {

	pieceLocationBbCopy, _ := bitboard.NewBitboard(pieceLocationsBb)
	for pieceLocationBbCopy.Value() != 0 {
		piecePosition := pieceLocationBbCopy.Lsb()
		pieceLocationBbCopy.RemoveBit(piecePosition)
		validMovesBb := genValidMovesFn(piecePosition, mvs.pos.AllOccupiedSqsBb(), mvs.ht)
		// can't move to square occupied by your own pieces
		validMovesBb.RemoveOverlappingBits(mvs.pos.ActiveSideOccupiedSqsBb())
		mvs.addValidMovesToArray(piecePosition, validMovesBb)
	}
}

func (mvs *LegalMoves) generateBishopMoves() {
	mvs.generateLegalMovesForSinglePiece(mvs.pos.GetActiveSidesBitboards().Bishops.Value(), generateValidDiagonalSlidingMovesBb)
}

func (mvs *LegalMoves) generateRookMoves() {
	mvs.generateLegalMovesForSinglePiece(mvs.pos.GetActiveSidesBitboards().Rooks.Value(), generateValidStraightSlidingMovesBb)
}

func (mvs *LegalMoves) generateQueenMoves() {
	mvs.generateLegalMovesForSinglePiece(mvs.pos.GetActiveSidesBitboards().Queen.Value(), generateSlidingMovesBb)
}

func (mvs *LegalMoves) generateKnightMoves() {
	mvs.generateLegalMovesForSinglePiece(mvs.pos.GetActiveSidesBitboards().Knights.Value(), getKnightMovesBb)
}

func (mvs *LegalMoves) generateKingMoves() {
	mvs.generateLegalMovesForSinglePiece(mvs.pos.GetActiveSidesBitboards().King.Value(), mvs.getKingMovesBb)
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

func getKnightMovesBb(index int, occSqsBb uint64, ht *hashtables.HashTables) *bitboard.Bitboard {
	bb, _ := bitboard.NewBitboard(ht.KnightAttackBbHash[index])
	return bb
}

func (mvs *LegalMoves) getKingMovesBb(index int, occSqsBb uint64, ht *hashtables.HashTables) *bitboard.Bitboard {
	return generateValidDirectionalMovesBb(index, ht.LegalKingMovesBbHash[mvs.pos.GetActiveSide()], occSqsBb)
}

func generateSlidingMovesBb(index int, occSqsBb uint64, ht *hashtables.HashTables) *bitboard.Bitboard {
	slidingMvs := generateValidDiagonalSlidingMovesBb(index, occSqsBb, ht)
	slidingMvs.Combine(generateValidStraightSlidingMovesBb(index, occSqsBb, ht))
	return slidingMvs
}

func generateValidDiagonalSlidingMovesBb(index int, occSqsBb uint64, ht *hashtables.HashTables) *bitboard.Bitboard {
	validDiagonalMoves, _ := bitboard.NewBitboard(
		generateValidDirectionalMovesBb(index, ht.NorthEastArrayBbHash, occSqsBb).Value() |
			generateValidDirectionalMovesBb(index, ht.NorthWestArrayBbHash, occSqsBb).Value() |
			generateValidDirectionalMovesBb(index, ht.SouthEastArrayBbHash, occSqsBb).Value() |
			generateValidDirectionalMovesBb(index, ht.SouthEastArrayBbHash, occSqsBb).Value())
	return validDiagonalMoves
}

func generateValidStraightSlidingMovesBb(index int, occSqsBb uint64, ht *hashtables.HashTables) *bitboard.Bitboard {
	validStraightMoves, _ := bitboard.NewBitboard(
		generateValidDirectionalMovesBb(index, ht.NorthArrayBbHash, occSqsBb).Value() |
			generateValidDirectionalMovesBb(index, ht.SouthArrayBbHash, occSqsBb).Value() |
			generateValidDirectionalMovesBb(index, ht.EastArrayBbHash, occSqsBb).Value() |
			generateValidDirectionalMovesBb(index, ht.WestArrayBbHash, occSqsBb).Value())
	return validStraightMoves
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
