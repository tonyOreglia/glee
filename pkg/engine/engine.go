package engine

import (
	"github.com/tonyOreglia/glee/pkg/bitboard"
	"github.com/tonyOreglia/glee/pkg/evaluate"
	"github.com/tonyOreglia/glee/pkg/generate"
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/moves"
	"github.com/tonyOreglia/glee/pkg/position"
)

var ht = hashtables.Lookup

type SearchParams struct {
	Depth           int
	Ply             int
	Pos             **position.Position
	EngineMove      *moves.Move
	Perft           *int
	SinglePlyPerft  *int
	EvaluationScore int
	Root            bool
}

func MinMax(p SearchParams) int {
	p.Root = p.Ply == p.Depth
	if p.Ply == 0 {
		*p.SinglePlyPerft++
		return evaluate.EvaluatePosition(*p.Pos)
	}
	p.EvaluationScore = 30000
	if (*p.Pos).IsWhitesTurn() {
		p.EvaluationScore = -30000
	}
	mvs := generate.GenerateMoves(*p.Pos)
	mvList := mvs.GetMovesList()
	for _, move := range mvList {
		if !MakeValidMove(move, p.Pos) {
			continue
		}
		evaluateMove(move, p)
	}
	return p.EvaluationScore
}

// evaluateMove checks if pseudo legal move is valid, and if so
// increments Ply and calls MinMax to continue the search
func evaluateMove(move moves.Move, p SearchParams) error {
	p.Ply = p.Ply - 1
	temp := MinMax(p)
	if temp > p.EvaluationScore {
		p.EvaluationScore = temp
		if p.Root {
			*p.EngineMove = move
		}
	}
	if p.Root {
		// move.Print()
		// fmt.Println(*p.SinglePlyPerft)
		*p.Perft += *p.SinglePlyPerft
		*p.SinglePlyPerft = 0
	}
	*p.Pos = (*p.Pos).UnMakeMove()
	return nil
}

func MakeValidMove(move moves.Move, pos **position.Position) bool {
	if (*pos).IsCastlingMove(move) {
		(*pos).Move(move)
		if !castlingMoveIsValid(move, pos) {
			*pos = (*pos).UnMakeMove()
			return false
		}
	} else {
		(*pos).Move(move)
	}
	legalMoves := generate.GenerateMoves(*pos)
	if (*pos).IsAttacked((*pos).InactiveSideKingBb(), legalMoves.AttackedSqsBb()) {
		*pos = (*pos).UnMakeMove()
		return false
	}
	return true
}

func castlingMoveIsValid(move moves.Move, pos **position.Position) bool {
	kingPosition := bitboard.NewBitboardFromIndex(move.Origin())
	legalMoves := generate.GenerateMoves(*pos)
	castlingSlidingSqBb := bitboard.NewBitboard(uint64(0))
	castlingSlidingSqBb.SetBit(int(ht.LookupCastlingSlidingSqByDest[uint64(move.Destination())]))
	castlingSlidingSqBb.Combine(kingPosition)
	pawnAttacks := generate.GeneratePotentialPawnAttacks(*pos, ht)
	if (*pos).IsAttacked(*castlingSlidingSqBb, bitboard.ReturnCombined(legalMoves.AttackedSqsBb(), pawnAttacks.AttackedSqsBb())) {
		return false
	}
	return true
}

func AlphaBetaMax(alpha int, beta int, ply int, p SearchParams) int {
	noMoves := true
	if ply == 0 {
		return evaluate.EvaluatePosition(*p.Pos)
	}
	p.Root = ply == p.Depth
	mvs := generate.GenerateMoves(*p.Pos).GetMovesList()
	for _, move := range mvs {
		if MakeValidMove(move, p.Pos) {
			noMoves = false
			score := AlphaBetaMin(alpha, beta, ply-1, p)
			*p.Pos = (*p.Pos).UnMakeMove()
			if score >= beta {
				return beta
			}
			if score > alpha {
				alpha = score
				if p.Root {
					*p.EngineMove = move
				}
			}
		}
	}
	if noMoves {
		// check for mate situation
		return -2500
	}
	return alpha
}

func AlphaBetaMin(alpha int, beta int, ply int, p SearchParams) int {
	noMoves := true
	if ply == 0 {
		return evaluate.EvaluatePosition(*p.Pos)
	}
	p.Root = ply == p.Depth
	mvs := generate.GenerateMoves(*p.Pos).GetMovesList()
	for _, move := range mvs {
		if MakeValidMove(move, p.Pos) {
			noMoves = false
			score := AlphaBetaMax(alpha, beta, ply-1, p)
			*p.Pos = (*p.Pos).UnMakeMove()
			if score <= alpha {
				return alpha
			}
			if score < beta {
				beta = score
				if p.Root {
					*p.EngineMove = move
				}
			}
		}
	}
	if noMoves {
		// check for mate situation
		return 2500
	}
	return beta
}
