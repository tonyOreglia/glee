package position

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyOreglia/glee/pkg/moves"
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
	position.MakeMoveAlgebraic("e2", "e3")
	assert.Equal(t, position.GetFenString(), "rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b KQkq - 1 1")
}

func TestWhiteCanCastleKingSide(t *testing.T) {
	position, _ := NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w K - 0 1")
	assert.True(t, position.WhiteCanCastleKingSide())
	position, _ = NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w k - 0 1")
	assert.False(t, position.WhiteCanCastleKingSide())
}

func TestWhiteCanCastleQueenSide(t *testing.T) {
	position, _ := NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w Q - 0 1")
	assert.True(t, position.WhiteCanCastleQueenSide())
	position, _ = NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1")
	assert.False(t, position.WhiteCanCastleQueenSide())
}

func TestBlackCanCastleKingSide(t *testing.T) {
	position, _ := NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b k - 0 1")
	assert.True(t, position.BlackCanCastleKingSide())
	position, _ = NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b - - 0 1")
	assert.False(t, position.BlackCanCastleKingSide())
}

func TestBlackCanCastleQueenSide(t *testing.T) {
	position, _ := NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b q - 0 1")
	assert.True(t, position.BlackCanCastleQueenSide())
	position, _ = NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR b Q - 0 1")
	assert.False(t, position.BlackCanCastleQueenSide())
}

func TestUnMakeMove(t *testing.T) {
	// unmake single move
	position, _ := NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	position.MakeMoveAlgebraic("e2", "e3")
	position = position.UnMakeMove()
	assert.Equal(t, position.GetFenString(), "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	// unmaking multiple moves in a row
	position, _ = NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	position.MakeMoveAlgebraic("e2", "e3")
	position.MakeMoveAlgebraic("e7", "e6")
	position.MakeMoveAlgebraic("d2", "d4")
	assert.Equal(t, "rnbqkbnr/pppp1ppp/4p3/8/3P4/4P3/PPP2PPP/RNBQKBNR b KQkq d3 2 1", position.GetFenString())
	position = position.UnMakeMove()
	position = position.UnMakeMove()
	position = position.UnMakeMove()
	assert.Equal(t, position.GetFenString(), "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	// unmake attacking move
	position, _ = NewPositionFen("7k/8/8/8/8/8/7p/6KR w q - 0 1")
	position.MakeMoveAlgebraic("h1", "h2")
	assert.Equal(t, position.GetFenString(), "7k/8/8/8/8/8/7R/6K1 b q - 1 1")
	position = position.UnMakeMove()
	assert.Equal(t, position.GetFenString(), "7k/8/8/8/8/8/7p/6KR w q - 0 1")

	//unmake en passante move
	position, _ = NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	position.MakeMoveAlgebraic("e2", "e4")
	assert.Equal(t, position.GetFenString(), "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 1 1")
	position = position.UnMakeMove()
	assert.Equal(t, position.GetFenString(), "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
}

func TestCastling(t *testing.T) {
	tests := map[string]struct {
		pos      string
		move     [2]string
		expected string
	}{
		"moving white king removes castling": {
			pos:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			move:     [2]string{"e1", "g1"},
			expected: "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R4RK1 b kq - 1 1",
		},
		"moving white king removes castling 2": {
			pos:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			move:     [2]string{"e1", "c1"},
			expected: "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/2KR3R b kq - 1 1",
		},
		"moving black king remove castling": {
			pos:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq - 0 1",
			move:     [2]string{"e8", "g8"},
			expected: "r4rk1/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQ - 0 1",
		},
		"moving black king removes castling rights 2": {
			pos:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq - 0 1",
			move:     [2]string{"e8", "c8"},
			expected: "2kr3r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQ - 0 1",
		},
		"moving white rook removes queenside castling rights": {
			pos:      "r3k2r/p1ppqNb1/bn2pnp1/3P4/4P3/2p2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			move:     [2]string{"a1", "b1"},
			expected: "r3k2r/p1ppqNb1/bn2pnp1/3P4/4P3/2p2Q1p/PPPBBPPP/1R2K2R b Kkq - 1 1",
		},
		"moving white rook removes kingside castling rights": {
			pos:      "r3k2r/p1ppqNb1/bn2pnp1/3P4/4P3/2p2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			move:     [2]string{"h1", "g1"},
			expected: "r3k2r/p1ppqNb1/bn2pnp1/3P4/4P3/2p2Q1p/PPPBBPPP/R3K1R1 b Qkq - 1 1",
		},
		"moving black rook removes queenside castling rights": {
			pos:      "r3k2r/p1ppqNb1/bn2pnp1/3P4/4P3/2p2Q1p/PPPBBPPP/R3K2R b KQkq - 0 1",
			move:     [2]string{"a8", "b8"},
			expected: "1r2k2r/p1ppqNb1/bn2pnp1/3P4/4P3/2p2Q1p/PPPBBPPP/R3K2R w KQk - 0 1",
		},
		"moving black rook removes kingside castling rights": {
			pos:      "r3k2r/p1ppqNb1/bn2pnp1/3P4/4P3/2p2Q1p/PPPBBPPP/R3K2R b KQkq - 0 1",
			move:     [2]string{"h8", "g8"},
			expected: "r3k1r1/p1ppqNb1/bn2pnp1/3P4/4P3/2p2Q1p/PPPBBPPP/R3K2R w KQq - 0 1",
		},
	}
	for tName, test := range tests {
		position, _ := NewPositionFen(test.pos)
		position.MakeMoveAlgebraic(test.move[0], test.move[1])
		assert.Equal(t, test.expected, position.GetFenString(), tName)
		position = position.UnMakeMove()
		assert.Equal(t, test.pos, position.GetFenString(), tName)
	}
}

func TestEnPassanteAttackMove(t *testing.T) {
	position, _ := NewPositionFen("r3k2r/p1ppqNb1/1n2pnp1/1b1P4/Pp2P3/2N2Q1p/1PPBBPPP/R3K2R b KQkq a3 0 1")
	position.MakeMoveAlgebraic("b4", "a3")
	assert.Equal(t, "r3k2r/p1ppqNb1/1n2pnp1/1b1P4/4P3/p1N2Q1p/1PPBBPPP/R3K2R w KQkq - 0 1", position.GetFenString())
}

func TestIsCastlingMove(t *testing.T) {
	position, _ := NewPositionFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	mv := moves.NewMove([]int{60, 62})
	assert.True(t, position.IsCastlingMove(*mv))

	position, _ = NewPositionFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	mv = moves.NewMove([]int{60, 58})
	assert.True(t, position.IsCastlingMove(*mv))

	position, _ = NewPositionFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq - 0 1")
	mv = moves.NewMove([]int{4, 2})
	assert.True(t, position.IsCastlingMove(*mv))

	position, _ = NewPositionFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R b KQkq - 0 1")
	mv = moves.NewMove([]int{4, 6})
	assert.True(t, position.IsCastlingMove(*mv))

	position, _ = NewPositionFen("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	mv = moves.NewMove([]int{60, 61})
	assert.False(t, position.IsCastlingMove(*mv))
}
