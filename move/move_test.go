package move

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveConstructor(t *testing.T) {
	var mvSlice = []int{1, 2}
	mv := NewMove(mvSlice)
	assert.Equal(t, mv.GetOrigin(), 1)
	assert.Equal(t, mv.GetDestination(), 2)
	assert.Equal(t, mv.GetPromoPiece(), 0)
	assert.Equal(t, mv.GetMoveSlice(), []int{1, 2})

	mvSlice = []int{62, 63, 3}
	mv = NewPromoMove(mvSlice)
	assert.Equal(t, mv.GetOrigin(), 62)
	assert.Equal(t, mv.GetDestination(), 63)
	assert.Equal(t, mv.GetPromoPiece(), 3)
	assert.Equal(t, mv.GetMoveSlice(), []int{62, 63})
}
