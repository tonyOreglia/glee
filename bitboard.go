package bitboard

import (
	"math/bits"
)

type Bitboard struct {
 	bitboard uint64
}

func NewBitboard(bb uint64) (*Bitboard, error) {
	bitboard := &Bitboard{ bitboard:bb }
	return bitboard, nil
}

func (b *Bitboard) Set(bb uint64) error {
	b.bitboard = bb
	return nil
}

func (b *Bitboard) SetBit(bit uint8) error {
	b.bitboard = 1 << bit
	return nil
}

func (b *Bitboard) Get() (uint64) {
	return b.bitboard
}

// count leading zeros
func (b *Bitboard) Clz() (uint8) {
	return uint8(bits.LeadingZeros64(b.Get()))
}

// count leading zeros
func (b *Bitboard) Ctz() (uint8) {
	return uint8(bits.TrailingZeros64(b.Get()))
}

// Does not check that this bit is set as speed is priority.
func (b *Bitboard) Pop(pos uint8) (error) {
	b.bitboard ^= 1 << pos
	return nil
}