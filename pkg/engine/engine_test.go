package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyOreglia/glee/pkg/moves"
	"github.com/tonyOreglia/glee/pkg/position"
)

var flagtests = []struct {
	fen           string
	depth         int
	expectedNodes int
	name          string
}{
	{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 1, 20, ""},
	{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 2, 400, ""},
	{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 3, 8902, ""},
	{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 4, 197281, ""},
	{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 1, 48, ""},
	{"n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 1, 24, ""},
	{"n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 2, 496, ""},
	{"r3k2r/p1ppqNb1/bn2pnp1/3P4/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq - 0 1", 4, 4164923, "test position 001"},
	{"r3k2r/p1ppqNb1/1n2pnp1/3P4/1p2P3/2N2Q1p/PPPBbPPP/R3K2R w KQkq - 0 1", 1, 41, "test position 001 --> a6e2"},
	{"1r2k2r/p1ppqNb1/bn2pnp1/3P4/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQ - 1 1", 1, 48, "test position 001 --> a8b8"},
	{"2r1k2r/p1ppqNb1/bn2pnp1/3P4/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQ - 1 1", 1, 48, "test position 001 --> a8c8"},
	{"r3k2r/p1ppqNb1/1n2pnp1/1b1P4/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 1 1", 2, 2084, "test position 001 --> a6b5"},
	{"r3k2r/p1ppqNb1/1n2pnp1/1b1P4/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 1 1", 3, 97761, "test position 001 --> a6b5"},
	{"r3k2r/p1ppqNb1/1n2pnp1/1b1P4/Pp2P3/2N2Q1p/1PPBBPPP/R3K2R b KQkq a3 0 1", 2, 2275, "test position 001 --> a6b5, a2a4"},
	{"r3k2r/p1ppqNb1/bn2pnp1/3P4/4P3/2p2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 2, 2120, "test position 001 --> b4c3"},
	{"r3k2r/p1ppqNb1/1n2pnp1/1b1P4/4P3/2p2Q1p/PPPBBPPP/1R2K2R w kq - 2 1", 2, 2073, "test position 002 (001 --> b4c3 --> a1b1 --> a6b5)"},
	{"r3k2r/p1ppqNb1/bn2pnp1/3P4/4P3/2p2Q1p/PPPBBPPP/1R2K2R b kq - 1 1", 2, 2039, "test position 001 --> b4c3 --> a1b1"},
	{"r3k2r/p1ppqNb1/bn2pnp1/3P4/4P3/2p2Q1p/PPPBBPPP/1R2K2R b kq - 1 1", 3, 90010, "test position 001 --> b4c3 --> a1b1"},
	{"r3k2r/p1ppqNb1/bn2pnp1/3P4/4P3/2p2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 3, 97799, "test position 001 --> b4c3"},
	{"r3k2r/p1ppqNb1/bn3np1/3p4/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 3, 99231, "test position 001 --> e6d5"},
	{"r3k2r/p1ppqNb1/bn3np1/3P4/1p6/2N2Q1p/PPPBBPPP/R3K2R b KQkq - 0 1", 2, 1951, "test position 001 --> e6d5 --> e4d5"},
	{"k7/8/6p1/8/8/7P/8/K7 b - - 0 1", 5, 4354, "test position 002"},
	{"n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 4, 182838, "test position 003"},
	{"R6r/8/8/2K5/5k2/8/8/r6R b - - 0 1", 4, 771368, "test position 004"},
}

func setup() (int, int) {
	perft := 0
	singlePlyPerft := 0
	return perft, singlePlyPerft
}

func TestMinMax(t *testing.T) {
	for _, tt := range flagtests {
		perft, singlePlyPerft := setup()
		pos, _ := position.NewPositionFen(tt.fen)
		MinMax(SearchParams{
			Depth:          tt.depth,
			Ply:            tt.depth,
			Pos:            &pos,
			Perft:          &perft,
			SinglePlyPerft: &singlePlyPerft,
		})
		assert.Equal(t, tt.expectedNodes, perft, tt.name)
	}
}

func TestMakeValidMove(t *testing.T) {
	tests := map[string]struct {
		move  *moves.Move
		pos   string
		legal bool
	}{
		"black cannot castle through check kingside": {
			move:  moves.NewMove([]int{4, 6}),
			pos:   "4k2r/8/8/8/8/8/5R2/7K b KQkq - 0 1",
			legal: false,
		},
		"white cannot castle through check kingside": {
			move:  moves.NewMove([]int{60, 62}),
			pos:   "7k/5r2/8/8/8/8/8/4K2R w K - 0 1",
			legal: false,
		},
		"black cannot castle through check by pawn kingside": {
			move:  moves.NewMove([]int{4, 6}),
			pos:   "4k2r/4P3/8/8/8/8/8/7K b KQkq - 0 1",
			legal: false,
		},
		"white cannot castle through check by pawn kingside": {
			move:  moves.NewMove([]int{60, 62}),
			pos:   "7k/8/8/8/8/8/4p3/4K2R w K - 0 1",
			legal: false,
		},
		"black cannot castle through check queenside": {
			move:  moves.NewMove([]int{4, 2}),
			pos:   "r3k2r/8/8/8/8/8/3R4/7K b q - 0 1",
			legal: false,
		},
		"white cannot castle through check queenside": {
			move:  moves.NewMove([]int{60, 58}),
			pos:   "7k/3r4/8/8/8/8/8/R3K2R w Q - 0 1",
			legal: false,
		},
		"black cannot castle through check by pawn queenside": {
			move:  moves.NewMove([]int{4, 2}),
			pos:   "r3k2r/4P3/8/8/8/8/8/7K b q - 0 1",
			legal: false,
		},
		"white cannot castle through check by pawn queenside": {
			move:  moves.NewMove([]int{60, 58}),
			pos:   "7k/8/8/8/8/8/4p3/R3K2R w Q - 0 1",
			legal: false,
		},
		"black can castle kingside": {
			move:  moves.NewMove([]int{4, 6}),
			pos:   "r3k2r/8/8/8/8/8/8/7K b q - 0 1",
			legal: true,
		},
		"white can castle kingside": {
			move:  moves.NewMove([]int{60, 62}),
			pos:   "7k/8/8/8/8/8/8/R3K2R w Q - 0 1",
			legal: true,
		},
		"black can castle queenside": {
			move:  moves.NewMove([]int{4, 2}),
			pos:   "r3k2r/8/8/8/8/8/8/7K b q - 0 1",
			legal: true,
		},
		"white can castle queenside": {
			move:  moves.NewMove([]int{60, 58}),
			pos:   "7k/8/8/8/8/8/8/R3K2R w Q - 0 1",
			legal: true,
		},
		"black cannot castle into check queenside": {
			move:  moves.NewMove([]int{4, 2}),
			pos:   "r3k2r/8/8/8/8/7B/8/7K b q - 0 1",
			legal: false,
		},
		"white cannot castle into check queenside": {
			move:  moves.NewMove([]int{60, 58}),
			pos:   "7k/8/7b/8/8/8/8/R3K2R w KQq - 0 1",
			legal: false,
		},
		"black cannot castle into check kingside": {
			move:  moves.NewMove([]int{4, 6}),
			pos:   "r3k2r/8/8/8/8/8/B7/7K b q - 0 1",
			legal: false,
		},
		"white cannot castle into check kingside": {
			move:  moves.NewMove([]int{60, 62}),
			pos:   "7k/8/8/8/8/5n2/8/R3K2R w KQq - 0 1",
			legal: false,
		},
		"black cannot castle out of check": {
			move:  moves.NewMove([]int{4, 6}),
			pos:   "r3k2r/8/8/4R3/8/8/8/7K b q - 0 1",
			legal: false,
		},
		"black king cannot castle kingside through check from pawn": {
			move:  moves.NewMove([]int{4, 6}),
			pos:   "r3k2r/ppppPppp/8/8/8/8/8/7K b kq - 0 1",
			legal: false,
		},
		"black king cannot castle queenside through check from pawn": {
			move:  moves.NewMove([]int{4, 2}),
			pos:   "r3k2r/pPpppppp/8/8/8/8/8/7K b kq - 0 1",
			legal: false,
		},
		"black king cannot castle queenside through check from pawn 2": {
			move:  moves.NewMove([]int{4, 2}),
			pos:   "r3k2r/ppPppppp/8/8/8/8/8/7K b kq - 0 1",
			legal: false,
		},
		"black king cannot castle queenside through check from pawn 3": {
			move:  moves.NewMove([]int{4, 2}),
			pos:   "r3k2r/pppppPpp/8/8/8/8/8/7K b kq - 0 1",
			legal: false,
		},
		"white king cannot castle kingside through check from pawn": {
			move:  moves.NewMove([]int{60, 62}),
			pos:   "7k/8/8/8/8/8/PPPPpPPP/R3K2R w KQ - 0 1",
			legal: false,
		},
	}
	for tName, test := range tests {
		pos, _ := position.NewPositionFen(test.pos)
		assert.Equal(t, test.legal, MakeValidMove(*test.move, &pos), tName)
	}
}
