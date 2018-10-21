package moves

import (
	"fmt"
)

type Move struct {
	origin      int
	destination int
	// Queen = 2 Bishops = 3 Knights = 4 Rooks = 5
	promotion int
}

func NewMove(singleMove []int) *Move {
	mv := &Move{origin: singleMove[0], destination: singleMove[1]}
	return mv
}

func NewPromoMove(singleMove []int) *Move {
	mv := &Move{origin: singleMove[0], destination: singleMove[1], promotion: singleMove[2]}
	return mv
}

func (m *Move) CopyMove() *Move {
	x := []int{m.origin, m.destination, m.promotion}
	return NewPromoMove(x)
}

func (m *Move) Origin() int {
	return m.origin
}

func (m *Move) Destination() int {
	return m.destination
}

func (m *Move) GetMoveSlice() []int {
	return []int{m.origin, m.destination}
}

func (m *Move) PromotionPiece() int {
	return m.promotion
}

func (m *Move) Print() {
	fmt.Print(m)
}
