package main

import (
	"fmt"
	"os"

	"github.com/tonyOreglia/glee/pkg/evaluate"
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/position"
)

var ht = hashtables.Lookup

func main() {
	command := make([]byte, 0, 100)
	pos := position.StartingPosition()
	for true {
		fmt.Print("glee: ")
		_, err := fmt.Scanln(&command)
		if err != nil {
			fmt.Print(err)
		}
		c := string(command)
		switch {
		case c == "help":
			printHelp()
		case c == "quit":
			os.Exit(0)
		case c == "fen":
			pos.PrintFen()
		case c == "new":
			pos = position.StartingPosition()
		case c == "disp":
			pos.Print()
		case c == "eval":
			fmt.Printf("score: %d\n", evaluate.EvaluatePosition(pos))
		case c == "uci":
			uci()
		default:
			handleMove(c)
		}
	}
}

func printHelp() {
	fmt.Println("Glee 0.0.0 - GoLang chEss Engine")
	fmt.Println("quit............terminates the program")
	fmt.Println("uci.............switch to uci-mode")
	fmt.Println("e2e4............moves piece")
	fmt.Println("e7e8Q...........promotion move resulting in Queen [Q,R,B,N]")
	// fmt.Println("st #............sets search time per move (1-300s)")
	// fmt.Println("sd #............sets search depth (1-9)")
	// fmt.Println("undo............takes back last move")
	fmt.Println("new.............resets board to initial state")
	// fmt.Println("clear...........redraws the console")
	fmt.Println("disp............shows the board")
	// fmt.Println("divide #........outputs the numbers of child moves")
	// fmt.Println("divide2 #.......outputs the total numbers of child moves")
	// fmt.Println("perft #.........counts nodes at given depth")
	// fmt.Println("perft2 #........counts all nodes to given depth")
	// fmt.Println("setboard <FEN>..reads a fen-string")
	fmt.Println("fen.............outputs FEN of board position")
	// fmt.Println("info............outputs data-structure")
	fmt.Println("eval............evaluates position")
	// fmt.Println("analyze.........infinite analysis")
	// fmt.Println("stack...........shows move-stack")
	// fmt.Println("sort............gives sorted move-list for Alpha-Beta")
	// fmt.Println("show............gives valid moves for current pos.")
}

func handleMove(mv string) {
	if len(mv) != 4 && len(mv) != 5 {
		badInput(mv)
	}
}

func uci() {
	command := make([]byte, 0, 100)
	fmt.Println("GLEE-GoLang chEss Engine")
	fmt.Println("tony.oreglia@gmail.com")
	fmt.Println("id name Glee 0.0.1")
	fmt.Println("id author Tony Oreglia")
	// fmt.Println("option name useInternalOpeningBook type check default false")
	// fmt.Println("option name ponder type check default true")
	// fmt.Println("option name Quiescent Search type check default true")
	// fmt.Println("option name Logging type check default true")
	// fmt.Println("option name NullMove Reduction type spin default 2 min 0 max 4")
	// fmt.Println("option name Research Value type spin default 50 min 0 max 150")
	// fmt.Println("option name UCI_Chess960 type check default false")
	fmt.Println("uciok")
	for true {
		_, err := fmt.Scanln(&command)
		if err != nil {
			fmt.Print(err)
			os.Exit(0)
		}
	}
}

func badInput(c string) {
	fmt.Printf("\ninput correct ?: %s\n\n", c)
}
