package generatemoves

import (
	"fmt"

	"github.com/tonyoreglia/glee/chessmoves"
	"github.com/tonyoreglia/glee/evaluate"
	"github.com/tonyoreglia/glee/hashtables"
	"github.com/tonyoreglia/glee/position"
)

// const Depth = 1

var ht = hashtables.CalculateAllLookupBbs()

func minMax(depth int, ply int, pos *position.Position, engineMove **chessmoves.Move, perft *int, singlePlyPerft *int) int {
	var value, tempValue int
	root := ply == depth
	moveList := chessmoves.NewLegalMoves(pos, ht)
	moveList.GenerateMoves()
	if ply == 0 {
		return evaluate.EvaluatePosition(pos, singlePlyPerft)
	}
	if pos.IsWhitesTurn() {
		value = -30000
		mvList := moveList.GetMovesList()
		for _, move := range mvList {
			pos.MakeMove(move.GetOrigin(), move.GetDestination(), pos.GetActiveSide())
			tempValue = minMax(depth, ply-1, pos, engineMove, perft, singlePlyPerft)
			if tempValue > value {
				value = tempValue
				if root {
					*engineMove = move.CopyMove()
				}
			}
			if root {
				fmt.Println("move: ", move, "-- perft: ", *singlePlyPerft)
				*perft += *singlePlyPerft
				*singlePlyPerft = 0
			}
			pos.UnMakeMove()
		}
		return value
	}
	value = 30000
	mvList := moveList.GetMovesList()
	for _, move := range mvList {
		pos.MakeMove(move.GetOrigin(), move.GetDestination(), pos.GetActiveSide())
		tempValue = minMax(depth, ply-1, pos, engineMove, perft, singlePlyPerft)
		if tempValue < value {
			value = tempValue
			if root {
				*engineMove = move.CopyMove()
			}
		}
		pos.UnMakeMove()
	}
	return value
}

// func searchMax(depth int, pos *position.Position) {
// 	if  depth == 0 {
// 		return evaluate.EvaluatePosition()
// 	}

// 	max := -30000
// 	root := depth == Depth

// 	// total := position->move_t.size()

// 	// if root {
// 	// 		move_count = position->move_t.size();
// 	// }
// 	// if(total == start) {
// 	// 		return -3000;
// 	// }
// 	for j := 0; j < total; j++ {
// 			if(make_move(j)) {
// 					if(root) {
// 							position->print_algebraic(position->game_t.size()-1);
// 							std::cout << "\t\t";
// 					}
// 					score = search_min(depth - 1, total);
// 					if(root) {
// 							std::cout << perft << std::endl;
// 							perft_total += perft;
// 							perft = 0;
// 					}
// 					unmake_move();
// 					position->clear_moves(total, position->move_t.size());
// 					if(score > max) {
// 							max = score;
// 							if(root) engine_move = j;
// 					}
// 			}
// 	}
// 	return max;
// }

// int game::search_min(char depth, int start) {
// 	if( depth == 0) return evaluate();
// 	int min = 30000;
// 	bool root = depth == DEPTH;
// 	int total = position->move_t.size();
// 	if(root) {
// 			move_count = position->move_t.size();
// 	}
// 	for(int j=start; j < total; j++) {
// 			if(make_move(j)) {
// 					if(root) {
// 							position->print_algebraic(position->game_t.size()-1);
// 							std::cout << "\t\t";
// 					}
// 					score = search_max(depth - 1, total);
// 					if(root) {
// 							std::cout << perft << std::endl;
// 							perft_total += perft;
// 							perft = 0;
// 					}
// 					unmake_move();
// 					position->clear_moves(total, position->move_t.size());
// 					if(score < min) {
// 							min = score;
// 							if(root) engine_move = j;
// 					}
// 			}
// 	}
// 	return min;
// }

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
