package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/tonyOreglia/glee/pkg/engine"
	"github.com/tonyOreglia/glee/pkg/evaluate"
	"github.com/tonyOreglia/glee/pkg/generate"
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/moves"
	"github.com/tonyOreglia/glee/pkg/position"
	"github.com/tonyOreglia/glee/pkg/utility"
)

const version = "0.0.1"

var ht = hashtables.Lookup

func main() {
	command := make([]byte, 0, 100)
	pos := position.StartingPosition()
	mvs := generate.GenerateMoves(pos)
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
		case "uci":
			uci()
		case "undo":
			undo(pos)
			mvs = generate.GenerateMoves(pos)
		case "search":
			search(pos, mvs)
			mvs = generate.GenerateMoves(pos)
		case "setboard":
			setboard(pos)
		case "playw":
			pos = position.StartingPosition()
			play(pos, 0)
		case "playb":
			pos = position.StartingPosition()
			play(pos, 1)
		default:
			handleMove(c, pos, mvs)
			mvs = generate.GenerateMoves(pos)
		}
	}
}

func play(p *position.Position, humanSide int) {
	move := make([]byte, 0, 100)
	for true {
		if p.GetActiveSide() == humanSide {
			for true {
				fmt.Print("move: ")
				_, err := fmt.Scan(&move)
				if err != nil {
					fmt.Print(err)
				}
				if string(move) == "quit" {
					os.Exit(0)
				}
				if handleMove(string(move), p, generate.GenerateMoves(p)) {
					break
				}
			}
		} else {
			search(p, generate.GenerateMoves(p))
		}
	}
}

func search(p *position.Position, mvs *moves.Moves) {
	perft := 0
	singlePlyPerft := 0
	params := engine.SearchParams{
		Depth:          5,
		Ply:            5,
		Pos:            &p,
		Perft:          &perft,
		SinglePlyPerft: &singlePlyPerft,
		EngineMove:     &moves.Move{},
	}
	if p.IsWhitesTurn() {
		engine.AlphaBetaMax(-10000, 10000, 5, params)
	} else {
		engine.AlphaBetaMin(-10000, 10000, 5, params)
	}
	p.Move(*params.EngineMove)
	p.Print()
	fmt.Print("glee move: ")
	params.EngineMove.Print()
}

func printHelp() {
	fmt.Printf("Glee %s - GoLang chEss Engine\n", version)
	fmt.Println("quit............terminates the program")
	fmt.Println("uci.............switch to uci-mode")
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

func setboard(p *position.Position) {
	in := bufio.NewReader(os.Stdin)
	fen, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	fen = fen[0 : len(fen)-1]
	p, err = position.NewPositionFen(string(fen))
	if err != nil {
		badInput(string(fen))
	}
	p.Print()
}

func handleMove(mv string, p *position.Position, mvs *moves.Moves) bool {
	lookupPromo := map[string]int{
		// Queen = 2 Bishops = 3 Knights = 4 Rooks = 5
		"Q": 2,
		"B": 3,
		"N": 4,
		"R": 5,
	}
	promotionPiece := 0
	if len(mv) != 4 && len(mv) != 5 {
		badInput(mv)
		return false
	}
	if len(mv) == 5 {
		promotionPiece = lookupPromo[string(mv[4])]
	}
	origin, err := utility.ConvertAlgebriacToIndex(mv[0:2])
	if err != nil {
		badInput(mv)
		return false
	}
	dest, err := utility.ConvertAlgebriacToIndex(mv[2:4])
	if err != nil {
		badInput(mv)
		return false
	}
	move, found := mvs.FindMove(origin, dest, promotionPiece)
	if !found {
		badInput(mv)
		return false
	}
	if !engine.MakeValidMove(move, &p) {
		badInput(mv)
		return false
	}
	p.Print()
	return true
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
		switch string(command) {
		case "uci":
			fmt.Println("not yet implemented")
		case "debug":
			fmt.Println("not yet implemented")
		case "isready":
			fmt.Println("readyok")
		case "setoption":
			fmt.Println("not yet implemented")
		case "register":
			fmt.Println("not yet implemented")
		case "later":
			fmt.Println("not yet implemented")
		case "name":
			fmt.Println("not yet implemented")
		case "code":
			fmt.Println("not yet implemented")
		case "ucinewgame":
		case "position":
			fmt.Println("not yet implemented")
		case "go":
			fmt.Println("not yet implemented")
		case "searchmoves":
			fmt.Println("not yet implemented")
		case "ponder":
			fmt.Println("not yet implemented")
		case "wtime":
			fmt.Println("not yet implemented")
		case "btime":
			fmt.Println("not yet implemented")
		case "winc":
			fmt.Println("not yet implemented")
		case "binc":
			fmt.Println("not yet implemented")
		case "movestogo":
			fmt.Println("not yet implemented")
		case "depth":
			fmt.Println("not yet implemented")
		case "nodes":
			fmt.Println("not yet implemented")
		case "mate":
			fmt.Println("not yet implemented")
		case "movetime":
			fmt.Println("not yet implemented")
		case "infinite":
			fmt.Println("not yet implemented")
		case "stop":
			fmt.Println("not yet implemented")
		case "ponderhit":
			fmt.Println("not yet implemented")
		case "quit":
			fmt.Println("not yet implemented")
		}
	}
}

func badInput(c string) {
	fmt.Printf("\ninput correct ?: %s\n\n", c)
}

func undo(p *position.Position) {
	newPos := p.UnMakeMove()
	if newPos != nil {
		p = newPos
		p.Print()
		return
	}
	badInput("no previous move to undo")
}
