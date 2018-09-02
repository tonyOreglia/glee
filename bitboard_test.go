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
	value := bitboard.Get()
	assert.Equal(t, value, uint64(100))
}

func TestBitboardClz(t *testing.T) {
	bitboard, _ := NewBitboard(1)
	position := bitboard.Clz()
	assert.Equal(t, uint8(63), position)
}

func TestBitboardCtz(t *testing.T) {
	bitboard, _ := NewBitboard(1)
	position := bitboard.Ctz()
	assert.Equal(t, uint8(0), position)
}

func TestBitboardSetBit(t *testing.T) {
	bitboard, _ := NewBitboard(0)
	bitboard.SetBit(7)
	assert.Equal(t, bitboard.Get(), uint64(128))
}

func TestBitboardPop(t *testing.T) {
	bitboard, _ := NewBitboard(0)
	bitboard.SetBit(7)
	bitboard.Pop(bitboard.Ctz())
	assert.Equal(t, bitboard.Get(), uint64(0))
}

