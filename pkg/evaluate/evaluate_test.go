package evaluate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyOreglia/glee/pkg/position"
)

func TestEvaluatePosition(t *testing.T) {
	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	score := EvaluatePosition(pos)
	assert.Equal(t, 0, score)

	pos, _ = position.NewPositionFen("k7/8/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	score = EvaluatePosition(pos)
	assert.True(t, score > 3000)

	pos, _ = position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/8/7K w KQkq - 0 1")
	score = EvaluatePosition(pos)
	assert.True(t, score < -3000)
}
