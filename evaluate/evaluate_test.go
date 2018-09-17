package evaluate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyoreglia/glee/position"
)

func TestEvaluatePosition(t *testing.T) {
	// starting position should be evaluated to zero
	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	score := EvaluatePosition(pos)
	assert.Equal(t, 0, score)

	pos, _ = position.NewPositionFen("k7/8/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	score = EvaluatePosition(pos)
	assert.Equal(t, 4010, score)

	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/8/7K w KQkq - 0 1")
	score = EvaluatePosition(pos)
	assert.Equal(t, -4010, score)\
}
