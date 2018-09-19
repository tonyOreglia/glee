// Package moves stores a list of chess moves.
// Each moves contains the necessary information needed for making
// and un-making the move.
package moves

import (
	"fmt"

	"github.com/tonyoreglia/glee/bitboard"
	"github.com/tonyoreglia/glee/move"
)

const MoveListSize = 100

type Moves struct {
	mvs []move.Move
}

// NewMovesList creates instance of Moves struct
func NewMovesList() *Moves {
	movesStruct := new(Moves)
	movesStruct.mvs = make([]move.Move, 0, 100)
	return movesStruct
}

func (m *Moves) GetMovesList() []move.Move {
	return m.mvs
}

// AddMove adds a move to instance of Moves struct
func (m *Moves) AddMove(origin int, dest int) {
	mv := move.NewMove([]int{origin, dest})
	m.mvs = append(m.mvs, *mv)
}

func (m *Moves) AddMoveFromSlice(singleMove []int) {
	mv := move.NewMove(singleMove)
	m.mvs = append(m.mvs, *mv)
}

func (m *Moves) AddPromotionMove(origin int, dest int, promo int) {
	mv := move.NewMove([]int{origin, dest, promo})
	m.mvs = append(m.mvs, *mv)
}

func (m *Moves) PopMove() *move.Move {
	mv := &m.mvs[len(m.mvs)-1]
	m.mvs = m.mvs[0 : len(m.mvs)-1]
	return mv
}

func (m *Moves) Length() int {
	return len(m.mvs)
}

func (m *Moves) Print() {
	fmt.Print(m.mvs)
}

func (m *Moves) GetMoves() [][]int {
	var mvList [][]int
	for _, move := range m.mvs {
		mvList = append(mvList, move.GetMoveSlice())
	}
	return mvList
}

// NewBbFromMovesSlice takes a legal moves slice and returns a bitboard representing those moves
func (m *Moves) GetBitboard() *bitboard.Bitboard {
	bitboard := &bitboard.Bitboard{}
	for _, v := range m.mvs {
		bitboard.SetBit(v.GetDestination())
	}
	return bitboard
}
