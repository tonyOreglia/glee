package moves

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMovesConstructor(t *testing.T) {
	mvs := NewMovesList()
	assert.Equal(t, 100, cap(mvs.mvs))
	assert.Equal(t, 0, len(mvs.mvs))
}

func TestAttackedSqsBb(t *testing.T) {
	mvs := NewMovesList()
	mvs.AddMove(1, 0)
	mvs.AddMove(1, 0)
	mvs.AddMove(0, 1)
	mvs.AddMove(0, 2)
	mvs.AddMove(0, 3)
	attSqsBb := mvs.AttackedSqsBb()
	expected := uint64(15)
	assert.Equal(t, expected, attSqsBb.Value())
}
