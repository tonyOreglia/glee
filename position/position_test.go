package position

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPositionContructorBuiltin(t *testing.T) {
	position := new(Position)
	assert.Equal(t, Position{}, *position)
}

func TestPositionContructorFen(t *testing.T) {
	position, _ := NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	assert.Equal(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR", position.GetFenString()) //  w KQkq - 0 1
}

// func TestPositionUpdate(t *testing.T) {
// 	postion, _ := NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
// 	position.MakeMove("e2", "e4")
// 	assert.Equal(t, position.OutputFen, "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
// }

func TestTokenizeFen(t *testing.T) {
	position, activeSide, castlingRights, enPassante, moveCt, halfMoveCt := convertFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	assert.Equal(t, position, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	assert.Equal(t, activeSide, White)
	assert.Equal(t, castlingRights, "KQkq")
	assert.Equal(t, enPassante, 64)
	assert.Equal(t, moveCt, 0)
	assert.Equal(t, halfMoveCt, 1)
}

func TestConvertAlgebriacToIndex(t *testing.T) {
	index := convertAlgebriacToIndex("a8")
	assert.Equal(t, 0, index)
	index = convertAlgebriacToIndex("h1")
	assert.Equal(t, 63, index)
	index = convertAlgebriacToIndex("h8")
	assert.Equal(t, 7, index)
}
