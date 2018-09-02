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
	value, _ := bitboard.Get()
	assert.Equal(t, value, uint64(100))
}

func TestBitboardLsb(t *testing.T) {
	bitboard, _ := NewBitboard(1)
	position, _ := bitboard.Lsb()
	assert.Equal(t, position, 0)
}

// func TestBitboardMsb(t *testing.T) {

// }
