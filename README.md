# glee
[![Go Report Card](https://goreportcard.com/badge/github.com/tonyOreglia/glee)](https://goreportcard.com/report/github.com/tonyOreglia/glee)

Tony's golang chess engine.


# To Do List
- [x] Breakdown the project into subtasks

- [x] Bitboard package
  - get
  - set
  - lsb
  - msb

- [x] Position package

- [ ] Generate valid moves
  - Initialise with position

- [ ] Update position

- [ ] Rate position
  - Simple piece wise
  - Positional

- [ ] Search position

# Rules of Chess

Pawn 
 - two forward on move 1 
 - one forward thereafter 
 - cannot move if blocked in front
 - attacks diagonal
 - can attack special en passante square
 - exchanges for any piece if moved to final rank 

Knights 
 - L shape
 - can jump pieces

Bishops, rooks, queen
 - sliding pices
 - cannot jump pieces 

King 
 - sliding piece, any direction, one space 
 - if unmoved, and rook is unmoved, no blocking pieces, not through check, can do special castling move 
 - cannot move into attacked space 
 - must move if attacked 
 - game over if attacked and no safe square to move to or possibility to block or take attacker 


How does this Chess engine calculate moves? 
  mvs.generatePawnMoves()
	mvs.generateKingMoves()
	mvs.generateQueenMoves()
	mvs.generateRookMoves()
	mvs.generateKnightMoves()
	mvs.generateBishopMoves()

