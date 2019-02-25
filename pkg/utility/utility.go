package utility

import (
	"fmt"
	"strconv"
)

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
