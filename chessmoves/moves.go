package chessmoves

import (
	"fmt"

	"github.com/tonyoreglia/glee/bitboard"
)

const MoveListSize = 100

type Moves struct {
	mvs []Move
}

// NewMovesList creates instance of Moves struct
func NewMovesList() *Moves {
	movesStruct := new(Moves)
	movesStruct.mvs = make([]Move, 0, 100)
	return movesStruct
}

func (m *Moves) GetMovesList() []Move {
	return m.mvs
}

// AddMove adds a move to instance of Moves struct
func (m *Moves) AddMove(origin int, dest int) {
	mv := NewMove([]int{origin, dest})
	m.mvs = append(m.mvs, *mv)
}

func (m *Moves) AddMoveFromSlice(singleMove []int) {
	mv := NewMove(singleMove)
	m.mvs = append(m.mvs, *mv)
}

func (m *Moves) AddPromotionMove(origin int, dest int, promo int) {
	mv := NewMove([]int{origin, dest, promo})
	m.mvs = append(m.mvs, *mv)
}

func (m *Moves) PopMove() *Move {
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
