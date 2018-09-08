// Package position leverages the bitboard library
// to represent a single chess position
package position

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tonyoreglia/glee/bitboard"
)

// White is the index of White's position bitboards in instance of Position Struct
const White = 0

// Black is the index of Black's position bitboards in instance of Position Struct
const Black = 1

type bitboards struct {
	occupiedSqs, king, queen, bishops, knights, rooks, pawns bitboard.Bitboard
}

type castlingRights struct {
	kingSide, queenSide bool
}

// Position struct represents a static chess position
type Position struct {
	bitboards      [2]bitboards
	castlingRights [2]castlingRights
	activeSide     int
	enPassanteSq   int
	moveCt         int
	halfMoveCt     int
}

// NewPositionFen constructs Position struct instance from Forth-Edwards Notation string
func NewPositionFen(fen string) (*Position, error) {
	p := new(Position)
	Position, activeSide, castlingRights, enPassanteSq, moveCount, halfMoveCount := getFenStringTokens(fen)
	p.setBitboardsFromFen(Position, activeSide)
	p.setActiveSide(activeSide)
	p.setCastlingRightsFromFen(castlingRights)
	p.enPassanteSq = enPassanteSq
	p.moveCt = moveCount
	p.halfMoveCt = halfMoveCount
	return p, nil
}

// MakeMove updates position with single chess move
func (p *Position) MakeMove(origin string, terminus string, activeSide int) {
	originIndex := convertAlgebriacToIndex(origin)
	terminusIndex := convertAlgebriacToIndex(terminus)
	p.updateBbsSingleMove(originIndex, terminusIndex, activeSide)
	p.switchActiveSide()
	if p.activeSide == Black {
		p.moveCt++
	}
	// need to update half move count (non capture or pawn plys)
	// need to update en passante
}

func (p *Position) switchActiveSide() {
	if p.activeSide == White {
		p.activeSide = Black
	} else {
		p.activeSide = White
	}
}

func (p *Position) updateBbsSingleMove(origin int, terminus int, activeSide int) {
	switch {
	case p.bitboards[activeSide].pawns.IsBitSet(origin):
		p.bitboards[activeSide].pawns.SetBit(terminus)
		p.bitboards[activeSide].pawns.RemoveBit(origin)
	case p.bitboards[activeSide].rooks.IsBitSet(origin):
		p.bitboards[activeSide].rooks.SetBit(terminus)
		p.bitboards[activeSide].rooks.RemoveBit(origin)
	case p.bitboards[activeSide].knights.IsBitSet(origin):
		p.bitboards[activeSide].knights.SetBit(terminus)
		p.bitboards[activeSide].knights.RemoveBit(origin)
	case p.bitboards[activeSide].bishops.IsBitSet(origin):
		p.bitboards[activeSide].bishops.SetBit(terminus)
		p.bitboards[activeSide].bishops.RemoveBit(origin)
	case p.bitboards[activeSide].queen.IsBitSet(origin):
		p.bitboards[activeSide].queen.SetBit(terminus)
		p.bitboards[activeSide].queen.RemoveBit(origin)
	case p.bitboards[activeSide].king.IsBitSet(origin):
		p.bitboards[activeSide].king.SetBit(terminus)
		p.bitboards[activeSide].king.RemoveBit(origin)
	}
	p.updatedOccupiedSqBitboard(activeSide)
}

// GetFenString converts position instance to it's Forsyth–Edwards Notation string
func (p *Position) GetFenString() string {
	fenPosition := convertBitboardsToFenString(p.bitboards)
	activeSideString, _ := p.convertActiveSideToString()
	castlingRightsString := p.convertCastlingRightsToFenString()
	enPassanteSqFenString := p.convertEnPassanteSqToFenString()
	fenPosition += " " + activeSideString +
		" " + castlingRightsString +
		" " + enPassanteSqFenString +
		" " + strconv.Itoa(p.moveCt) +
		" " + strconv.Itoa(p.halfMoveCt)
	return fenPosition
}

func (p *Position) convertEnPassanteSqToFenString() string {
	if p.enPassanteSq == 64 {
		return "-"
	}
	return convertIndexToAlgebraic(p.enPassanteSq)
}

func (p *Position) convertCastlingRightsToFenString() string {
	castlingRightsFenString := ""
	if p.castlingRights[White].kingSide {
		castlingRightsFenString += "K"
	}
	if p.castlingRights[White].queenSide {
		castlingRightsFenString += "Q"
	}
	if p.castlingRights[Black].kingSide {
		castlingRightsFenString += "k"
	}
	if p.castlingRights[Black].queenSide {
		castlingRightsFenString += "q"
	}
	if castlingRightsFenString == "" {
		castlingRightsFenString = "-"
	}
	return castlingRightsFenString
}

