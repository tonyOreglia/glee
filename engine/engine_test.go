package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyoreglia/glee/moves"
	"github.com/tonyoreglia/glee/position"
)

func TestMinMax(t *testing.T) {
	// depth 1 starting position
	depth := 1
	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	perft := 0
	singlePlyPerft := 0
	mv := new(moves.Move)
	minMax(depth, depth, &pos, &mv, &perft, &singlePlyPerft)
	assert.Equal(t, 20, perft)

	// depth 2 starting position
	depth = 2
	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	perft = 0
	mv = new(moves.Move)
	minMax(depth, depth, &pos, &mv, &perft, &singlePlyPerft)
	assert.Equal(t, 400, perft)

	// depth 3 starting position
	depth = 3
	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	perft = 0
	mv = new(moves.Move)
	minMax(depth, depth, &pos, &mv, &perft, &singlePlyPerft)
	assert.Equal(t, 8902, perft)

	// good test position depth 1
	depth = 1
	pos, _ = position.NewPositionFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	perft = 0
	mv = new(moves.Move)
	minMax(depth, depth, &pos, &mv, &perft, &singlePlyPerft)
	assert.Equal(t, 48, perft)

	// good test position depth 2
	// depth = 2
	// pos, _ = position.NewPositionFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	// perft = 0
	// mv = new(moves.Move)
	// minMax(depth, depth, &pos, &mv, &perft, &singlePlyPerft)
	// fmt.Println("Perft: ", perft)
	// fmt.Print(mv)
	// assert.Equal(t, 2039, perft)

	// pawn promo testing
	// depth = 1
	// pos, _ = position.NewPositionFen("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1")
	// perft = 0
	// mv = new(moves.Move)
	// minMax(depth, depth, &pos, &mv, &perft, &singlePlyPerft)
	// assert.Equal(t, 24, perft)

	// depth = 2
	// pos, _ = position.NewPositionFen("n1n5/PPPk4/8/8/8/8/4Kppp/5N1N b - - 0 1")
	// perft = 0
	// mv = new(moves.Move)
	// minMax(depth, depth, &pos, &mv, &perft, &singlePlyPerft)
	// assert.Equal(t, 496, perft)
}
