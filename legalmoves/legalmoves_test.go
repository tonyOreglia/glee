package legalmoves

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tonyoreglia/glee/hashtables"
)

// func TestGenerateLegalMoves(t *testing.T) {
// 	ht := hashtables.CalculateAllLookupBbs()
// 	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
// 	moves := NewLegalMoves(pos, ht)
// 	var expectedLegalMoves [][2]int
// 	expectedLegalMoves = make([][2]int, 0, 100)
// 	expectedLegalMoves = append(expectedLegalMoves, [2]int{1, 2})
// 	assert.Equal(t, expectedLegalMoves, moves)
// }

var ht = hashtables.CalculateAllLookupBbs()

func TestGenerateValidDirectionalMovesBb(t *testing.T) {
	// piece sliding southwest from h8 unblocked
	ht := hashtables.CalculateAllLookupBbs()
	sw := ht.SouthWestArrayBbHash
	index := 7
	validMvsBb := generateValidDirectionalMovesBb(index, sw, uint64(0))
	validMvsBb.Print()
	expectedValidMvsBb := uint64(0x102040810204000)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())

	// piece sliding north from a3, blocked on a7
	north := ht.NorthArrayBbHash
	index = 40
	validMvsBb = generateValidDirectionalMovesBb(index, north, uint64(0x100))
	validMvsBb.Print()
	expectedValidMvsBb = uint64(0x101010100)
	assert.Equal(t, expectedValidMvsBb, validMvsBb.Value())
}
