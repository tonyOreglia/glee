package engine

import (
	"fmt"
	"log"

	"github.com/tonyOreglia/glee/pkg/bitboard"
	"github.com/tonyOreglia/glee/pkg/evaluate"
	"github.com/tonyOreglia/glee/pkg/generate"
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/moves"
	"github.com/tonyOreglia/glee/pkg/position"
)

var ht = hashtables.Lookup

func minMax(depth int, ply int, pos **position.Position, engineMove **moves.Move, perft *int, singlePlyPerft *int) int {
	var value, tempValue int
	root := ply == depth
	if root {
		fmt.Println("")
		fmt.Print("Perft Divide: ")
		fmt.Println((*pos).GetFenString())
	}
	generate.GenerateMoves(*pos)
	// moveGenerator.GenerateMoves()
	if ply == 0 {
		return evaluate.EvaluatePosition(*pos, singlePlyPerft)
	}
	if (*pos).IsWhitesTurn() {
		value = -30000
		generateMoves(&value, &tempValue, root, depth, ply, pos, engineMove, perft, singlePlyPerft)
		return value
	}
	value = 30000
	generateMoves(&value, &tempValue, root, depth, ply, pos, engineMove, perft, singlePlyPerft)
	return value
}

func castlingMoveIsValid(move moves.Move, value *int, tempValue *int, root bool, depth int, ply int, pos **position.Position, engineMove **moves.Move, perft *int, singlePlyPerft *int) bool {
	tempPos := *pos
	*pos = (*pos).UnMakeMove()
	legalMoves := generate.GenerateMoves(*pos)
	// mg.GenerateMoves()
	castlingSlidingSqBb, err := bitboard.NewBitboard(uint64(0))
	castlingSlidingSqBb.SetBit(int(ht.LookupCastlingSlidingSqByDest[uint64(move.Destination())]))
	if err != nil {
		log.Fatal(err.Error())
	}
	if (*pos).IsAttacked(*castlingSlidingSqBb, legalMoves.AttackedSqsBb()) {
		*pos = tempPos
		return false
	}
	*pos = tempPos
	return true
}

func evaluateMove(move moves.Move, value *int, tempValue *int, root bool, depth int, ply int, pos **position.Position, engineMove **moves.Move, perft *int, singlePlyPerft *int) error {
	if (*pos).IsCastlingMove(move) {
		if !castlingMoveIsValid(move, value, tempValue, root, depth, ply, pos, engineMove, perft, singlePlyPerft) {
			return nil
		}
	}
	(*pos).Move(move)
	legalMoves := generate.GenerateMoves(*pos)
	// mg.GenerateMoves()
	if (*pos).IsAttacked((*pos).InactiveSideKingBb(), legalMoves.AttackedSqsBb()) {
		*pos = (*pos).UnMakeMove()
		return nil
	}
	*tempValue = minMax(depth, ply-1, pos, engineMove, perft, singlePlyPerft)
	if *tempValue > *value {
		value = tempValue
		if root {
			*engineMove = move.CopyMove()
		}
	}
	if root {
		move.Print()
		fmt.Println(*singlePlyPerft)
		*perft += *singlePlyPerft
		*singlePlyPerft = 0
	}
	*pos = (*pos).UnMakeMove()
	return nil
}

func generateMoves(value *int, tempValue *int, root bool, depth int, ply int, pos **position.Position, engineMove **moves.Move, perft *int, singlePlyPerft *int) error {
	legalMoves := generate.GenerateMoves(*pos)
	mvList := legalMoves.GetMovesList()
	for _, move := range mvList {
		evaluateMove(move, value, tempValue, root, depth, ply, pos, engineMove, perft, singlePlyPerft)
	}
	return nil
}
