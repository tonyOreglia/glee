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
		validMovesBb := genValidMovesFn(piecePosition, mvs.pos.AllOccupiedSqsBb().Value(), mvs.ht)
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
	piecePosition := mvs.pos.GetActiveSidesBitboards().King.Lsb()
	validMovesBb := mvs.getKingMovesBb(piecePosition, mvs.pos.AllOccupiedSqsBb().Value(), mvs.ht)
	validMovesBb.RemoveOverlappingBits(mvs.pos.ActiveSideOccupiedSqsBb())
	validMovesBb.RemoveOverlappingBits(mvs.pos.GetActiveSideCastlingRightsBb())
	mvs.addValidMovesToArray(piecePosition, validMovesBb)
}

// Positive 1 shift direction for white, negative 1 for black.
func (mvs *LegalMoves) generatePawnMoves() {
	var getShiftedBb func(*bitboard.Bitboard, uint) *bitboard.Bitboard
	var directionOfMovement int
	if mvs.pos.GetActiveSide() == position.White {
		getShiftedBb = bitboard.GetShiftedRightBb
		directionOfMovement = 1
	} else {
		getShiftedBb = bitboard.GetShiftedLeftBb
		directionOfMovement = -1
	}

	pawnPosBb := mvs.pos.GetActiveSidesBitboards().Pawns
	hFileBb, _ := bitboard.NewBitboard(mvs.ht.HfileBb)
	aFileBb, _ := bitboard.NewBitboard(mvs.ht.AfileBb)

	pawnAttackBb := getShiftedBb(&pawnPosBb, 9)
	pawnAttackBb.RemoveOverlappingBits(hFileBb).BitwiseAnd(mvs.pos.InactiveSideOccupiedSqsBb())
	mvs.addPawnMovesToArray(9, directionOfMovement, pawnAttackBb)

	pawnAttackBb = getShiftedBb(&pawnPosBb, 7)
	pawnAttackBb.RemoveOverlappingBits(aFileBb).BitwiseAnd(mvs.pos.InactiveSideOccupiedSqsBb())
	mvs.addPawnMovesToArray(7, directionOfMovement, pawnAttackBb)

	pawnPushBb := getShiftedBb(&pawnPosBb, 8)
	pawnPushBb.RemoveOverlappingBits(mvs.pos.AllOccupiedSqsBb())
	doubleRankpawnPushBb := getShiftedBb(pawnPushBb, 8)
	mvs.addPawnMovesToArray(8, directionOfMovement, pawnPushBb)
	doubleRankpawnPushBb.RemoveOverlappingBits(mvs.pos.AllOccupiedSqsBb())
	mvs.addPawnMovesToArray(16, directionOfMovement, doubleRankpawnPushBb)
}

func (mvs *LegalMoves) addPawnMovesToArray(shift int, shiftDirection int, pawnPushBb *bitboard.Bitboard) {
	shift = shift * shiftDirection
	for pawnPushBb.Value() != 0 {
		dest := pawnPushBb.Lsb()
		pawnPushBb.RemoveBit(dest)
		origin := dest + shift
		mvs.moves = append(mvs.moves, [2]int{origin, dest})
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
	return generateValidDirectionalMovesBb(index, ht.NorthEastArrayBbHash, occSqsBb).Combine(
		generateValidDirectionalMovesBb(index, ht.NorthWestArrayBbHash, occSqsBb).Combine(
			generateValidDirectionalMovesBb(index, ht.SouthEastArrayBbHash, occSqsBb).Combine(
				generateValidDirectionalMovesBb(index, ht.SouthEastArrayBbHash, occSqsBb))))
}

func generateValidStraightSlidingMovesBb(index int, occSqsBb uint64, ht *hashtables.HashTables) *bitboard.Bitboard {
	return generateValidDirectionalMovesBb(index, ht.NorthArrayBbHash, occSqsBb).Combine(
		generateValidDirectionalMovesBb(index, ht.SouthArrayBbHash, occSqsBb).Combine(
			generateValidDirectionalMovesBb(index, ht.EastArrayBbHash, occSqsBb).Combine(
				generateValidDirectionalMovesBb(index, ht.WestArrayBbHash, occSqsBb))))
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
