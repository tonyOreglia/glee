package main

import (
	"github.com/tonyoreglia/glee/hashtables"
	"github.com/tonyoreglia/glee/position"
)

func main() {
	hashtables.CalculateAllLookupBbs()
}

func verifyStartingPositionVisually() {
	position, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	position.PrintBitboards()
}
