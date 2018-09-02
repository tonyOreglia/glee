package bitboard

import ()

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

func (b *Bitboard) Get() (uint64, error) {
	return b.bitboard, nil
}