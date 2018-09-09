package legalmoves

import (
	"fmt"
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

func TestGenerateKingMovesBb(t *testing.T) {
	// should be zero legal king moves from starting position
	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	mvs := NewLegalMoves(pos, ht)
	mvs.generateKingMoves()
	assert.Equal(t, 0, len(mvs.moves))

	// legal white King moves in d4
	pos, _ = position.NewPositionFen("7k/8/8/8/3K4/8/8/3B4 w KQkq - 0 1")
	mvs = NewLegalMoves(pos, ht)
	mvs.generateKingMoves()
	fmt.Print(mvs.moves)
	bb := bitboard.NewBbFromMovesSlice(mvs.moves)
	bb.Print()
	expectedMvs := [][2]int{{35, 26}, {35, 27}, {35, 28}, {35, 34}, {35, 36}, {35, 42}, {35, 43}, {35, 44}}
	assert.ElementsMatch(t, expectedMvs, mvs.moves)

	// king side castling
	pos, _ = position.NewPositionFen("7k/8/8/8/8/8/PPPPPPPP/3QK3 w K - 0 1")
	mvs = NewLegalMoves(pos, ht)
	mvs.generateKingMoves()
	fmt.Print(mvs.moves)
	bb = bitboard.NewBbFromMovesSlice(mvs.moves)
	bb.Print()
	expectedMvs = [][2]int{{60, 61}, {60, 62}}
	assert.ElementsMatch(t, expectedMvs, mvs.moves)
}

func TestGenerateBishopMovesBb(t *testing.T) {
	// Legal white bishop moves from h1 unblocked
	pos, _ := position.NewPositionFen("7k/8/8/8/8/8/8/6KB w KQkq - 0 1")
	mvs := NewLegalMoves(pos, ht)
	mvs.generateBishopMoves()
	expectedMvs := [][2]int{{63, 0}, {63, 9}, {63, 18}, {63, 27}, {63, 36}, {63, 45}, {63, 54}}
	assert.ElementsMatch(t, expectedMvs, mvs.moves)

	// should be zero legal Bishop moves from starting position
	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	mvs = NewLegalMoves(pos, ht)
	mvs.generateBishopMoves()
	assert.Equal(t, 0, len(mvs.moves))

	// legal white bishop moves blocked on right array
	pos, _ = position.NewPositionFen("7k/8/8/8/8/5r2/8/3B3K w KQkq - 0 1")
	mvs = NewLegalMoves(pos, ht)
	mvs.generateBishopMoves()
	expectedMvs = [][2]int{{59, 32}, {59, 41}, {59, 45}, {59, 50}, {59, 52}}
	assert.ElementsMatch(t, expectedMvs, mvs.moves)
}

func TestGenerateRookMovesBb(t *testing.T) {
	// Legal rook moves from a1 blocked horizontally
	pos, _ := position.NewPositionFen("7k/8/8/8/8/8/8/RK6 w KQkq - 0 1")
	mvs := NewLegalMoves(pos, ht)
	mvs.generateRookMoves()
	fmt.Print(mvs.moves)
	expectedMvs := [][2]int{{56, 0}, {56, 16}, {56, 8}, {56, 24}, {56, 32}, {56, 40}, {56, 48}}
	assert.ElementsMatch(t, expectedMvs, mvs.moves)

	// should be zero legal Rook moves from starting position
	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	mvs = NewLegalMoves(pos, ht)
	mvs.generateRookMoves()
	assert.Equal(t, 0, len(mvs.moves))

	// legal rook moves from d4 unblocked
	pos, _ = position.NewPositionFen("7k/8/8/8/3R4/8/8/7K w KQkq - 0 1")
	mvs = NewLegalMoves(pos, ht)
	mvs.generateRookMoves()
	expectedMvs = [][2]int{{35, 3}, {35, 11}, {35, 19}, {35, 27}, {35, 32}, {35, 33}, {35, 34}, {35, 36}, {35, 37}, {35, 38}, {35, 39}, {35, 43}, {35, 51}, {35, 59}}
	assert.ElementsMatch(t, expectedMvs, mvs.moves)
}

func TestGenerateQueenMovesBb(t *testing.T) {
	// should be zero legal queen moves from starting position
	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	mvs := NewLegalMoves(pos, ht)
	mvs.generateQueenMoves()
	assert.Equal(t, 0, len(mvs.moves))

	// Legal queen moves from a1 blocked horizontally
	pos, _ = position.NewPositionFen("7k/8/8/8/8/8/8/QK6 w KQkq - 0 1")
	mvs = NewLegalMoves(pos, ht)
	mvs.generateQueenMoves()
	expectedMvs := [][2]int{{56, 0}, {56, 7}, {56, 8}, {56, 14}, {56, 16},
		{56, 21}, {56, 24}, {56, 28}, {56, 32}, {56, 35}, {56, 40}, {56, 42},
		{56, 48}, {56, 49}}
	assert.ElementsMatch(t, expectedMvs, mvs.moves)
}

func TestGenerateKnightMovesBb(t *testing.T) {
	// should be four legal knight moves from starting position
	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	mvs := NewLegalMoves(pos, ht)
	mvs.generateKnightMoves()
	assert.Equal(t, 4, len(mvs.moves))

	// Legal knight moves from a1 blocked at c2 as white
	pos, _ = position.NewPositionFen("7k/8/8/8/8/8/2B5/NK6 w KQkq - 0 1")
	mvs = NewLegalMoves(pos, ht)
	mvs.generateKnightMoves()
	fmt.Print(mvs.moves)
	expectedMvs := [][2]int{{56, 41}}
	assert.ElementsMatch(t, expectedMvs, mvs.moves)

	// Legal knight moves from a1 blocked at c2 as black
	pos, _ = position.NewPositionFen("7k/8/8/8/8/1B6/2b5/nK6 b KQkq - 0 1")
	mvs = NewLegalMoves(pos, ht)
	mvs.generateKnightMoves()
	fmt.Print(mvs.moves)
	expectedMvs = [][2]int{{56, 41}}
	assert.ElementsMatch(t, expectedMvs, mvs.moves)
}

func TestGenerateValidDirectionalMovesBb(t *testing.T) {
	// piece sliding southwest from h8 unblocked
	sw := ht.SouthWestArrayBbHash
	index := 7
	validMvsBb := generateValidDirectionalMovesBb(index, sw, uint64(0))
	expectedValidMvsBb := uint64(0x102040810204000)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())

	// piece sliding north from a3, blocked on a7
	north := ht.NorthArrayBbHash
	index = 40
	validMvsBb = generateValidDirectionalMovesBb(index, north, uint64(0x100))
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
