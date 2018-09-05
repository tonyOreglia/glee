package position

import (
	"strings"
	"strconv"
	"log"
	"errors"
	"github.com/tonyoreglia/glee/bitboard"
)

const White = 0
const Black = 1

type Bitboards struct {
	occupiedSqs, king, queen, bishops, knights, rooks, pawns bitboard.Bitboard
}

type CastlingRights struct {
	kingSide, queenSide bool
}

type Position struct {
	bitboards [2]Bitboards
	castlingRights [2]CastlingRights
	activeSide int
	enPassanteSq int
	moveCt int
	halfMoveCt int
}

func NewPositionFen(fen string) (*Position, error) {
	p := new(Position)
	position, activeSide, castlingRights, enPassanteSq, moveCount, halfMoveCount := convertFen(fen)
	p.setBitboardsFromFen(position, activeSide)
	p.setActiveSide(activeSide)
	p.setCastlingRightsFromFen(castlingRights)
	p.enPassanteSq = enPassanteSq
	p.moveCt = moveCount
	p.halfMoveCt = halfMoveCount
	return p, nil
}

func (p *Position) GetFenString() (string) {
	fenPosition := convertBitboardsToFenString(p.bitboards)
	return fenPosition
}

func (p *Position) setActiveSide(activeSide int) {
	p.activeSide = activeSide
}

func (p *Position) setCastlingRightsFromFen(castlingRights string) {
	if castlingRights == "-" { return }
	for i := 0; i < len(castlingRights); i++ {
		singleCastlingRight := string(castlingRights[i])
		switch singleCastlingRight {
			case "K": p.castlingRights[White].kingSide = true
			case "Q": p.castlingRights[White].queenSide = true
			case "k": p.castlingRights[Black].kingSide = true
			case "q": p.castlingRights[Black].queenSide = true
		}
	}
} 

func (p *Position) setBitboardsFromFen(fenPosition string, activeSide int) {
	var boardIndex uint8 = 0
	for fenStringIndex := 0; fenStringIndex < len(fenPosition); fenStringIndex++ {
		letter := string(fenPosition[fenStringIndex])
		switch letter {
				case "p" : p.bitboards[Black].pawns.SetBit(boardIndex)
				case "r" : p.bitboards[Black].rooks.SetBit(boardIndex)
				case "n" : p.bitboards[Black].knights.SetBit(boardIndex)
				case "b" : p.bitboards[Black].bishops.SetBit(boardIndex)
				case "q" : p.bitboards[Black].queen.SetBit(boardIndex)
				case "k" : p.bitboards[Black].king.SetBit(boardIndex)
				case "P" : p.bitboards[White].pawns.SetBit(boardIndex)
				case "R" : p.bitboards[White].rooks.SetBit(boardIndex)
				case "N" : p.bitboards[White].knights.SetBit(boardIndex)
				case "B" : p.bitboards[White].bishops.SetBit(boardIndex)
				case "Q" : p.bitboards[White].queen.SetBit(boardIndex)
				case "K" : p.bitboards[White].king.SetBit(boardIndex)
				case "/" : boardIndex -= 1
				case "1" :
				case "2" : boardIndex += 1
				case "3" : boardIndex += 2
				case "4" : boardIndex += 3
				case "5" : boardIndex += 4
				case "6" : boardIndex += 5
				case "7" : boardIndex += 6
				case "8" : boardIndex += 7
		}
		boardIndex++
	}
}


func convertFen(fen string) (string, int, string, int, int, int) {
	var activeSide int
	fenTokens := strings.Split(fen, " ")
	moveCount, err := strconv.Atoi(fenTokens[4])
	if err != nil { log.Fatal(err) }
	halfMoveCount, err := strconv.Atoi(fenTokens[5])
	if err != nil { log.Fatal(err) }
	enPassnantSq := convertAlgebriacToIndex(fenTokens[3])
	switch fenTokens[1] {
		case "w": activeSide = White
		case "b": activeSide = Black
		default: log.Fatal(errors.New("Active side encoded in Fen must be either 'w' or 'b'"))
	}
	err = validateFen(fenTokens[0], activeSide, fenTokens[2], enPassnantSq, moveCount, halfMoveCount) 
	if err != nil {
		log.Fatal(err)
	}
	return fenTokens[0], activeSide, fenTokens[2], enPassnantSq, moveCount, halfMoveCount 
}

