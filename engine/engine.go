package engine

import (
	"fmt"

	"github.com/tonyoreglia/glee/evaluate"
	"github.com/tonyoreglia/glee/generate"
	"github.com/tonyoreglia/glee/hashtables"
	"github.com/tonyoreglia/glee/moves"
	"github.com/tonyoreglia/glee/position"
)

// const Depth = 1

var ht = hashtables.Lookup

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
		mvList := moveGenerator.GetMovesList()
		for _, move := range mvList {
			(*pos).Move(move)
			mg := generate.NewLegalMoveGenerator(*pos)
			mg.GenerateMoves()
			if (*pos).IsAttacked((*pos).WhiteKingBb(), mg.MovesStruct().AttackedSqsBb()) {
				*pos = (*pos).UnMakeMove()
				continue
			}

			tempValue = minMax(depth, ply-1, pos, engineMove, perft, singlePlyPerft)
			if tempValue > value {
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
		}
		return value
	}
	value = 30000
	mvList := moveGenerator.GetMovesList()
	for _, move := range mvList {
		(*pos).Move(move)
		mg := generate.NewLegalMoveGenerator(*pos)
		mg.GenerateMoves()
		if (*pos).IsAttacked((*pos).BlackKingBb(), mg.MovesStruct().AttackedSqsBb()) {
			*pos = (*pos).UnMakeMove()
			continue
		}
		tempValue = minMax(depth, ply-1, pos, engineMove, perft, singlePlyPerft)
		if tempValue < value {
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
	}
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
