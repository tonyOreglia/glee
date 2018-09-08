package legalmoves

type LegalMoves struct {
	// slice with capacity of 100 moves, starting len of 0
	moves [][2]int
}

func NewLegalMoves() *LegalMoves {
	legalMoves := &LegalMoves{}
	legalMoves.moves = make([][2]int, 0, 100)
	return legalMoves
}
