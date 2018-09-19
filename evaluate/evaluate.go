package evaluate

import (
	"github.com/tonyoreglia/glee/position"
)

func EvaluatePosition(pos *position.Position, perft *int) int {
	// pos.Print()
	score := 0
	*perft++

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

	return score
}
