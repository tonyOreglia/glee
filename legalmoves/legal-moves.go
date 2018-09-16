// Package legalmoves calculates all legal chess moves
// from a single given chess position, which is represented
// using the position package.
package legalmoves

import (
	"github.com/tonyoreglia/glee/bitboard"
	"github.com/tonyoreglia/glee/hashtables"
	"github.com/tonyoreglia/glee/moves"
	"github.com/tonyoreglia/glee/position"
)

// LegalMoves stores the legal moves from a given position
type LegalMoves struct {
	movesList *moves.Moves
	pos       *position.Position
	ht        *hashtables.HashTables
}

// NewLegalMoves exposes functionality to generate legal moves from a specific position
func NewLegalMoves(pos *position.Position, ht *hashtables.HashTables) *LegalMoves {
	resources := &LegalMoves{}
	resources.movesList = moves.NewMovesList()
	// moves.moves = make([][2]int, 0, 100)
	resources.pos = pos
	resources.ht = ht
	return resources
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

func (mvs *LegalMoves) generateMoves() {
	mvs.generatePawnMoves()
	mvs.generateKingMoves()
	mvs.generateQueenMoves()
	mvs.generateRookMoves()
	mvs.generateKnightMoves()
	mvs.generateBishopMoves()
}

func (mvs *LegalMoves) generateBishopMoves() {
	mvs.generateLegalMovesForSinglePiece(mvs.pos.GetActiveSidesBitboards()[position.Bishops].Value(), generateValidDiagonalSlidingMovesBb)
}

func (mvs *LegalMoves) generateRookMoves() {
	mvs.generateLegalMovesForSinglePiece(mvs.pos.GetActiveSidesBitboards()[position.Rooks].Value(), generateValidStraightSlidingMovesBb)
}

func (mvs *LegalMoves) generateQueenMoves() {
	mvs.generateLegalMovesForSinglePiece(mvs.pos.GetActiveSidesBitboards()[position.Queen].Value(), generateSlidingMovesBb)
}

func (mvs *LegalMoves) generateKnightMoves() {
	mvs.generateLegalMovesForSinglePiece(mvs.pos.GetActiveSidesBitboards()[position.Knights].Value(), getKnightMovesBb)
}

func (mvs *LegalMoves) generateKingMoves() {
	kingPosition := mvs.pos.GetActiveSidesBitboards()[position.King].Lsb()
	kingMovesLookup, _ := bitboard.NewBitboard(mvs.ht.LegalKingMovesBbHash[mvs.pos.GetActiveSide()][kingPosition])
	occSqsBb := mvs.pos.AllOccupiedSqsBb().Value()

	// check if castling moves blocked
	validMovesBb := generateValidDirectionalMovesBb(kingPosition, mvs.ht.EastArrayBbHash, occSqsBb, getLsb).
		Combine(generateValidDirectionalMovesBb(kingPosition, mvs.ht.WestArrayBbHash, occSqsBb, getMsb)).
		BitwiseAnd(kingMovesLookup).
		Combine(kingMovesLookup).
		RemoveOverlappingBits(mvs.pos.ActiveSideOccupiedSqsBb()).
		RemoveOverlappingBits(mvs.pos.GetActiveSideCastlingRightsBb())

	mvs.addValidMovesToArray(kingPosition, validMovesBb)
}

func (mvs *LegalMoves) getKingMovesBb(index int, occSqsBb uint64, ht *hashtables.HashTables) *bitboard.Bitboard {
	kingMoves := ht.LegalKingMovesBbHash[mvs.pos.GetActiveSide()][index]
	kingMovesBb, _ := bitboard.NewBitboard(kingMoves)

	return generateValidDirectionalMovesBb(index, ht.EastArrayBbHash, occSqsBb, getLsb).
		Combine(generateValidDirectionalMovesBb(index, ht.WestArrayBbHash, occSqsBb, getMsb)).
		BitwiseAnd(kingMovesBb).
		Combine(kingMovesBb)
}

func (mvs *LegalMoves) generatePawnMoves() {
	var getShiftedBb func(*bitboard.Bitboard, uint) *bitboard.Bitboard
	var directionOfMovement int
	var doublePushMask *bitboard.Bitboard
	var promotionRank *bitboard.Bitboard
	if mvs.pos.GetActiveSide() == position.White {
		getShiftedBb = bitboard.GetShiftedRightBb
		directionOfMovement = 1
		doublePushMask, _ = bitboard.NewBitboard(mvs.ht.FourthRankBb)
		promotionRank, _ = bitboard.NewBitboard(mvs.ht.EighthRankBb)
	} else {
		getShiftedBb = bitboard.GetShiftedLeftBb
		directionOfMovement = -1
		doublePushMask, _ = bitboard.NewBitboard(mvs.ht.FifthRankBb)
		promotionRank, _ = bitboard.NewBitboard(mvs.ht.FirstRankBb)
	}

	pawnPosBb := mvs.pos.GetActiveSidesBitboards()[position.Pawns]
	hFileBb, _ := bitboard.NewBitboard(mvs.ht.HfileBb)
	aFileBb, _ := bitboard.NewBitboard(mvs.ht.AfileBb)

	pawnAttackBb := getShiftedBb(&pawnPosBb, 9)
	pawnAttackBb.RemoveOverlappingBits(hFileBb).BitwiseAnd(mvs.pos.InactiveSideOccupiedSqsBb())
	mvs.addPawnMovesToArray(9, directionOfMovement, pawnAttackBb, promotionRank)

	pawnAttackBb = getShiftedBb(&pawnPosBb, 7)
	pawnAttackBb.RemoveOverlappingBits(aFileBb).BitwiseAnd(mvs.pos.InactiveSideOccupiedSqsBb())
	mvs.addPawnMovesToArray(7, directionOfMovement, pawnAttackBb, promotionRank)

	pawnPushBb := getShiftedBb(&pawnPosBb, 8)
	pawnPushBb.RemoveOverlappingBits(mvs.pos.AllOccupiedSqsBb())
	doubleRankpawnPushBb := getShiftedBb(pawnPushBb, 8)
	mvs.addPawnMovesToArray(8, directionOfMovement, pawnPushBb, promotionRank)
	doubleRankpawnPushBb.BitwiseAnd(doublePushMask)
	doubleRankpawnPushBb.RemoveOverlappingBits(mvs.pos.AllOccupiedSqsBb())
	mvs.addPawnMovesToArray(16, directionOfMovement, doubleRankpawnPushBb, promotionRank)
}

func (mvs *LegalMoves) addPawnMovesToArray(shift int, shiftDirection int, pawnPushBb *bitboard.Bitboard, promoRank *bitboard.Bitboard) {
	shift = shift * shiftDirection
	for pawnPushBb.Value() != 0 {
		dest := pawnPushBb.Lsb()
		pawnPushBb.RemoveBit(dest)
		origin := dest + shift
		destBb := new(bitboard.Bitboard)
		destBb.SetBit(dest)
		if destBb.BitwiseAnd(promoRank).Value() != 0 {
			mvs.movesList.AddPromotionMove(origin, dest, position.Queen)
			mvs.movesList.AddPromotionMove(origin, dest, position.Rooks)
			mvs.movesList.AddPromotionMove(origin, dest, position.Knights)
			mvs.movesList.AddPromotionMove(origin, dest, position.Bishops)
		} else {
			mvs.movesList.AddMove(origin, dest)
		}
	}
}

// AddValidMovesToArray save subset of valid moves from current position
func (mvs *LegalMoves) addValidMovesToArray(index int, validMovesBb *bitboard.Bitboard) {
	var validMove int
	for validMovesBb.Value() != 0 {
		validMove = validMovesBb.Lsb()
		validMovesBb.RemoveBit(validMove)
		mvs.movesList.AddMove(index, validMove)
		// mvs.moves = append(mvs.moves, [2]int{index, validMove})
	}
}

func getKnightMovesBb(index int, occSqsBb uint64, ht *hashtables.HashTables) *bitboard.Bitboard {
	bb, _ := bitboard.NewBitboard(ht.KnightAttackBbHash[index])
	return bb
}

func generateSlidingMovesBb(index int, occSqsBb uint64, ht *hashtables.HashTables) *bitboard.Bitboard {
	slidingMvs := generateValidDiagonalSlidingMovesBb(index, occSqsBb, ht)
	slidingMvs.Combine(generateValidStraightSlidingMovesBb(index, occSqsBb, ht))
	return slidingMvs
}

func generateValidDiagonalSlidingMovesBb(index int, occSqsBb uint64, ht *hashtables.HashTables) *bitboard.Bitboard {
	return generateValidDirectionalMovesBb(index, ht.NorthEastArrayBbHash, occSqsBb, getMsb).Combine(
		generateValidDirectionalMovesBb(index, ht.NorthWestArrayBbHash, occSqsBb, getMsb).Combine(
			generateValidDirectionalMovesBb(index, ht.SouthEastArrayBbHash, occSqsBb, getLsb).Combine(
				generateValidDirectionalMovesBb(index, ht.SouthEastArrayBbHash, occSqsBb, getLsb))))
}

func generateValidStraightSlidingMovesBb(index int, occSqsBb uint64, ht *hashtables.HashTables) *bitboard.Bitboard {
	return generateValidDirectionalMovesBb(index, ht.NorthArrayBbHash, occSqsBb, getMsb).Combine(
		generateValidDirectionalMovesBb(index, ht.SouthArrayBbHash, occSqsBb, getLsb).Combine(
			generateValidDirectionalMovesBb(index, ht.EastArrayBbHash, occSqsBb, getLsb).Combine(
				generateValidDirectionalMovesBb(index, ht.WestArrayBbHash, occSqsBb, getMsb))))
}

func generateValidDirectionalMovesBb(
	index int, directionalHash [64]uint64, occupiedSqsBb uint64, sigBit func(*bitboard.Bitboard) int) *bitboard.Bitboard {
	var validDirectionalMoves *bitboard.Bitboard
	occupiedSqsOverlapsDirectionalArray, _ := bitboard.NewBitboard(occupiedSqsBb & directionalHash[index])
	if occupiedSqsOverlapsDirectionalArray.Value() != 0 {
		sigBit := sigBit(occupiedSqsOverlapsDirectionalArray)
		// msb := occupiedSqsOverlapsDirectionalArray.Msb()
		validDirectionalMoves, _ = bitboard.NewBitboard(directionalHash[index] ^ directionalHash[sigBit])
		return validDirectionalMoves
	}
	validDirectionalMoves, _ = bitboard.NewBitboard(directionalHash[index])
	return validDirectionalMoves
}

func getLsb(bb *bitboard.Bitboard) int {
	return bb.Lsb()
}

func getMsb(bb *bitboard.Bitboard) int {
	return bb.Msb()
}
