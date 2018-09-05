package bitboard

import (
	"math/bits"
	"fmt"
	"strconv"
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

func (b *Bitboard) GetBit(bit uint8) int {
	if ((uint64(1) << bit) & b.bitboard) != uint64(0) { return 1 }
	return 0
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

func (b *Bitboard) Print() {
	var i uint8
	fmt.Println("")
	for i = 0; i<64; i++ {
			fmt.Print(strconv.Itoa(b.GetBit(i)))
			if ((i+1) % 8) == 0 { fmt.Println("") }
	}
	fmt.Println("")
	fmt.Println("")
}