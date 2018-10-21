package utility

import "strconv"

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
