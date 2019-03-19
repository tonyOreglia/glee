package evaluate

import (
	"github.com/tonyOreglia/glee/pkg/position"
)

func EvaluatePosition(pos *position.Position) int {
	score := 0

	whiteKingBb := pos.GetWhiteBitboards()[position.King]
	whiteQueenBb := pos.GetWhiteBitboards()[position.Queen]
	whiteRooksBb := pos.GetWhiteBitboards()[position.Rooks]
	whiteKnightsBb := pos.GetWhiteBitboards()[position.Knights]
	whiteBishopsBb := pos.GetWhiteBitboards()[position.Bishops]
	whitePawnsBb := pos.GetWhiteBitboards()[position.Pawns]
	blackKingBb := pos.GetBlackBitboards()[position.King]
	blackQueensBb := pos.GetBlackBitboards()[position.Queen]
	blackRooksBb := pos.GetBlackBitboards()[position.Rooks]
	blackKnightsBb := pos.GetBlackBitboards()[position.Knights]
	blackBishopsBb := pos.GetBlackBitboards()[position.Bishops]
	blackPawnsBb := pos.GetBlackBitboards()[position.Pawns]

	score += 20000 * (whiteKingBb.PopulationCount() - blackKingBb.PopulationCount())
	score += 510 * (whiteRooksBb.PopulationCount() - blackRooksBb.PopulationCount())
	score += 320 * (whiteKnightsBb.PopulationCount() - blackKnightsBb.PopulationCount())
	score += 330 * (whiteBishopsBb.PopulationCount() - blackBishopsBb.PopulationCount())
	score += 100 * (whitePawnsBb.PopulationCount() - blackPawnsBb.PopulationCount())
	score += 890 * (whiteQueenBb.PopulationCount() - blackQueensBb.PopulationCount())

	score += kingBonusWhite[whiteKingBb.Msb()]
	score -= kingBonusBlack[blackKingBb.Msb()]

	if whiteBishopsBb.PopulationCount() > 1 {
		score += 15
	}

	if blackBishopsBb.PopulationCount() > 1 {
		score -= 15
	}

	for !whiteBishopsBb.IsZero() {
		msb := whiteBishopsBb.Msb()
		score += bishopBonusWhite[msb]
		whiteBishopsBb.RemoveBit(msb)
	}
	for !blackBishopsBb.IsZero() {
		msb := blackBishopsBb.Msb()
		score -= bishopBonusBlack[msb]
		blackBishopsBb.RemoveBit(msb)
	}

	for !whiteKnightsBb.IsZero() {
		msb := whiteKnightsBb.Msb()
		score += knightBonus[msb]
		whiteKnightsBb.RemoveBit(msb)
	}
	for !blackKnightsBb.IsZero() {
		msb := blackKnightsBb.Msb()
		score -= knightBonus[msb]
		blackKnightsBb.RemoveBit(msb)
	}

	for !whitePawnsBb.IsZero() {
		msb := whitePawnsBb.Msb()
		score += pawnBonusWhite[msb]
		whitePawnsBb.RemoveBit(msb)
	}
	for !blackPawnsBb.IsZero() {
		msb := blackPawnsBb.Msb()
		score -= pawnBonusBlack[msb]
		blackPawnsBb.RemoveBit(msb)
	}

	return score
}
