package legalmoves

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/tonyoreglia/glee/hashtables"
// 	"github.com/tonyoreglia/glee/position"
// )

// func TestGenerateLegalMoves(t *testing.T) {
// 	ht := hashtables.CalculateAllLookupBbs()
// 	pos, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
// 	moves := NewLegalMoves(pos, ht)
// 	var expectedLegalMoves [][2]int
// 	expectedLegalMoves = make([][2]int, 0, 100)
// 	expectedLegalMoves = append(expectedLegalMoves, [2]int{1, 2})
// 	assert.Equal(t, expectedLegalMoves, moves)
// }
