// Package bitboard implements utility routines for
// using chess engine bitboards
package bitboard

import (
	"fmt"
	"math/bits"
)

// Bitboard struct exposes uint64 "bitboard" with associated getter, setter, and helper fxns
type Bitboard struct {
	bitboard uint64
}

func NewBitboard(bb uint64) (*Bitboard, error) {
	bitboard := &Bitboard{bitboard: bb}
	return bitboard, nil
}

// NewBbFromMovesSlice takes a legal moves slice and returns a bitboard representing those moves
func NewBbFromMovesSlice(mvs [][2]int) *Bitboard {
	bitboard := &Bitboard{}
	for i := 0; i < len(mvs); i++ {
		bitboard.SetBit(mvs[i][1])
	}
	return bitboard
}

func (b *Bitboard) Set(bb uint64) error {
	b.bitboard = bb
	return nil
}

func (b *Bitboard) Combine(bb *Bitboard) {
	b.bitboard |= bb.Value()
}

func (b *Bitboard) RemoveOverlappingBits(bb *Bitboard) {
	b.bitboard &^= bb.Value()
}

func (b *Bitboard) SetBit(bit int) error {
	b.bitboard |= 1 << uint(bit)
	return nil
}

func (b *Bitboard) IsBitSet(bit int) bool {
	if ((uint64(1) << uint(bit)) & b.bitboard) != uint64(0) {
		return true
	}
	return false
}

func (b *Bitboard) GetBitValue(bit int) uint {
	if b.IsBitSet(bit) {
		return 1
	}
	return 0
}

func (b *Bitboard) Value() uint64 {
	return b.bitboard
}

// count leading zeros
func (b *Bitboard) Msb() int {
	return 63 - int(bits.LeadingZeros64(b.Value()))
}

// count leading zeros
func (b *Bitboard) Lsb() int {
	return int(bits.TrailingZeros64(b.Value()))
}

// Does not check that this bit is set as speed is priority.
func (b *Bitboard) RemoveBit(pos int) error {
	b.bitboard ^= 1 << uint(pos)
	return nil
}

func (b *Bitboard) Print() {
	var i int
	fmt.Println("")
	for i = 0; i < 64; i++ {
		var sq int
		if b.IsBitSet(i) {
			sq = 1
		}
		fmt.Print(sq)
		if ((i + 1) % 8) == 0 {
			fmt.Println("")
		}
	}
	fmt.Println("")
	fmt.Println("")
}
