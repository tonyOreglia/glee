package moves

import (
	"fmt"

	"github.com/tonyOreglia/glee/pkg/bitboard"
)

const MoveListSize = 100

type Moves struct {
	mvs []Move
}

// NewMovesList creates instance of Moves struct
// Should be renamed to New()
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
	mv := NewPromoMove([]int{origin, dest, promo})
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
	for _, move := range m.mvs {
		move.Print()
		fmt.Println()
	}
}

func (m *Moves) AttackedSqsBb() *bitboard.Bitboard {
	bb := new(bitboard.Bitboard)
	for _, move := range m.mvs {
		bb.SetBit(move.destination)
	}
	return bb
}

func (m *Moves) GetMoves() [][]int {
	var mvList [][]int
	for _, move := range m.mvs {
		mvList = append(mvList, move.GetMoveSlice())
	}
	return mvList
}

func (m *Moves) FindMove(origin int, dest int) bool {
	for _, move := range m.mvs {
		if move.origin == origin && move.destination == dest {
			return true
		}
	}
	return false
}

// NewBbFromMovesSlice takes a legal moves slice and returns a bitboard representing those moves
func (m *Moves) GetBitboard() *bitboard.Bitboard {
	bitboard := &bitboard.Bitboard{}
	for _, v := range m.mvs {
		bitboard.SetBit(v.Destination())
	}
	return bitboard
}
