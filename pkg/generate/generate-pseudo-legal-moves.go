// Package generate calculates all legal chess moves
// from a single given chess position using position library.
package generate

import (
	"github.com/tonyOreglia/glee/pkg/bitboard"
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/moves"
	"github.com/tonyOreglia/glee/pkg/position"
)

// NewLegalMoveGenerator exposes functionality to generate legal moves from a specific position
func GenerateMoves(pos *position.Position) *moves.Moves {
	movesList := moves.NewMovesList()
	ht := hashtables.Lookup
	GenerateAllMoves(pos, movesList, ht)
	return movesList
}

// GenerateMoves generates a list of pseudo legal moves for a given postion.
// Moves that expose the moving sides king to check may be included in the move list.
// This condition is checked for by the engine when thinking.
func GenerateAllMoves(pos *position.Position, mvsList *moves.Moves, ht *hashtables.HashTables) {
	GeneratePawnMoves(pos, mvsList, ht)
	GenerateKingMoves(pos, mvsList, ht)
	GenerateQueenMoves(pos, mvsList, ht)
	GenerateRookMoves(pos, mvsList, ht)
	GenerateKnightMoves(pos, mvsList, ht)
	GenerateBishopMoves(pos, mvsList, ht)
}

func generateLegalMovesForSinglePiece(
	pos *position.Position, movesList *moves.Moves, pieceLocationsBb uint64, genValidMovesFn func(int, uint64, *hashtables.HashTables) *bitboard.Bitboard, ht *hashtables.HashTables) {

	pieceLocationBbCopy, _ := bitboard.NewBitboard(pieceLocationsBb)
	for pieceLocationBbCopy.Value() != 0 {
		piecePosition := pieceLocationBbCopy.Lsb()
		pieceLocationBbCopy.RemoveBit(piecePosition)
		validMovesBb := genValidMovesFn(piecePosition, pos.AllOccupiedSqsBb().Value(), ht)
		// can't move to square occupied by your own pieces
		validMovesBb.RemoveOverlappingBits(pos.ActiveSideOccupiedSqsBb())
		addValidMovesToArray(movesList, piecePosition, validMovesBb)
	}
}

func GenerateBishopMoves(pos *position.Position, mvsList *moves.Moves, ht *hashtables.HashTables) {
	generateLegalMovesForSinglePiece(pos, mvsList, pos.GetActiveSidesBitboards()[position.Bishops].Value(), generateValidDiagonalSlidingMovesBb, ht)
}

func GenerateRookMoves(pos *position.Position, mvsList *moves.Moves, ht *hashtables.HashTables) {
	generateLegalMovesForSinglePiece(pos, mvsList, pos.GetActiveSidesBitboards()[position.Rooks].Value(), generateValidStraightSlidingMovesBb, ht)
}

func GenerateQueenMoves(pos *position.Position, mvsList *moves.Moves, ht *hashtables.HashTables) {
	generateLegalMovesForSinglePiece(pos, mvsList, pos.GetActiveSidesBitboards()[position.Queen].Value(), generateSlidingMovesBb, ht)
}

func GenerateKnightMoves(pos *position.Position, mvsList *moves.Moves, ht *hashtables.HashTables) {
	generateLegalMovesForSinglePiece(pos, mvsList, pos.GetActiveSidesBitboards()[position.Knights].Value(), getKnightMovesBb, ht)
}

// GenerateKingMovesFromInitialPosition checks valid moves with castling taken into account
func GenerateKingMovesFromInitialPosition(pos *position.Position, mvsList *moves.Moves, ht *hashtables.HashTables) *bitboard.Bitboard {
	kingBb := pos.GetActiveSidesBitboards()[position.King]
	kingPosition := kingBb.Lsb()
	kingMovesLookup, _ := bitboard.NewBitboard(ht.LegalKingMovesBbHash[pos.GetActiveSide()][kingPosition])
	castlingBits, _ := bitboard.NewBitboard(ht.CastlingBits[0] | ht.CastlingBits[1])
	occSqsBb := pos.AllOccupiedSqsBb().Value()
	finalRankBb, _ := bitboard.NewBitboard(ht.EighthRankBb | ht.FirstRankBb)
	return generateValidDirectionalMovesBb(kingPosition, ht.EastArrayBbHash, occSqsBb, getLsb).
		Combine(generateValidDirectionalMovesBb(kingPosition, ht.WestArrayBbHash, occSqsBb, getMsb)).
		// only allow sliding moves that are withing the legal king moves
		BitwiseAnd(kingMovesLookup).
		// cannot move to square that is occupied by same side
		RemoveOverlappingBits(pos.ActiveSideOccupiedSqsBb()).
		// remove castling move if permissions is not set
		RemoveOverlappingBits(pos.GetActiveSideCastlingRightsBb()).
		// remove castling moves if there is an opposing piece there
		RemoveOverlappingBits(bitboard.ReturnBitwiseAnd(pos.InactiveSideOccupiedSqsBb(), castlingBits)).
		// add in moves that are not on the final rank
		Combine(bitboard.ReturnBitwiseAnd(kingMovesLookup.RemoveOverlappingBits(pos.ActiveSideOccupiedSqsBb()), finalRankBb.ReturnBitsFlipped()))
}