func (p *Position) convertActiveSideToString() (string, error) {
	if p.activeSide == White {
		return "w", nil
	}
	if p.activeSide == Black {
		return "b", nil
	}
	return "", fmt.Errorf("position class ActiveSide value invalid: %v", p.activeSide)
}

func (p *Position) setActiveSide(activeSide int) {
	fmt.Println(activeSide)
	p.activeSide = activeSide
}

func (p *Position) setCastlingRightsFromFen(castlingRights string) {
	if castlingRights == "-" {
		return
	}
	for i := 0; i < len(castlingRights); i++ {
		singleCastlingRight := string(castlingRights[i])
		switch singleCastlingRight {
		case "K":
			p.castlingRights[White].kingSide = true
		case "Q":
			p.castlingRights[White].queenSide = true
		case "k":
			p.castlingRights[Black].kingSide = true
		case "q":
			p.castlingRights[Black].queenSide = true
		}
	}
}

func (p *Position) setBitboardsFromFen(fenPosition string, activeSide int) {
	var boardIndex int
	for fenStringIndex := 0; fenStringIndex < len(fenPosition); fenStringIndex++ {
		letter := string(fenPosition[fenStringIndex])
		switch letter {
		case "p":
			p.bitboards[Black].pawns.SetBit(boardIndex)
		case "r":
			p.bitboards[Black].rooks.SetBit(boardIndex)
		case "n":
			p.bitboards[Black].knights.SetBit(boardIndex)
		case "b":
			p.bitboards[Black].bishops.SetBit(boardIndex)
		case "q":
			p.bitboards[Black].queen.SetBit(boardIndex)
		case "k":
			p.bitboards[Black].king.SetBit(boardIndex)
		case "P":
			p.bitboards[White].pawns.SetBit(boardIndex)
		case "R":
			p.bitboards[White].rooks.SetBit(boardIndex)
		case "N":
			p.bitboards[White].knights.SetBit(boardIndex)
		case "B":
			p.bitboards[White].bishops.SetBit(boardIndex)
		case "Q":
			p.bitboards[White].queen.SetBit(boardIndex)
		case "K":
			p.bitboards[White].king.SetBit(boardIndex)
		case "/":
			boardIndex--
		case "1":
		case "2":
			boardIndex++
		case "3":
			boardIndex += 2
		case "4":
			boardIndex += 3
		case "5":
			boardIndex += 4
		case "6":
			boardIndex += 5
		case "7":
			boardIndex += 6
		case "8":
			boardIndex += 7
		}
		boardIndex++
	}
	p.updatedOccupiedSqBitboard(White)
	p.updatedOccupiedSqBitboard(Black)
}

func getFenStringTokens(fen string) (string, int, string, int, int, int) {
	var activeSide int
	fenTokens := strings.Split(fen, " ")
	moveCount, err := strconv.Atoi(fenTokens[4])
	if err != nil {
		log.Fatal(err)
	}
	halfMoveCount, err := strconv.Atoi(fenTokens[5])
	if err != nil {
		log.Fatal(err)
	}
	enPassnantSq := convertAlgebriacToIndex(fenTokens[3])
	switch fenTokens[1] {
	case "w":
		activeSide = White
	case "b":
		activeSide = Black
	default:
		log.Fatal(errors.New("Active side encoded in Fen must be either 'w' or 'b'"))
	}
	err = validateFenTokens(fenTokens[0], activeSide, fenTokens[2], enPassnantSq, moveCount, halfMoveCount)
	if err != nil {
		log.Fatal(err)
	}
	return fenTokens[0], activeSide, fenTokens[2], enPassnantSq, moveCount, halfMoveCount
}

func validateFenTokens(Position string, activeSide int, castlingRights string, enPassanteSq int, moveCount int, halfMoveCount int) error {
	if moveCount < 0 {
		return errors.New("Move count encoded in FEN string is less than zero")
	}
	if halfMoveCount < 0 {
		return errors.New("Half move count encoded in FEN string is less than zero")
	}
	if len(Position) > 71 {
		return errors.New("Position string encoded in FEN is longer than 71 characters in length")
	}
	if len(Position) < 15 {
		return errors.New("Position string encoded in FEN is shorter than 15 characters in length")
	}
	if enPassanteSq > 64 {
		return errors.New("Invalid en passante square encoded in FEN")
	}
	if activeSide != 0 && activeSide != 1 {
		return errors.New("invalid active side encoded in FEN string")
	}
	return nil
}

func convertAlgebriacToIndex(algebraic string) int {
	if algebraic == "-" {
		return 64
	}
	column := string(algebraic[0])
	row, _ := strconv.Atoi(string(algebraic[1]))
	row--
	var index, columnValue, rowValue int
	switch column {
	case "a":
		columnValue = 0
	case "b":
		columnValue = 1
	case "c":
		columnValue = 2
	case "d":
		columnValue = 3
	case "e":
		columnValue = 4
	case "f":
		columnValue = 5
	case "g":
		columnValue = 6
	case "h":
		columnValue = 7
	}
	rowValue = ((7 - row) * 8)
	index = columnValue + rowValue
	return index
}

