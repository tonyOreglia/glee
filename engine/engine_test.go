package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyoreglia/glee/moves"
	"github.com/tonyoreglia/glee/position"
)

var flagtests = []struct {
	pos           string
	depth         int
	expectedNodes int
}{
	{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 1, 20},
	{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 2, 400},
	{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 1, 48},
	{"n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 1, 24},
	{"n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1", 2, 496},
	{"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1", 2, 2039},
	{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1", 3, 8092},
}

func setup() (int, int, *moves.Move) {
	perft := 0
	singlePlyPerft := 0
	mv := new(moves.Move)
	return perft, singlePlyPerft, mv
}

func TestMinMax(t *testing.T) {
	for _, tt := range flagtests {
		perft, singlePlyPerft, mv := setup()
		pos, _ := position.NewPositionFen(tt.pos)
		minMax(tt.depth, tt.depth, &pos, &mv, &perft, &singlePlyPerft)
		assert.Equal(t, tt.expectedNodes, perft)
	}
}
