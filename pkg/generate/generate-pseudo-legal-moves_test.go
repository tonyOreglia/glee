package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyOreglia/glee/pkg/bitboard"
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/position"
)

func TestGenerateMoves(t *testing.T) {
	tests := map[string]struct {
		pos           string
		generateMoves func(*LegalMoveGenerator)
		assertion     func(*LegalMoveGenerator)
	}{
		// ## ALL PIECES ##
		"starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.GenerateMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 20, mvs.movesList.Length())
			},
		},
		"this is 25 because a pseudo legal move is included which leaves the king in check": {
			pos: "n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.GenerateMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 25, mvs.movesList.Length())
			},
		},
		"knight moves": {
			pos: "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.GenerateMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 48, mvs.movesList.Length())
			},
		},
		// ## PAWN MOVES ##
		"pawn moves from starting postion": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generatePawnMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 16, mvs.movesList.Length())
			},
		},
		"white pawn attacks diagonal left and moves forward one": {
			pos: "k7/8/r7/1P7/8/8/8/K7 w - - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generatePawnMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 2, mvs.movesList.Length())
			},
		},
		"white pawns reaches final rank and has four options for promotion": {
			pos: "k7/7P/8/8/8/8/8/K7 w - - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generatePawnMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 4, mvs.movesList.Length())
			},
		},
		"black pawns reaches final rank and has four options for promotion": {
			pos: "k7/8/8/8/8/8/7p/K7 b - - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generatePawnMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 4, mvs.movesList.Length())
			},
		},
		// ## QUEEN MOVES ##
		"should be zero legal queen moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateQueenMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 0, mvs.movesList.Length())
			},
		},
		"Legal queen moves from a1 blocked horizontally": {
			pos: "7k/8/8/8/8/8/8/QK6 w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateQueenMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				expectedMvs := [][]int{{56, 0}, {56, 7}, {56, 8}, {56, 14}, {56, 16},
					{56, 21}, {56, 24}, {56, 28}, {56, 32}, {56, 35}, {56, 40}, {56, 42},
					{56, 48}, {56, 49}}
				assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
			},
		},
		// ## KING MOVES ##
		"should be zero legal king moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateKingMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 0, mvs.movesList.Length())
			},
		},
		"legal white King moves from middle of board unblocked": {
			pos: "7k/8/8/8/3K4/8/8/3B4 w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateKingMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				expectedMvs := [][]int{{35, 26}, {35, 27}, {35, 28}, {35, 34}, {35, 36}, {35, 42}, {35, 43}, {35, 44}}
				assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
			},
		},
		"white king castling king-side": {
			pos: "7k/8/8/8/8/8/PPPPPPPP/3QK3 w K - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateKingMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				expectedMvs := [][]int{{60, 61}, {60, 62}}
				assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
			},
		},
		"white king w/o castling permission": {
			pos: "7k/8/8/8/8/8/PPPPPPPP/3QK3 w - - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateKingMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				expectedMvs := [][]int{{60, 61}}
				assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
			},
		},
		"black king castling king-side": {
			pos: "3rk3/pppppppp/8/8/8/8/PPPPPPPP/3QK3 b kq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateKingMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				expectedMvs := [][]int{{4, 6}, {4, 5}}
				assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
			},
		},
		// ## BISHOP MOVES ##
		"should be zero legal Bishop moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateBishopMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 0, mvs.movesList.Length())
			},
		},
		"bishop  moves wide open spaces": {
			pos: "7k/8/8/8/8/8/8/6KB w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateBishopMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				expectedMvs := [][]int{{63, 0}, {63, 9}, {63, 18}, {63, 27}, {63, 36}, {63, 45}, {63, 54}}
				assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
			},
		},
		"legal white bishop moves blocked on right array": {
			pos: "7k/8/8/8/8/5r2/8/3B3K w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateBishopMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				expectedMvs := [][]int{{59, 32}, {59, 41}, {59, 45}, {59, 50}, {59, 52}}
				assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
			},
		},
		// ## ROOK MOVES ##
		"should be zero legal white Rook moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateRookMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 0, mvs.movesList.Length())
			},
		},
		"should be zero legal black Rook moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 2",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateRookMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 0, mvs.movesList.Length())
			},
		},
		"Legal rook moves from a1 blocked horizontally": {
			pos: "7k/8/8/8/8/8/8/RK6 w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateRookMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				expectedMvs := [][]int{{56, 0}, {56, 16}, {56, 8}, {56, 24}, {56, 32}, {56, 40}, {56, 48}}
				assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
			},
		},
		"legal rook moves from d4 unblocked": {
			pos: "7k/8/8/8/3R4/8/8/7K w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateRookMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				expectedMvs := [][]int{{35, 3}, {35, 11}, {35, 19}, {35, 27}, {35, 32}, {35, 33}, {35, 34}, {35, 36}, {35, 37}, {35, 38}, {35, 39}, {35, 43}, {35, 51}, {35, 59}}
				assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
			},
		},
		// ## KNIGHT MOVES ##
		"should be four legal knight moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateKnightMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				assert.Equal(t, 4, mvs.movesList.Length())
			},
		},
		"Legal knight moves from a1 blocked at c2 as white": {
			pos: "7k/8/8/8/8/8/2B5/NK6 w KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateKnightMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				expectedMvs := [][]int{{56, 41}}
				assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
			},
		},
		"Legal knight moves from a1 blocked at c2 as black": {
			pos: "7k/8/8/8/8/1B6/2b5/nK6 b KQkq - 0 1",
			generateMoves: func(mvs *LegalMoveGenerator) {
				mvs.generateKnightMoves()
			},
			assertion: func(mvs *LegalMoveGenerator) {
				expectedMvs := [][]int{{56, 41}}
				assert.ElementsMatch(t, expectedMvs, mvs.movesList.GetMoves())
			},
		},
	}
	for _, test := range tests {
		pos, _ := position.NewPositionFen(test.pos)
		mvs := NewLegalMoveGenerator(pos)
		test.generateMoves(mvs)
		mvs.movesList.Print()
		bb := mvs.movesList.GetBitboard()
		bb.Print()
		test.assertion(mvs)
	}
}

