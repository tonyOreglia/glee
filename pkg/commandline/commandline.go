package commandline

import (
	"fmt"
	"os"

	"github.com/tonyOreglia/glee/pkg/evaluate"
	"github.com/tonyOreglia/glee/pkg/generate"
	"github.com/tonyOreglia/glee/pkg/moves"
	"github.com/tonyOreglia/glee/pkg/position"
)

const version = "0.0.1"

func CLI() {
	command := make([]byte, 0, 100)
	pos := position.StartingPosition()
	mvs := generate.GenerateMoves(pos)
	var move *moves.Move
	for true {
		fmt.Print("glee: ")
		_, err := fmt.Scan(&command)
		if err != nil {
			fmt.Print(err)
		}
		c := string(command)
		switch c {
		case "help":
			printHelp()
		case "quit":
			os.Exit(0)
		case "fen":
			pos.PrintFen()
		case "new":
			pos = position.StartingPosition()
		case "disp":
			pos.Print()
		case "eval":
			fmt.Printf("score: %d\n", evaluate.EvaluatePosition(pos))
		case "undo":
			undo(pos)
			mvs = generate.GenerateMoves(pos)
		case "search":
			pos, move = search(pos, mvs)
			move.Print()
			mvs = generate.GenerateMoves(pos)
		case "setboard":
			setboard(pos)
		case "playw":
			pos = position.StartingPosition()
			pos.Print()
			pos = play(pos, 0)
		case "playb":
			pos = position.StartingPosition()
			pos = play(pos, 1)
		default:
			handleMove(c, pos, mvs)
			pos.Print()
			mvs = generate.GenerateMoves(pos)
		}
	}
}

func printHelp() {
	fmt.Printf("Glee %s - GoLang chEss Engine\n", version)
	fmt.Println("quit............terminates the program")
	// fmt.Println("uci.............switch to uci-mode")
	fmt.Println("e2e4............moves piece")
	fmt.Println("e7e8Q...........promotion move resulting in Queen [Q,R,B,N]")
	// fmt.Println("st #............sets search time per move (1-300s)")
	// fmt.Println("sd #............sets search depth (1-9)")
	fmt.Println("undo............takes back last move")
	fmt.Println("new.............resets board to initial state")
	fmt.Println("disp............shows the board")
	fmt.Println("search..........engine plays the current position")
	fmt.Println("playw...........play white vs engine as black")
	fmt.Println("playb...........play black vs engine as white")
	// fmt.Println("divide #........outputs the numbers of child moves")
	// fmt.Println("divide2 #.......outputs the total numbers of child moves")
	// fmt.Println("perft #.........counts nodes at given depth")
	// fmt.Println("perft2 #........counts all nodes to given depth")
	fmt.Println("setboard <FEN>..reads a fen-string")
	fmt.Println("fen.............outputs FEN of board position")
	// fmt.Println("info............outputs data-structure")
	fmt.Println("eval............evaluates position")
	// fmt.Println("analyze.........infinite analysis")
	// fmt.Println("stack...........shows move-stack")
	// fmt.Println("sort............gives sorted move-list for Alpha-Beta")
	// fmt.Println("show............gives valid moves for current pos.")
}

func play(p *position.Position, humanSide int) *position.Position {
	move := make([]byte, 0, 100)
	for true {
		if p.GetActiveSide() == humanSide {
			for true {
				fmt.Print("human move: ")
				_, err := fmt.Scan(&move)
				if err != nil {
					fmt.Print(err)
				}
				if string(move) == "quit" {
					return p
				}
				if string(move) == "fen" {
					p.PrintFen()
					break
				}
				if string(move) == "help" {
					fmt.Printf("\nGlee %s - GoLang chEss Engine\n", version)
					fmt.Println("e2e4............moves piece")
					fmt.Println("e7e8Q...........promotion move resulting in Queen [Q,R,B,N]")
					fmt.Println("fen.............outputs FEN of board position")
					fmt.Println("quit............terminates the game")
					break;
				}
				if handleMove(string(move), p, generate.GenerateMoves(p)) {
					break
				}
			}
		} else {
			var mv *moves.Move
			p, mv = search(p, generate.GenerateMoves(p))
			p.Move(*mv)
			p.Print()
			fmt.Print("glee move: ")
			mv.Print()
		}
	}
	return p
}
