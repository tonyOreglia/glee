package generatemoves

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyoreglia/glee/chessmoves"
	"github.com/tonyoreglia/glee/position"
)

func TestMinMax(t *testing.T) {
	// depth 1 starting position
	depth := 1
	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	perft := 0
	singlePlyPerft := 0
	mv := new(chessmoves.Move)
	minMax(depth, depth, &pos, &mv, &perft, &singlePlyPerft)
	fmt.Println("Perft: ", perft)
	fmt.Print(mv)
	assert.Equal(t, 20, perft)

	// depth 2 starting position
	depth = 2
	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	perft = 0
	mv = new(chessmoves.Move)
	minMax(depth, depth, &pos, &mv, &perft, &singlePlyPerft)
	fmt.Println("Perft: ", perft)
	fmt.Print(mv)
	assert.Equal(t, 400, perft)

	// depth 3 starting position
	depth = 3
	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	perft = 0
	mv = new(chessmoves.Move)
	minMax(depth, depth, &pos, &mv, &perft, &singlePlyPerft)
	fmt.Println("Perft: ", perft)
	fmt.Print(mv)
	assert.Equal(t, 8902, perft)
}
