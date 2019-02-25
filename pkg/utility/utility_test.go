package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertAlgebriacToIndex(t *testing.T) {
	index, _ := ConvertAlgebriacToIndex("a8")
	assert.Equal(t, 0, index)
	index, _ = ConvertAlgebriacToIndex("h1")
	assert.Equal(t, 63, index)
	index, _ = ConvertAlgebriacToIndex("h8")
	assert.Equal(t, 7, index)
}
