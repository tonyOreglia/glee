package bitboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitboardContructor(t *testing.T) {
	bitboard, _ := NewBitboard(100)
	assert.Equal(t, bitboard.bitboard, uint64(100))
}

func TestBitboardSet(t *testing.T) {
	bitboard, _ := NewBitboard(100)
	bitboard.Set(200)
	assert.Equal(t, bitboard.bitboard, uint64(200))
}

func TestBitboardGet(t *testing.T) {
	bitboard, _ := NewBitboard(100)
	value := bitboard.Value()
	assert.Equal(t, value, uint64(100))
}

func TestBitboardRemoveBit(t *testing.T) {
	// should remove a set bit
	bitboard, _ := NewBitboard(1)
	assert.Equal(t, uint64(1), bitboard.Value())
	bitboard.RemoveBit(0)
	assert.Equal(t, uint64(0), bitboard.Value())

	// should not switch bit if it's not set
	bitboard.RemoveBit(0)
	assert.Equal(t, uint64(0), bitboard.Value())
}

func TestBitboardLsb(t *testing.T) {
	bitboard, _ := NewBitboard(1)
	position := bitboard.Lsb()
	assert.Equal(t, 0, position)
}

func TestBitboardMsb(t *testing.T) {
	bitboard, _ := NewBitboard(1)
	position := bitboard.Msb()
	assert.Equal(t, 0, position)
}

func TestBitboardSetBit(t *testing.T) {
	bitboard, _ := NewBitboard(0)
	bitboard.SetBit(7)
	assert.Equal(t, bitboard.Value(), uint64(128))
}

func TestBitboardPop(t *testing.T) {
	bitboard, _ := NewBitboard(0)
	bitboard.SetBit(7)
	bitboard.RemoveBit(bitboard.Msb())
	assert.Equal(t, bitboard.Value(), uint64(0))
}

func TestBitboardGetBit(t *testing.T) {
	bitboard, _ := NewBitboard(1)
	assert.Equal(t, bitboard.BitIsSet(0), true)
	assert.Equal(t, bitboard.BitIsSet(63), false)
	bitboard, _ = NewBitboard(2)
	assert.Equal(t, bitboard.BitIsSet(1), true)
	assert.Equal(t, bitboard.BitIsSet(63), false)
	bitboard, _ = NewBitboard(4)
	assert.Equal(t, bitboard.BitIsSet(2), true)
	assert.Equal(t, bitboard.BitIsSet(63), false)
	bitboard, _ = NewBitboard(8)
	assert.Equal(t, bitboard.BitIsSet(3), true)
	assert.Equal(t, bitboard.BitIsSet(63), false)
	bitboard, _ = NewBitboard(16)
	assert.Equal(t, bitboard.BitIsSet(4), true)
	assert.Equal(t, bitboard.BitIsSet(63), false)
	bitboard, _ = NewBitboard(9223372036854775808)
	assert.Equal(t, bitboard.BitIsSet(4), false)
	assert.Equal(t, bitboard.BitIsSet(62), false)
	assert.Equal(t, bitboard.BitIsSet(63), true)
}