func TestGenerateValidDirectionalMovesBb(t *testing.T) {
	// piece sliding southwest from h8 unblocked
	sw := hashtables.Lookup.SouthWestArrayBbHash
	index := 7
	validMvsBb := generateValidDirectionalMovesBb(index, sw, uint64(0), getMsb)
	expectedValidMvsBb := uint64(0x102040810204000)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())

	// piece sliding north from a3, blocked on a7
	north := hashtables.Lookup.NorthArrayBbHash
	index = 40
	validMvsBb = generateValidDirectionalMovesBb(index, north, uint64(0x100), getLsb)
	expectedValidMvsBb = uint64(0x101010100)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())
}

func TestGenerateValidDiagonalSlidingMovesBb(t *testing.T) {
	// piece sliding diagonally from e4 unblocked
	index := 36
	occSqsVal := uint64(0)
	validMvsBb := generateValidDiagonalSlidingMovesBb(index, occSqsVal, hashtables.Lookup)
	expectedValidMvsBb, _ := bitboard.NewBitboard(uint64(0x8244280028448201))
	expectedValidMvsBb.Print()
	validMvsBb.Print()
	assert.Equal(t, expectedValidMvsBb.Value(), validMvsBb.Value())

	// piece sliding diagonally from a3 blocked at c5
	index = 40
	occSqsBb, _ := bitboard.NewBitboard(0)
	occSqsBb.SetBit(26)
	validMvsBb = generateValidDiagonalSlidingMovesBb(index, occSqsBb.Value(), hashtables.Lookup)
	expectedValidMvsBb, _ = bitboard.NewBitboard(uint64(0x402000204000000))
	expectedValidMvsBb.Print()
	validMvsBb.Print()
	assert.Equal(t, expectedValidMvsBb.Value(), validMvsBb.Value())
}