func GenerateKingMoves(pos *position.Position, mvsList *moves.Moves, ht *hashtables.HashTables) {
	kingBb := pos.GetActiveSidesBitboards()[position.King]
	kingPosition := kingBb.Lsb()
	kingMovesLookup, _ := bitboard.NewBitboard(ht.LegalKingMovesBbHash[pos.GetActiveSide()][kingPosition])
	var validMovesBb *bitboard.Bitboard
	if kingPosition == 4 || kingPosition == 60 {
		validMovesBb = GenerateKingMovesFromInitialPosition(pos, mvsList, ht)
	} else {
		validMovesBb = kingMovesLookup.RemoveOverlappingBits(pos.ActiveSideOccupiedSqsBb())
	}

	addValidMovesToArray(mvsList, kingPosition, validMovesBb)
}

func GeneratePawnMoves(pos *position.Position, mvsList *moves.Moves, ht *hashtables.HashTables) {
	var getShiftedBb func(*bitboard.Bitboard, uint) *bitboard.Bitboard
	var directionOfMovement int
	var doublePushMask *bitboard.Bitboard
	var promotionRank *bitboard.Bitboard
	var attackRightShift uint
	var attackLeftShift uint
	enPassanteBB := bitboard.NewBitboardFromIndex(pos.EnPassante())
	if pos.GetActiveSide() == position.White {
		getShiftedBb = bitboard.GetShiftedRightBb
		directionOfMovement = 1
		doublePushMask, _ = bitboard.NewBitboard(ht.FourthRankBb)
		promotionRank, _ = bitboard.NewBitboard(ht.EighthRankBb)
		attackRightShift = 7
		attackLeftShift = 9
	} else {
		getShiftedBb = bitboard.GetShiftedLeftBb
		directionOfMovement = -1
		doublePushMask, _ = bitboard.NewBitboard(ht.FifthRankBb)
		promotionRank, _ = bitboard.NewBitboard(ht.FirstRankBb)
		attackRightShift = 9
		attackLeftShift = 7
	}

	pawnPosBb := pos.GetActiveSidesBitboards()[position.Pawns]
	hFileBb, _ := bitboard.NewBitboard(ht.HfileBb)
	aFileBb, _ := bitboard.NewBitboard(ht.AfileBb)

	pawnAttackBb := getShiftedBb(&pawnPosBb, attackLeftShift).
		RemoveOverlappingBits(hFileBb).
		BitwiseAnd(bitboard.ReturnCombined(pos.InactiveSideOccupiedSqsBb(), enPassanteBB))
	addPawnMovesToArray(mvsList, int(attackLeftShift), directionOfMovement, pawnAttackBb, promotionRank)

	pawnAttackBb = getShiftedBb(&pawnPosBb, attackRightShift).
		RemoveOverlappingBits(aFileBb).
		BitwiseAnd(bitboard.ReturnCombined(pos.InactiveSideOccupiedSqsBb(), enPassanteBB))
	addPawnMovesToArray(mvsList, int(attackRightShift), directionOfMovement, pawnAttackBb, promotionRank)

	pawnPushBb := getShiftedBb(&pawnPosBb, 8)
	pawnPushBb.RemoveOverlappingBits(pos.AllOccupiedSqsBb())
	doubleRankpawnPushBb := getShiftedBb(pawnPushBb, 8)
	addPawnMovesToArray(mvsList, 8, directionOfMovement, pawnPushBb, promotionRank)
	doubleRankpawnPushBb.BitwiseAnd(doublePushMask)
	doubleRankpawnPushBb.RemoveOverlappingBits(pos.AllOccupiedSqsBb())
	addPawnMovesToArray(mvsList, 16, directionOfMovement, doubleRankpawnPushBb, promotionRank)
}

func addPawnMovesToArray(movesList *moves.Moves, shift int, shiftDirection int, pawnPushBb *bitboard.Bitboard, promoRank *bitboard.Bitboard) {
	shift = shift * shiftDirection
	for pawnPushBb.Value() != 0 {
		dest := pawnPushBb.Lsb()
		pawnPushBb.RemoveBit(dest)
		origin := dest + shift
		destBb := new(bitboard.Bitboard)
		destBb.SetBit(dest)
		if destBb.BitwiseAnd(promoRank).Value() != uint64(0) {
			movesList.AddPromotionMove(origin, dest, position.Queen)
			movesList.AddPromotionMove(origin, dest, position.Rooks)
			movesList.AddPromotionMove(origin, dest, position.Knights)
			movesList.AddPromotionMove(origin, dest, position.Bishops)
		} else {
			movesList.AddMove(origin, dest)
		}
	}
}

// AddValidMovesToArray save subset of valid moves from current position
func addValidMovesToArray(movesList *moves.Moves, index int, validMovesBb *bitboard.Bitboard) {
	var validMove int
	for validMovesBb.Value() != 0 {
		validMove = validMovesBb.Lsb()
		validMovesBb.RemoveBit(validMove)
		movesList.AddMove(index, validMove)
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
				generateValidDirectionalMovesBb(index, ht.SouthWestArrayBbHash, occSqsBb, getLsb))))
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

func generateValidDirectionalKingMovesBb(
	index int, directionalHash [64]uint64, occupiedSqsBb uint64, sigBit func(*bitboard.Bitboard) int) *bitboard.Bitboard {
	var validDirectionalMoves *bitboard.Bitboard
	occupiedSqsOverlapsDirectionalArray, _ := bitboard.NewBitboard(occupiedSqsBb & directionalHash[index])
	if occupiedSqsOverlapsDirectionalArray.Value() != 0 {
		sigBit := sigBit(occupiedSqsOverlapsDirectionalArray)
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
