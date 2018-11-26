package engine

import (
	"fmt"
	"log"

	"github.com/tonyoreglia/glee/bitboard"
	"github.com/tonyoreglia/glee/evaluate"
	"github.com/tonyoreglia/glee/generate"
	"github.com/tonyoreglia/glee/hashtables"
	"github.com/tonyoreglia/glee/moves"
	"github.com/tonyoreglia/glee/position"
)

// const Depth = 1

var ht = hashtables.Lookup

func castlingMoveIsValid(move moves.Move, value *int, tempValue *int, root bool, depth int, ply int, pos **position.Position, engineMove **moves.Move, perft *int, singlePlyPerft *int) bool {
	tempPos := *pos
	*pos = (*pos).UnMakeMove()
	mg := generate.NewLegalMoveGenerator(*pos)
	mg.GenerateMoves()
	castlingSlidingSqBb, err := bitboard.NewBitboard(uint64(0))
	castlingSlidingSqBb.SetBit(int(ht.LookupCastlingSlidingSqByDest[uint64(move.Destination())]))
	if err != nil {
		log.Fatal(err.Error())
	}
	if (*pos).IsAttacked(*castlingSlidingSqBb, mg.MovesStruct().AttackedSqsBb()) {
		*pos = tempPos
		return false
	}
	*pos = tempPos
	return true
}

func evaluateMove(move moves.Move, value *int, tempValue *int, root bool, depth int, ply int, pos **position.Position, engineMove **moves.Move, perft *int, singlePlyPerft *int) error {
	var mg *generate.LegalMoveGenerator
	if (*pos).IsCastlingMove(move) {
		if !castlingMoveIsValid(move, value, tempValue, root, depth, ply, pos, engineMove, perft, singlePlyPerft) {
			return nil
		}
	}
	(*pos).Move(move)
	mg = generate.NewLegalMoveGenerator(*pos)
	mg.GenerateMoves()
	if (*pos).IsAttacked((*pos).InactiveSideKingBb(), mg.MovesStruct().AttackedSqsBb()) {
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
	moveGenerator := generate.NewLegalMoveGenerator(*pos)
	moveGenerator.GenerateMoves()
	mvList := moveGenerator.GetMovesList()
	for _, move := range mvList {
		evaluateMove(move, value, tempValue, root, depth, ply, pos, engineMove, perft, singlePlyPerft)
	}
	return nil
}

func minMax(depth int, ply int, pos **position.Position, engineMove **moves.Move, perft *int, singlePlyPerft *int) int {
	var value, tempValue int
	root := ply == depth
	if root {
		fmt.Println("")
		fmt.Print("Perft Divide: ")
		fmt.Println((*pos).GetFenString())
	}
	moveGenerator := generate.NewLegalMoveGenerator(*pos)
	moveGenerator.GenerateMoves()
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

// int game::alpha_beta_max(int alpha, int beta, char depth, int start) {
// 	bool no_moves = true;
// 	if(depth == 0) return evaluate(); //even: light. odd: dark
// 	bool root = depth == DEPTH;
// 	int total = position->move_t.size();
// 	if(root) {
// 			move_count = position->move_t.size();
// 	}
// 	for(int j = start; j < total; j++) {
// 			if(make_move(j)) {
// 					no_moves = false;
// 					score = alpha_beta_min(alpha, beta, depth - 1, total);
// 					unmake_move();
// 					position->clear_moves(total, position->move_t.size());
// 					if(score >= beta)
// 							return beta;
// 					if(score > alpha) {
// 							alpha = score;
// 							if(root) engine_move = j;
// 					}
// 			}
// 	}
// 	if(no_moves) {
// 			if(mate_check()) {
// 					return (-25000 - depth);
// 			}
// 	}
// 	return alpha;
// }

// int game::alpha_beta_min(int alpha, int beta, char depth, int start) {
// 	bool no_moves = true;
// 	if( depth == 0) return evaluate(); //even: dark. odd: light
// 	bool root = depth == DEPTH;
// 	int total = position->move_t.size();
// 	if(root) {
// 			move_count = position->move_t.size();
// 	}
// 	for(int j=start; j < total; j++) {
// 			if(make_move(j)) {
// 					no_moves = false;
// 					score = alpha_beta_max(alpha, beta, depth - 1, total);
// 					unmake_move();
// 					position->clear_moves(total, position->move_t.size());
// 					if(score <= alpha)
// 							return alpha;
// 					if(score < beta) {
// 							beta = score;
// 							if(root) engine_move = j;
// 					}
// 			}
// 	}
// 	if(no_moves) {
// 			if(mate_check()) {
// 					return (25000 + depth);
// 			}
// 					//return 1000 + depth;
// 	}
// 	return beta;
// }
