package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/moves"
	"github.com/tonyOreglia/glee/pkg/position"
)

var ht = hashtables.Lookup

func TestGeneratePsuedoLegalMoves(t *testing.T) {
	tests := map[string]struct {
		pos           string
		generateMoves func(*position.Position) *moves.Moves
		assertion     func(*moves.Moves, string)
	}{
		// ## ALL PIECES ##
		"starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				return GenerateMoves(pos)
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 20, mvs.Length(), msg)
			},
		},
		"this is 25 because a pseudo legal move is included which leaves the king in check": {
			pos: "n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				return GenerateMoves(pos)
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 25, mvs.Length(), msg)
			},
		},
		"knight moves": {
			pos: "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				return GenerateMoves(pos)
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 48, mvs.Length(), msg)
			},
		},
		// ## PAWN MOVES ##
		"pawn moves from starting postion": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GeneratePawnMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 16, mvs.Length(), msg)
			},
		},
		"white pawn attacks diagonal left and moves forward one": {
			pos: "k7/8/r7/1P6/8/8/8/K7 w - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GeneratePawnMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 2, mvs.Length(), msg)
			},
		},
		"black pawn attacks diagonal left and moves forward one": {
			pos: "k7/8/r7/1p6/R7/8/8/K7 b - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GeneratePawnMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 2, mvs.Length(), msg)
			},
		},
		"white pawns reaches final rank and has four options for promotion": {
			pos: "k7/7P/8/8/8/8/8/K7 w - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GeneratePawnMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 4, mvs.Length(), msg)
			},
		},
		"black pawns reaches final rank and has four options for promotion": {
			pos: "k7/8/8/8/8/8/7p/K7 b - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GeneratePawnMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 4, mvs.Length(), msg)
			},
		},
		"white pawn blocked from moving to final rank": {
			pos: "r3k2r/ppppPppp/8/8/8/8/8/K7 w - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GeneratePawnMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 0, mvs.Length(), msg)
			},
		},
		"black pawn blocked from reaching final rank": {
			pos: "k7/8/8/8/8/8/PPPPpPPP/R3K2R b - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GeneratePawnMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 0, mvs.Length(), msg)
			},
		},
		"black pawns attacks en passante square": {
			pos: "k7/8/8/8/6pP/8/8/K7 b - h3 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GeneratePawnMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{38, 47}, {38, 46}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves())
			},
		},
		"white pawns attacks en passante square": {
			pos: "k7/8/8/pP6/8/8/8/K7 w - a6 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GeneratePawnMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{25, 17}, {25, 16}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves())
			},
		},
		// ## QUEEN MOVES ##
		"should be zero legal queen moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateQueenMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 0, mvs.Length(), msg)
			},
		},
		"Legal queen moves from a1 blocked horizontally": {
			pos: "7k/8/8/8/8/8/8/QK6 w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateQueenMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{56, 0}, {56, 7}, {56, 8}, {56, 14}, {56, 16},
					{56, 21}, {56, 24}, {56, 28}, {56, 32}, {56, 35}, {56, 40}, {56, 42},
					{56, 48}, {56, 49}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves())
			},
		},
		// ## KING MOVES ##
		"no white king moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 0, mvs.Length(), msg)
			},
		},
		"no black king moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 2",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 0, mvs.Length(), msg)
			},
		},
		"white king middle of board unblocked": {
			pos: "7k/8/8/8/3K4/8/8/3B4 w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{35, 26}, {35, 27}, {35, 28}, {35, 34}, {35, 36}, {35, 42}, {35, 43}, {35, 44}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"black king middle of board unblocked": {
			pos: "7K/8/8/8/3k4/8/8/3B4 b - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{35, 26}, {35, 27}, {35, 28}, {35, 34}, {35, 36}, {35, 42}, {35, 43}, {35, 44}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"white king middle of board completely blocked": {
			pos: "7k/8/8/PPPPPPPP/RRRKRRRR/RRRRRRRR/8/3B4 w - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"black king middle of board completely blocked": {
			pos: "8/8/8/pppppppp/rrrkrrrr/rrrrrrrr/8/3B4 b - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"white king middle of board surrounded by opposition": {
			pos: "7k/8/8/pppppppp/rrrKrrrr/rrrrrrrr/8/3B4 w - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{35, 26}, {35, 27}, {35, 28}, {35, 34}, {35, 36}, {35, 42}, {35, 43}, {35, 44}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"black king middle of board surrounded by opposition": {
			pos: "8/8/8/PPPPPPPP/RRRkRRRR/RRRRRRRR/8/3B4 b - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{35, 26}, {35, 27}, {35, 28}, {35, 34}, {35, 36}, {35, 42}, {35, 43}, {35, 44}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"white castling king-side": {
			pos: "7k/8/8/8/8/8/PPPPPPPP/3QK3 w K - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{60, 61}, {60, 62}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"black castling king-side": {
			pos: "3rk3/pppppppp/8/8/8/8/PPPPPPPP/3QK3 b kq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{4, 6}, {4, 5}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"white castling queen-side": {
			pos: "7k/8/8/8/8/8/PPPPPPPP/R3KQ2 w KQ - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				mvs.Print()
				expectedMvs := [][]int{{60, 59}, {60, 58}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"black castling queen-side": {
			pos: "r3kr2/pppppppp/8/8/8/8/PPPPPPPP/3QK3 b kq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{4, 3}, {4, 2}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"white king w/o castling permission": {
			pos: "7k/8/8/8/8/8/PPPPPPPP/R3K2R w - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{60, 61}, {60, 59}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"black w/o castling permission": {
			pos: "r3k2r/ppppppppp/8/8/8/8/PPPPPPPP/3QK3 b - - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{4, 3}, {4, 5}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"white castling both-sides": {
			pos: "7k/8/8/8/8/8/PPPPPPPP/R3K2R w KQ - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{60, 59}, {60, 58}, {60, 61}, {60, 62}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"black castling both-sides": {
			pos: "r3k2r/ppppppppp/8/8/8/8/PPPPPPPP/3QK3 b KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{4, 3}, {4, 5}, {4, 2}, {4, 6}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"white has castling rights but blocked by own pieces both sides": {
			pos: "7k/8/8/8/8/8/PPPPPPPP/R2QKQ1R w KQ - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"black has castling rights but blocked by own pieces both sides": {
			pos: "r2qkq1r/ppppppppp/8/8/8/8/PPPPPPPP/3QK3 b KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"white has castling rights but blocked by opposition pieces both sides": {
			pos: "7k/8/8/8/8/8/PPPPPPPP/R2qKq1R w KQ - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{60, 61}, {60, 59}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"black has castling rights but blocked by opposition pieces both sides": {
			pos: "r2QkQ1r/ppppppppp/8/8/8/8/PPPPPPPP/3QK3 b KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{4, 3}, {4, 5}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"white has castling rights but blocked by opposition pieces on both sides with a space to move": {
			pos: "7k/8/8/8/8/8/PPPPPPPP/R1q1K1qR w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{60, 61}, {60, 59}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"black has castling rights but blocked by opposition pieces on both sides with a space to move": {
			pos: "r1Q1k1Qr/ppppppppp/8/8/8/8/PPPPPPPP/3QK3 b kq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{4, 3}, {4, 5}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		"king moves castling and an attack an seventh rank": {
			pos: "r3k2r/5N2/8/8/8/8/8/7K b KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 7, mvs.Length(), msg)
			},
		},
		"pseudo legal moves allow king to castle out of check from pawn": {
			pos: "r3k2r/ppppPppp/8/8/8/8/8/7K b kq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKingMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{4, 12}, {4, 5}, {4, 3}, {4, 6}, {4, 2}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves(), msg)
			},
		},
		// ## BISHOP MOVES ##
		"should be zero legal Bishop moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateBishopMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 0, mvs.Length(), msg)
			},
		},
		"bishop  moves wide open spaces": {
			pos: "7k/8/8/8/8/8/8/6KB w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateBishopMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{63, 0}, {63, 9}, {63, 18}, {63, 27}, {63, 36}, {63, 45}, {63, 54}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves())
			},
		},
		"legal white bishop moves blocked on right array": {
			pos: "7k/8/8/8/8/5r2/8/3B3K w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateBishopMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{59, 32}, {59, 41}, {59, 45}, {59, 50}, {59, 52}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves())
			},
		},
		// ## ROOK MOVES ##
		"should be zero legal white Rook moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateRookMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 0, mvs.Length(), msg)
			},
		},
		"should be zero legal black Rook moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b KQkq - 0 2",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateRookMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 0, mvs.Length(), msg)
			},
		},
		"Legal rook moves from a1 blocked horizontally": {
			pos: "7k/8/8/8/8/8/8/RK6 w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateRookMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{56, 0}, {56, 16}, {56, 8}, {56, 24}, {56, 32}, {56, 40}, {56, 48}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves())
			},
		},
		"legal rook moves from d4 unblocked": {
			pos: "7k/8/8/8/3R4/8/8/7K w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateRookMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{35, 3}, {35, 11}, {35, 19}, {35, 27}, {35, 32}, {35, 33}, {35, 34}, {35, 36}, {35, 37}, {35, 38}, {35, 39}, {35, 43}, {35, 51}, {35, 59}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves())
			},
		},
		// ## KNIGHT MOVES ##
		"should be four legal knight moves from starting position": {
			pos: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKnightMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				assert.Equal(t, 4, mvs.Length(), msg)
			},
		},
		"Legal knight moves from a1 blocked at c2 as white": {
			pos: "7k/8/8/8/8/8/2B5/NK6 w KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKnightMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{56, 41}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves())
			},
		},
		"Legal knight moves from a1 blocked at c2 as black": {
			pos: "7k/8/8/8/8/1B6/2b5/nK6 b KQkq - 0 1",
			generateMoves: func(pos *position.Position) *moves.Moves {
				mvs := moves.NewMovesList()
				GenerateKnightMoves(pos, mvs, ht)
				return mvs
			},
			assertion: func(mvs *moves.Moves, msg string) {
				expectedMvs := [][]int{{56, 41}}
				assert.ElementsMatch(t, expectedMvs, mvs.GetMoves())
			},
		},
	}
	for tName, test := range tests {
		pos, _ := position.NewPositionFen(test.pos)
		legalMoves := test.generateMoves(pos)
		test.assertion(legalMoves, tName)
	}
}