func validateFen(position string, activeSide int, castlingRights string, enPassanteSq int, moveCount int, halfMoveCount int) (error) {
	if moveCount < 0 { return errors.New("Move count encoded in FEN string is less than zero") }
 	if halfMoveCount < 0 { return errors.New("Half move count encoded in FEN string is less than zero") }
	if len(position) > 71 { return errors.New("Position string encoded in FEN is longer than 71 characters in length") }
	if len(position) < 15 { return errors.New("Position string encoded in FEN is shorter than 15 characters in length") }
	if enPassanteSq > 64 { return errors.New("Invalid en passante square encoded in FEN") }
	if activeSide != 0 && activeSide != 1 { return errors.New("Invalide active side encoded in FEN string.") }
	return nil
}

func convertAlgebriacToIndex(algebraic string) int {
	if algebraic == "-" { return 64 }
	column := string(algebraic[0])
	row, _ := strconv.Atoi(string(algebraic[1]))
	row --
	var index, columnValue, rowValue int
	switch column {
			case "a": columnValue = 0;
			case "b": columnValue = 1;
			case "c": columnValue = 2;
			case "d": columnValue = 3;
			case "e": columnValue = 4;
			case "f": columnValue = 5;
			case "g": columnValue = 6;
			case "h":	columnValue = 7;	
	}
	rowValue = (7 - row) * 8
	index = columnValue + rowValue
	return index
}

func convertBitboardsToFenString(bb [2]Bitboards) string {
	fenString := ""
	emptySqs := 0
	for j := uint(0); j < 8; j++ {
		for i := uint(0); i < 8; i++ {
			if (bb[White].king.Get() & (uint64(1) << (j*8+i))) != 0 {fenString += string(emptySqs); fenString += "K"; continue}
			if (bb[White].queen.Get() & (uint64(1) << (j*8+i))) != 0 {fenString += string(emptySqs); fenString += "Q"; continue}
			if (bb[White].bishops.Get() & (uint64(1) << (j*8+i))) != 0 {fenString += string(emptySqs); fenString += "B"; continue}
			if (bb[White].rooks.Get() & (uint64(1) << (j*8+i))) != 0 {fenString += string(emptySqs); fenString += "R"; continue}
			if (bb[White].knights.Get() & (uint64(1) << (j*8+i))) != 0 {fenString += string(emptySqs); fenString += "N"; continue}
			if (bb[White].pawns.Get() & (uint64(1) << (j*8+i))) != 0 {fenString += string(emptySqs); fenString += "P"; continue}
			if (bb[Black].king.Get() & (uint64(1) << (j*8+i))) != 0 {fenString += string(emptySqs); fenString += "k"; continue}
			if (bb[Black].queen.Get() & (uint64(1) << (j*8+i))) != 0 {fenString += string(emptySqs); fenString += "q"; continue}
			if (bb[Black].bishops.Get() & (uint64(1) << (j*8+i))) != 0 {fenString += string(emptySqs); fenString += "b"; continue}
			if (bb[Black].rooks.Get() & (uint64(1) << (j*8+i))) != 0 {fenString += string(emptySqs); fenString += "r"; continue}
			if (bb[Black].knights.Get() & (uint64(1) << (j*8+i))) != 0 {fenString += string(emptySqs); fenString += "n"; continue}
			if (bb[Black].pawns.Get() & (uint64(1) << (j*8+i))) != 0 {fenString += string(emptySqs); fenString += "p"; continue}
			emptySqs++
		}
		fenString += "/"
		emptySqs = 0
	}
	log.Print(fenString)
	return fenString
}
