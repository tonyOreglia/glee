package moves

import (
	"fmt"
	"strconv"
)

type Move struct {
	origin      int
	destination int
	// Queen = 2 Bishops = 3 Knights = 4 Rooks = 5
	promotion int
}

func NewMove(singleMove []int) *Move {
	mv := &Move{origin: singleMove[0], destination: singleMove[1], promotion: 0}
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
	fmt.Println(ConvertIndexToAlgebraic(m.origin) + ConvertIndexToAlgebraic(m.destination))
}

func (m *Move) String() string {
	mvString := ConvertIndexToAlgebraic(m.origin) + ConvertIndexToAlgebraic(m.destination)
	if (m.promotion > 0) && (m.destination < 8) {
		mvString += "Q"
	}
	if (m.promotion > 0) && (m.destination > 55) {
		mvString += "q"
	}
	return mvString
}

func ConvertIndexToAlgebraic(index int) string {
	var algebraic string
	column := index % 8
	row := 8 - (index / 8)
	switch column {
	case 0:
		algebraic = "a"
	case 1:
		algebraic = "b"
	case 2:
		algebraic = "c"
	case 3:
		algebraic = "d"
	case 4:
		algebraic = "e"
	case 5:
		algebraic = "f"
	case 6:
		algebraic = "g"
	case 7:
		algebraic = "h"
	}
	algebraic += strconv.Itoa(row)
	return algebraic
}

func ConvertAlgebriacToIndex(algebraic string) (int, error) {
	if algebraic == "-" {
		return 64, nil
	}
	column := string(algebraic[0])
	row, _ := strconv.Atoi(string(algebraic[1]))
	row--
	var index, columnValue, rowValue int
	switch column {
	case "a":
		columnValue = 0
	case "b":
		columnValue = 1
	case "c":
		columnValue = 2
	case "d":
		columnValue = 3
	case "e":
		columnValue = 4
	case "f":
		columnValue = 5
	case "g":
		columnValue = 6
	case "h":
		columnValue = 7
	default:
		return 64, fmt.Errorf("invalid algebraic notation")
	}
	rowValue = ((7 - row) * 8)
	index = columnValue + rowValue
	return index, nil
}
