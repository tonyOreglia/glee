package main

import (
	"fmt"

	ht "github.com/tonyoreglia/glee/hashtables"
	"github.com/tonyoreglia/glee/position"
)

func main() {
	fmt.Println(ht.Lookup.EastArrayBbHash[7])
	fmt.Println("hello")
}

func verifyStartingPositionVisually() {
	position, _ := position.NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	position.PrintBitboards()
}
