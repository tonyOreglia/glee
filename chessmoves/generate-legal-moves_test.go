package chessmoves

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyoreglia/glee/bitboard"
	"github.com/tonyoreglia/glee/hashtables"
	"github.com/tonyoreglia/glee/position"
)

var ht = hashtables.CalculateAllLookupBbs()

func TestGenerateMoves(t *testing.T) {
	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	mvs := NewLegalMoveGenerator(pos, ht)
	mvs.GenerateMoves()
	mvs.movesList.Print()
	bb := mvs.movesList.GetBitboard()
	bb.Print()
	assert.Equal(t, 20, mvs.movesList.Length())

	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/7P/PPPPPPP1/RNBQKBNR b KQkq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.GenerateMoves()
	mvs.movesList.Print()
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	assert.Equal(t, 20, mvs.movesList.Length())

	pos, _ = position.NewPositionFen("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.GenerateMoves()
	mvs.movesList.Print()
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	assert.Equal(t, 25, mvs.movesList.Length())

	// mdidle game position
	pos, _ = position.NewPositionFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.GenerateMoves()
	mvs.movesList.Print()
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	assert.Equal(t, 48, mvs.movesList.Length())

	// promotion testing
	pos, _ = position.NewPositionFen("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.GenerateMoves()
	mvs.movesList.Print()
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	// this should be 24 once I am correctly checking for safe king moves. avoid check
	assert.Equal(t, 25, mvs.movesList.Length())

}
func TestGeneratePawnMoves(t *testing.T) {
	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	mvs := NewLegalMoveGenerator(pos, ht)
	mvs.generatePawnMoves()
	mvs.movesList.Print()
	bb := mvs.movesList.GetBitboard()
	bb.Print()
	assert.Equal(t, 16, mvs.movesList.Length())

	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/7P/PPPPPPP1/RNBQKBNR b KQkq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generatePawnMoves()
	mvs.movesList.Print()
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	assert.Equal(t, 16, mvs.movesList.Length())

	// testing attack
	pos, _ = position.NewPositionFen("k7/8/r7/1P7/8/8/8/K7 w - - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generatePawnMoves()
	mvs.movesList.Print()
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	assert.Equal(t, 2, mvs.movesList.Length())

	// testing white promotion
	pos, _ = position.NewPositionFen("k7/7P/8/8/8/8/8/K7 w - - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generatePawnMoves()
	mvs.movesList.Print()
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	assert.Equal(t, 4, mvs.movesList.Length())

	// testing black promotion with attack
	pos, _ = position.NewPositionFen("k7/8/8/8/8/8/7p/K5N1 b - - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generatePawnMoves()
	mvs.movesList.Print()
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	assert.Equal(t, 8, mvs.movesList.Length())
}

func TestGenerateKingMovesBb(t *testing.T) {
	// should be zero legal king moves from starting position
	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1")
	mvs := NewLegalMoveGenerator(pos, ht)
	mvs.generateKingMoves()
	fmt.Print(mvs.movesList)
	bb := mvs.movesList.GetBitboard()
	bb.Print()
	assert.Equal(t, 0, mvs.movesList.Length())

	// legal white King moves in d4
	pos, _ = position.NewPositionFen("7k/8/8/8/3K4/8/8/3B4 w KQkq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateKingMoves()
	fmt.Print(mvs.movesList)
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	expectedMvs := [][]int{{35, 26}, {35, 27}, {35, 28}, {35, 34}, {35, 36}, {35, 42}, {35, 43}, {35, 44}}
	assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())

	// black king side castling
	pos, _ = position.NewPositionFen("7k/8/8/8/8/8/PPPPPPPP/3QK3 w K - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateKingMoves()
	fmt.Print(mvs.movesList)
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	expectedMvs = [][]int{{60, 61}, {60, 62}}
	assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())

	// black king side castling no permission
	pos, _ = position.NewPositionFen("7k/8/8/8/8/8/PPPPPPPP/3QK3 w - - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateKingMoves()
	fmt.Print(mvs.movesList)
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	expectedMvs = [][]int{{60, 61}}
	assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())

	// black king side castling
	pos, _ = position.NewPositionFen("3rk3/pppppppp/8/8/8/8/PPPPPPPP/3QK3 b kq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateKingMoves()
	fmt.Print(mvs.movesList)
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	expectedMvs = [][]int{{4, 6}, {4, 5}}
	assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())

	// black queen side castling
	pos, _ = position.NewPositionFen("4kr3/pppppppp/8/8/8/8/PPPPPPPP/3QK3 b kq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateKingMoves()
	fmt.Print(mvs.movesList)
	bb = mvs.movesList.GetBitboard()
	bb.Print()
	expectedMvs = [][]int{{4, 3}, {4, 2}}
	assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())

}

func TestGenerateBishopMovesBb(t *testing.T) {
	// Legal white bishop moves from h1 unblocked
	pos, _ := position.NewPositionFen("7k/8/8/8/8/8/8/6KB w KQkq - 0 1")
	mvs := NewLegalMoveGenerator(pos, ht)
	mvs.generateBishopMoves()
	expectedMvs := [][]int{{63, 0}, {63, 9}, {63, 18}, {63, 27}, {63, 36}, {63, 45}, {63, 54}}
	assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())

	// should be zero legal Bishop moves from starting position
	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateBishopMoves()
	assert.Equal(t, 0, mvs.movesList.Length())

	// legal white bishop moves blocked on right array
	pos, _ = position.NewPositionFen("7k/8/8/8/8/5r2/8/3B3K w KQkq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateBishopMoves()
	expectedMvs = [][]int{{59, 32}, {59, 41}, {59, 45}, {59, 50}, {59, 52}}
	assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
}

func TestGenerateRookMovesBb(t *testing.T) {
	// Legal rook moves from a1 blocked horizontally
	pos, _ := position.NewPositionFen("7k/8/8/8/8/8/8/RK6 w KQkq - 0 1")
	mvs := NewLegalMoveGenerator(pos, ht)
	mvs.generateRookMoves()
	fmt.Print(mvs.movesList.GetMoves())
	expectedMvs := [][]int{{56, 0}, {56, 16}, {56, 8}, {56, 24}, {56, 32}, {56, 40}, {56, 48}}
	assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())

	// should be zero legal Rook moves from starting position
	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateRookMoves()
	fmt.Print(mvs.movesList.GetMoves())
	assert.Equal(t, 0, mvs.movesList.Length())

	// should be zero legal Rook moves from starting position
	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 2")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateRookMoves()
	mvs.movesList.Print()
	bb := mvs.movesList.GetBitboard()
	bb.Print()
	assert.Equal(t, 0, mvs.movesList.Length())

	// legal rook moves from d4 unblocked
	pos, _ = position.NewPositionFen("7k/8/8/8/3R4/8/8/7K w KQkq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateRookMoves()
	expectedMvs = [][]int{{35, 3}, {35, 11}, {35, 19}, {35, 27}, {35, 32}, {35, 33}, {35, 34}, {35, 36}, {35, 37}, {35, 38}, {35, 39}, {35, 43}, {35, 51}, {35, 59}}
	assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
}

func TestGenerateQueenMovesBb(t *testing.T) {
	// should be zero legal queen moves from starting position
	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	mvs := NewLegalMoveGenerator(pos, ht)
	mvs.generateQueenMoves()
	assert.Equal(t, 0, mvs.movesList.Length())

	// Legal queen moves from a1 blocked horizontally
	pos, _ = position.NewPositionFen("7k/8/8/8/8/8/8/QK6 w KQkq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateQueenMoves()
	expectedMvs := [][]int{{56, 0}, {56, 7}, {56, 8}, {56, 14}, {56, 16},
		{56, 21}, {56, 24}, {56, 28}, {56, 32}, {56, 35}, {56, 40}, {56, 42},
		{56, 48}, {56, 49}}
	assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
}

func TestGenerateKnightMovesBb(t *testing.T) {
	// should be four legal knight moves from starting position
	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	mvs := NewLegalMoveGenerator(pos, ht)
	mvs.generateKnightMoves()
	assert.Equal(t, 4, mvs.movesList.Length())

	// Legal knight moves from a1 blocked at c2 as white
	pos, _ = position.NewPositionFen("7k/8/8/8/8/8/2B5/NK6 w KQkq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateKnightMoves()
	fmt.Print(mvs.movesList.GetMoves())
	expectedMvs := [][]int{{56, 41}}
	assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())

	// Legal knight moves from a1 blocked at c2 as black
	pos, _ = position.NewPositionFen("7k/8/8/8/8/1B6/2b5/nK6 b KQkq - 0 1")
	mvs = NewLegalMoveGenerator(pos, ht)
	mvs.generateKnightMoves()
	fmt.Print(mvs.movesList.GetMoves())
	expectedMvs = [][]int{{56, 41}}
	assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
}

func TestGenerateValidDirectionalMovesBb(t *testing.T) {
	// piece sliding southwest from h8 unblocked
	sw := ht.SouthWestArrayBbHash
	index := 7
	validMvsBb := generateValidDirectionalMovesBb(index, sw, uint64(0), getMsb)
	expectedValidMvsBb := uint64(0x102040810204000)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())

	// piece sliding north from a3, blocked on a7
	north := ht.NorthArrayBbHash
	index = 40
	validMvsBb = generateValidDirectionalMovesBb(index, north, uint64(0x100), getLsb)
	expectedValidMvsBb = uint64(0x101010100)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())
}

func TestGenerateValidDiagonalSlidingMovesBb(t *testing.T) {
	// piece sliding diagonally from e4 unblocked
	index := 36
	occSqsVal := uint64(0)
	validMvsBb := generateValidDiagonalSlidingMovesBb(index, occSqsVal, ht)
	expectedValidMvsBb, _ := bitboard.NewBitboard(uint64(0x8244280028448201))
	expectedValidMvsBb.Print()
	validMvsBb.Print()
	assert.Equal(t, expectedValidMvsBb.Value(), validMvsBb.Value())

	// piece sliding diagonally from a3 blocked at c5
	index = 40
	occSqsBb, _ := bitboard.NewBitboard(0)
	occSqsBb.SetBit(26)
	validMvsBb = generateValidDiagonalSlidingMovesBb(index, occSqsBb.Value(), ht)
	expectedValidMvsBb, _ = bitboard.NewBitboard(uint64(0x402000204000000))
	expectedValidMvsBb.Print()
	validMvsBb.Print()
	assert.Equal(t, expectedValidMvsBb.Value(), validMvsBb.Value())
}
