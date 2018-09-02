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

func (b *Bitboard) Get() (uint64) {
	return b.bitboard
}

// count leading zeros
func (b *Bitboard) Clz() (int) {
	return bits.LeadingZeros64(b.Get())
}

// count leading zeros
func (b *Bitboard) Ctz() (int) {
	return bits.TrailingZeros64(b.Get())
}