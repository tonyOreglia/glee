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

type searchParams struct {
	depth          int
	ply            int
	pos            **position.Position
	engineMove     **moves.Move
	perft          *int
	singlePlyPerft *int
}

func minMax(p searchParams) int {
	var value, tempValue int
	root := p.ply == p.depth
	if root {
		fmt.Println("")
		fmt.Print("Perft Divide: ")
		fmt.Println((*p.pos).GetFenString())
	}
	generate.GenerateMoves(*p.pos)
	if p.ply == 0 {
		return evaluate.EvaluatePosition(*p.pos, p.singlePlyPerft)
	}
	if (*p.pos).IsWhitesTurn() {
		value = -30000
		generateMoves(&value, &tempValue, root, p)
		return value
	}
	value = 30000
	generateMoves(&value, &tempValue, root, p)
	return value
}

func generateMoves(value *int, tempValue *int, root bool, p searchParams) error {
	mvList := generate.GenerateMoves(*p.pos).GetMovesList()
	for _, move := range mvList {
		evaluateMove(move, value, tempValue, root, p.depth, p.ply, p.pos, p.engineMove, p.perft, p.singlePlyPerft)
	}
	return nil
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
	*tempValue = minMax(searchParams{
		depth:          depth,
		ply:            ply - 1,
		pos:            pos,
		engineMove:     engineMove,
		perft:          perft,
		singlePlyPerft: singlePlyPerft,
	})
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
