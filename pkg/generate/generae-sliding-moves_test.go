package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyOreglia/glee/pkg/bitboard"
	"github.com/tonyOreglia/glee/pkg/hashtables"
)

func TestGenerateValidDirectionalMovesBb(t *testing.T) {
	// piece sliding southwest from h8 unblocked
	sw := hashtables.Lookup.SouthWestArrayBbHash
	index := 7
	validMvsBb := generateValidDirectionalMovesBb(index, sw, uint64(0), getMsb)
	expectedValidMvsBb := uint64(0x102040810204000)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())

	// piece sliding north from a3, blocked on a7
	north := hashtables.Lookup.NorthArrayBbHash
	index = 40
	validMvsBb = generateValidDirectionalMovesBb(index, north, uint64(0x100), getLsb)
	expectedValidMvsBb = uint64(0x101010100)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())
}

func TestGenerateValidDiagonalSlidingMovesBb(t *testing.T) {
	// piece sliding diagonally from e4 unblocked
	index := 36
	occSqsVal := uint64(0)
	validMvsBb := generateValidDiagonalSlidingMovesBb(index, occSqsVal, hashtables.Lookup)
	expectedValidMvsBb, _ := bitboard.NewBitboard(uint64(0x8244280028448201))
	assert.Equal(t, expectedValidMvsBb.Value(), validMvsBb.Value())

	// piece sliding diagonally from a3 blocked at c5
	index = 40
	occSqsBb, _ := bitboard.NewBitboard(0)
	occSqsBb.SetBit(26)
	validMvsBb = generateValidDiagonalSlidingMovesBb(index, occSqsBb.Value(), hashtables.Lookup)
	expectedValidMvsBb, _ = bitboard.NewBitboard(uint64(0x402000204000000))
	assert.Equal(t, expectedValidMvsBb.Value(), validMvsBb.Value())
}