func convertIndexToAlgebraic(index int) string {
	var algebraic string
	column := index % 8
	row := 8 - (index / 8)
	switch column {
	case 0:
		algebraic = "a"
	case 1:
		algebraic = "b"
	case 2:
		algebraic = "c"
	case 3:
		algebraic = "d"
	case 4:
		algebraic = "e"
	case 5:
		algebraic = "f"
	case 6:
		algebraic = "g"
	case 7:
		algebraic = "h"
	}
	algebraic += strconv.Itoa(row)
	return algebraic
}

// PrintBitboards prints all Position bitboards in 8x8 array
func (p *Position) PrintBitboards() {
	for i := 0; i < 2; i++ {
		if i == 0 {
			fmt.Println(":::::White::::::")
		} else {
			fmt.Println(":::::Black::::")
		}
		fmt.Println("Occupied Squares")
		p.bitboards[i].occupiedSqs.Print()
		fmt.Println("Pawns")
		p.bitboards[i].pawns.Print()
		fmt.Println("Rooks")
		p.bitboards[i].rooks.Print()
		fmt.Println("Knights")
		p.bitboards[i].knights.Print()
		fmt.Println("Bishops")
		p.bitboards[i].bishops.Print()
		fmt.Println("queen")
		p.bitboards[i].queen.Print()
		fmt.Println("king")
		p.bitboards[i].king.Print()
	}

}

func (p *Position) updatedOccupiedSqBitboard(activeSide int) {
	p.bitboards[activeSide].occupiedSqs.Set(
		p.bitboards[activeSide].pawns.Value() |
			p.bitboards[activeSide].rooks.Value() |
			p.bitboards[activeSide].knights.Value() |
			p.bitboards[activeSide].bishops.Value() |
			p.bitboards[activeSide].king.Value() |
			p.bitboards[activeSide].queen.Value())
}

func convertSingleBbIndexToFen(i int, j int, fenString *string, emptySqs *int, bb [2]bitboards) {
	index := int(j*8 + i)
	switch {
	case bb[White].king.IsBitSet(index):
		*fenString += "K"
		*emptySqs = 0
	case bb[White].queen.IsBitSet(index):
		*fenString += "Q"
		*emptySqs = 0
	case bb[White].bishops.IsBitSet(index):
		*fenString += "B"
		*emptySqs = 0
	case bb[White].rooks.IsBitSet(index):
		*fenString += "R"
		*emptySqs = 0
	case bb[White].knights.IsBitSet(index):
		*fenString += "N"
		*emptySqs = 0
	case bb[White].pawns.IsBitSet(index):
		*fenString += "P"
		*emptySqs = 0
	case bb[Black].king.IsBitSet(index):
		*fenString += "k"
		*emptySqs = 0
	case bb[Black].queen.IsBitSet(index):
		*fenString += "q"
		*emptySqs = 0
	case bb[Black].bishops.IsBitSet(index):
		*fenString += "b"
		*emptySqs = 0
	case bb[Black].rooks.IsBitSet(index):
		*fenString += "r"
		*emptySqs = 0
	case bb[Black].knights.IsBitSet(index):
		*fenString += "n"
		*emptySqs = 0
	case bb[Black].pawns.IsBitSet(index):
		*fenString += "p"
		*emptySqs = 0
	default:
		*emptySqs++
		nextWhiteSqOcc := bb[White].occupiedSqs.IsBitSet(index + 1)
		nexBlackSqOcc := bb[Black].occupiedSqs.IsBitSet(index + 1)
		nextSqOcc := nextWhiteSqOcc || nexBlackSqOcc
		notHFile := i != 7
		if nextSqOcc && notHFile {
			*fenString += strconv.Itoa(*emptySqs)
		}
	}
}

func convertSingleBbRowToFenString(rank int, fenString *string, emptySqs *int, bb [2]bitboards) {
	for file := int(0); file < 8; file++ {
		convertSingleBbIndexToFen(file, rank, fenString, emptySqs, bb)
	}
	if !bb[White].occupiedSqs.IsBitSet((rank+1)*8-1) && !bb[Black].occupiedSqs.IsBitSet((rank+1)*8-1) {
		*fenString += strconv.Itoa(*emptySqs)
	}
	notFirstRank := rank != 7
	if notFirstRank {
		*fenString += "/"
	}
	*emptySqs = 0
}

func convertBitboardsToFenString(bb [2]bitboards) string {
	fenString := ""
	emptySqs := 0
	for rank := int(0); rank < 8; rank++ {
		convertSingleBbRowToFenString(rank, &fenString, &emptySqs, bb)
	}
	return fenString
}
