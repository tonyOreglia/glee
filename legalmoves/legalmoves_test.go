package legalmoves

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyoreglia/glee/bitboard"
	"github.com/tonyoreglia/glee/hashtables"
	"github.com/tonyoreglia/glee/position"
)

// func TestGenerateLegalMoves(t *testing.T) {
// 	ht := hashtables.CalculateAllLookupBbs()
// 	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
// 	moves := NewLegalMoves(pos, ht)
// 	var expectedLegalMoves [][2]int
// 	expectedLegalMoves = make([][2]int, 0, 100)
// 	expectedLegalMoves = append(expectedLegalMoves, [2]int{1, 2})
// 	assert.Equal(t, expectedLegalMoves, moves)
// }

var ht = hashtables.CalculateAllLookupBbs()

func TestGenerateBishopMovesBb(t *testing.T) {
	// Legal white bishop moves from h1 unblocked
	pos, _ := position.NewPositionFen("7k/8/8/8/8/8/8/6KB w KQkq - 0 1")
	mvs := NewLegalMoves(pos, ht)
	mvs.generateBishopMoves()
	expectedMvs := [][2]int{{63, 0}, {63, 9}, {63, 18}, {63, 27}, {63, 36}, {63, 45}, {63, 54}}
	assert.Equal(t, expectedMvs, mvs.moves)

	// should be zero legal Bishop moves from starting position
	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	mvs = NewLegalMoves(pos, ht)
	mvs.generateBishopMoves()
	assert.Equal(t, 0, len(mvs.moves))
}

func TestGenerateValidDirectionalMovesBb(t *testing.T) {
	// piece sliding southwest from h8 unblocked
	sw := ht.SouthWestArrayBbHash
	index := 7
	validMvsBb := generateValidDirectionalMovesBb(index, sw, uint64(0))
	validMvsBb.Print()
	expectedValidMvsBb := uint64(0x102040810204000)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())

	// piece sliding north from a3, blocked on a7
	north := ht.NorthArrayBbHash
	index = 40
	validMvsBb = generateValidDirectionalMovesBb(index, north, uint64(0x100))
	validMvsBb.Print()
	expectedValidMvsBb = uint64(0x101010100)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())
}

func TestGenerateValidDiagonalSlidingMovesBb(t *testing.T) {
	// piece sliding diagonally from e4 unblocked
	index := 36
	occSqsVal := uint64(0)
	validMvsBb := generateValidDiagonalSlidingMovesBb(index, occSqsVal, ht)
	expectedValidMvsBb := uint64(0x8040200028448201)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())

	// piece sliding diagonally from a3 blocked at c5
	index = 40
	occSqsBb, _ := bitboard.NewBitboard(0)
	occSqsBb.SetBit(26)
	validMvsBb = generateValidDiagonalSlidingMovesBb(index, occSqsBb.Value(), ht)
	expectedValidMvsBb = uint64(0x402000204000000)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())
}
