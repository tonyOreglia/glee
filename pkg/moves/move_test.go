package moves

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveConstructor(t *testing.T) {
	var mvSlice = []int{1, 2}
	mv := NewMove(mvSlice)
	assert.Equal(t, mv.Origin(), 1)
	assert.Equal(t, mv.Destination(), 2)
	assert.Equal(t, mv.PromotionPiece(), 0)
	assert.Equal(t, mv.GetMoveSlice(), []int{1, 2})

	mvSlice = []int{62, 63, 3}
	mv = NewPromoMove(mvSlice)
	assert.Equal(t, mv.Origin(), 62)
	assert.Equal(t, mv.Destination(), 63)
	assert.Equal(t, mv.PromotionPiece(), 3)
	assert.Equal(t, mv.GetMoveSlice(), []int{62, 63})
}

func TestConvertAlgebriacToIndex(t *testing.T) {
	index, _ := ConvertAlgebriacToIndex("a8")
	assert.Equal(t, 0, index)
	index, _ = ConvertAlgebriacToIndex("h1")
	assert.Equal(t, 63, index)
	index, _ = ConvertAlgebriacToIndex("h8")
	assert.Equal(t, 7, index)
}
