package position

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertAlgebriacToIndex(t *testing.T) {
	index := convertAlgebriacToIndex("a8")
	assert.Equal(t, 0, index)
	index = convertAlgebriacToIndex("h1")
	assert.Equal(t, 63, index)
	index = convertAlgebriacToIndex("h8")
	assert.Equal(t, 7, index)
}

func TestTokenizeFen(t *testing.T) {
	position, activeSide, castlingRights, enPassante, moveCt, halfMoveCt := getFenStringTokens("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	assert.Equal(t, position, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR")
	assert.Equal(t, activeSide, White)
	assert.Equal(t, castlingRights, "KQkq")
	assert.Equal(t, enPassante, 64)
	assert.Equal(t, moveCt, 0)
	assert.Equal(t, halfMoveCt, 1)

	position, activeSide, castlingRights, enPassante, moveCt, halfMoveCt = getFenStringTokens("rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b q e3 1 2")
	assert.Equal(t, position, "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R")
	assert.Equal(t, activeSide, Black)
	assert.Equal(t, castlingRights, "q")
	assert.Equal(t, enPassante, 44)
	assert.Equal(t, moveCt, 1)
	assert.Equal(t, halfMoveCt, 2)
}

func TestPositionContructorFen(t *testing.T) {
	position, _ := NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Kkq - 0 1")
	assert.Equal(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Kkq - 0 1", position.GetFenString())

	position, _ = NewPositionFen("rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b q e3 1 2")
	assert.Equal(t, "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b q e3 1 2", position.GetFenString())

	position, _ = NewPositionFen("rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b - - 1 2")
	assert.Equal(t, "rnbqkbnr/pp1ppppp/8/2p5/4P3/5N2/PPPP1PPP/RNBQKB1R b - - 1 2", position.GetFenString())

	position, _ = NewPositionFen("7k/8/8/8/8/8/8/6KB w q - 0 1")
	assert.Equal(t, "7k/8/8/8/8/8/8/6KB w q - 0 1", position.GetFenString())

	position, _ = NewPositionFen("7k/8/8/8/8/8/8/Rq6 w KQkq - 0 1")
	assert.Equal(t, "7k/8/8/8/8/8/8/Rq6 w KQkq - 0 1", position.GetFenString())
}

func TestPositionUpdate(t *testing.T) {
	position, _ := NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	position.MakeMove("e2", "e3", White)
	assert.Equal(t, position.GetFenString(), "rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b KQkq - 1 1")
}
