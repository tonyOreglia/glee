// Package position leverages the bitboard library
// to represent a single chess position.
package position

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tonyOreglia/glee/pkg/bitboard"
	"github.com/tonyOreglia/glee/pkg/hashtables"
	"github.com/tonyOreglia/glee/pkg/moves"
)

// don't worry, this does not cause hash tables to be generated more than once
var ht = hashtables.Lookup

// White is the index of White's position bitboards in instance of Position Struct
const White = 0

// Black is the index of Black's position bitboards in instance of Position Struct
const Black = 1

const WhiteKingSideCastlingRightsBit = 62
const WhiteQueenSideCastlingRightsBit = 58
const BlackKingSideCastlingRightsBit = 6
const BlackQueenSideCastlingRightsBit = 2

const OccupiedSqs = 0
const King = 1
const Queen = 2
const Bishops = 3
const Knights = 4
const Rooks = 5
const Pawns = 6

// Position struct represents a static chess position
type Position struct {
	bitboards      [2][]bitboard.Bitboard
	castlingRights []bitboard.Bitboard
	activeSide     int
	enPassanteSq   int
	moveCt         int
	halfMoveCt     int
	previousPos    *Position
}

func StartingPosition() *Position {
	p, _ := NewPositionFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	return p
}

// NewPositionFen constructs Position struct instance from Forth-Edwards Notation string
func NewPositionFen(fen string) (*Position, error) {
	p := new(Position)
	p.bitboards[0] = make([]bitboard.Bitboard, 7)
	p.bitboards[1] = make([]bitboard.Bitboard, 7)
	p.castlingRights = make([]bitboard.Bitboard, 2)
	Position, activeSide, castlingRights, enPassanteSq, moveCount, halfMoveCount := getFenStringTokens(fen)
	p.setBitboardsFromFen(Position, activeSide)
	p.setActiveSide(activeSide)
	p.setCastlingRightsFromFen(castlingRights)
	p.enPassanteSq = enPassanteSq
	p.moveCt = moveCount
	p.halfMoveCt = halfMoveCount
	p.previousPos = nil
	return p, nil
}

func (p *Position) Copy() *Position {
	pCopy := new(Position)
	*pCopy = *p
	pCopy.bitboards[0] = make([]bitboard.Bitboard, 7)
	pCopy.bitboards[1] = make([]bitboard.Bitboard, 7)
	pCopy.castlingRights = make([]bitboard.Bitboard, 2)
	copy(pCopy.bitboards[0], p.bitboards[0])
	copy(pCopy.bitboards[1], p.bitboards[1])
	copy(pCopy.castlingRights, p.castlingRights)
	return pCopy
}

// AllOccupiedSqsBb returns bitboard representing which squares havea piece
func (p *Position) AllOccupiedSqsBb() *bitboard.Bitboard {
	return bitboard.ReturnCombined(&p.bitboards[White][OccupiedSqs], &p.bitboards[Black][OccupiedSqs])
}

// ActiveSideOccupiedSqsBb returns bitboard representing which squares the active side occupies
func (p *Position) ActiveSideOccupiedSqsBb() *bitboard.Bitboard {
	return &p.bitboards[p.activeSide][OccupiedSqs]
}

func (p *Position) IsWhitesTurn() bool {
	return p.activeSide == White
}

func (p *Position) IsBlacksTurn() bool {
	return p.activeSide == Black
}

func (p *Position) EnPassante() int {
	return p.enPassanteSq
}

func (p *Position) InactiveSideOccupiedSqsBb() *bitboard.Bitboard {
	if p.activeSide == White {
		return &p.bitboards[Black][OccupiedSqs]
	}
	return &p.bitboards[White][OccupiedSqs]
}

// GetActiveSide returns the currently active side
func (p *Position) GetActiveSide() int {
	return p.activeSide
}

func (p *Position) GetActiveSideCastlingRightsBb() *bitboard.Bitboard {
	return &p.castlingRights[p.activeSide]
}

// GetActiveSidesBitboards returns the position bitboards for the currently active side
func (p *Position) GetActiveSidesBitboards() []bitboard.Bitboard {
	return p.bitboards[p.activeSide]
}

func (p *Position) ActiveSideKingBb() bitboard.Bitboard {
	return p.bitboards[p.activeSide][King]
}

func (p *Position) InactiveSideKingBb() bitboard.Bitboard {
	if p.activeSide == White {
		return p.bitboards[Black][King]
	}
	return p.bitboards[White][King]
}

func (p *Position) WhiteKingBb() bitboard.Bitboard {
	return p.bitboards[White][King]
}

func (p *Position) BlackKingBb() bitboard.Bitboard {
	return p.bitboards[Black][King]
}

func (p *Position) GetWhiteBitboards() []bitboard.Bitboard {
	return p.bitboards[White]
}

func (p *Position) GetBlackBitboards() []bitboard.Bitboard {
	return p.bitboards[Black]
}

func (p *Position) UnMakeMove() *Position {
	return p.previousPos
}

func (p *Position) Move(mv moves.Move) {
	sideToMove := p.activeSide
	p.MakeMove(mv.Origin(), mv.Destination())
	if mv.PromotionPiece() != 0 {
		p.promotePawn(mv.Destination(), mv.PromotionPiece(), sideToMove)
	}
}

func (p *Position) promotePawn(sq int, piece int, sideToMove int) {
	p.bitboards[sideToMove][piece].SetBit(sq)
	p.bitboards[sideToMove][Pawns].RemoveBit(sq)
	p.updatedOccupiedSqBitboard(sideToMove)
}

func (p *Position) MakeMove(originIndex int, terminusIndex int) {
	p.previousPos = p.Copy()
	// double pawn push move, set en passante
	p.enPassanteSq = 64
	doublePawnPush := p.bitboards[p.activeSide][Pawns].BitIsSet(originIndex) && (terminusIndex-originIndex == -16 || terminusIndex-originIndex == 16)
	if doublePawnPush {
		p.enPassanteSq = (terminusIndex-originIndex)/2 + originIndex
	}
	movingPiece := p.updateMovingSidesBbs(originIndex, terminusIndex)

	if movingPiece == King {
		diff := terminusIndex - originIndex
		// king side castle move
		if diff == 2 {
			_ = p.updateMovingSidesBbs(terminusIndex+1, terminusIndex-1)
		}
		// queen side castle move
		if diff == -2 {
			_ = p.updateMovingSidesBbs(terminusIndex-2, terminusIndex+1)
		}
		p.revokeQueenSideCastlingRight()
		p.revokeKingSideCastlingRight()
	}
	if movingPiece == Rooks {
		if originIndex == 0 || originIndex == 56 {
			p.revokeQueenSideCastlingRight()
		}
		if originIndex == 7 || originIndex == 63 {
			p.revokeKingSideCastlingRight()
		}
	}
	p.updatedOccupiedSqBitboard(p.activeSide)
	p.switchActiveSide()
	attackedPiece := p.removeAttackedPieceFromBbs(terminusIndex)
	if attackedPiece == Rooks {
		if (terminusIndex % 8) == 7 {
			p.revokeKingSideCastlingRight()
		} else if (terminusIndex % 8) == 0 {
			p.revokeQueenSideCastlingRight()
		}
	}
	enPassanteAttack := movingPiece == Pawns && (terminusIndex-originIndex)%8 != 0 && attackedPiece == 0
	if enPassanteAttack {
		capturnedPawnIndex := terminusIndex + 8
		if originIndex >= 32 && originIndex < 40 {
			capturnedPawnIndex = terminusIndex - 8
		}
		p.removeAttackedPieceFromBbs(capturnedPawnIndex)
	}
	p.updatedOccupiedSqBitboard(p.activeSide)
	if p.activeSide == Black {
		p.moveCt++
	}
}

// MakeMove updates position with single chess move
func (p *Position) MakeMoveAlgebraic(origin string, terminus string) {
	originIndex := convertAlgebriacToIndex(origin)
	terminusIndex := convertAlgebriacToIndex(terminus)
	p.MakeMove(originIndex, terminusIndex)
}

func (p *Position) IsAttacked(kingBb bitboard.Bitboard, destSqsBb *bitboard.Bitboard) bool {
	return kingBb.BitwiseAnd(destSqsBb).Value() != uint64(0)
}

func (p *Position) IsCastlingMove(mv moves.Move) bool {
	if p.IsKingMove(mv) {
		if (mv.Destination()-mv.Origin() == 2) || (mv.Destination()-mv.Origin() == -2) {
			return true
		}
	}
	return false
}

func (p *Position) IsKingMove(mv moves.Move) bool {
	activeSideKingBb := p.ActiveSideKingBb()
	originBb := bitboard.NewBitboard(ht.SingleIndexBbHash[mv.Origin()])
	return bitboard.ReturnBitwiseAnd(&activeSideKingBb, originBb).Value() != uint64(0)
}

func (p *Position) switchActiveSide() {
	if p.activeSide == White {
		p.activeSide = Black
	} else {
		p.activeSide = White
	}
}

func (p *Position) removeAttackedPieceFromBbs(terminus int) int {
	switch {
	case p.bitboards[p.activeSide][Pawns].BitIsSet(terminus):
		p.bitboards[p.activeSide][Pawns].RemoveBit(terminus)
		return Pawns
	case p.bitboards[p.activeSide][Rooks].BitIsSet(terminus):
		p.bitboards[p.activeSide][Rooks].RemoveBit(terminus)
		return Rooks
	case p.bitboards[p.activeSide][Knights].BitIsSet(terminus):
		p.bitboards[p.activeSide][Knights].RemoveBit(terminus)
		return Knights
	case p.bitboards[p.activeSide][Bishops].BitIsSet(terminus):
		p.bitboards[p.activeSide][Bishops].RemoveBit(terminus)
		return Bishops
	case p.bitboards[p.activeSide][Queen].BitIsSet(terminus):
		p.bitboards[p.activeSide][Queen].RemoveBit(terminus)
		return Queen
	case p.bitboards[p.activeSide][King].BitIsSet(terminus):
		p.bitboards[p.activeSide][King].RemoveBit(terminus)
		return King
	}
	return 0
}

func (p *Position) updateMovingSidesBbs(origin int, terminus int) int {
	switch {
	case p.bitboards[p.activeSide][Pawns].BitIsSet(origin):
		p.bitboards[p.activeSide][Pawns].SetBit(terminus)
		p.bitboards[p.activeSide][Pawns].RemoveBit(origin)
		return Pawns
	case p.bitboards[p.activeSide][Rooks].BitIsSet(origin):
		p.bitboards[p.activeSide][Rooks].SetBit(terminus)
		p.bitboards[p.activeSide][Rooks].RemoveBit(origin)
		return Rooks
	case p.bitboards[p.activeSide][Knights].BitIsSet(origin):
		p.bitboards[p.activeSide][Knights].SetBit(terminus)
		p.bitboards[p.activeSide][Knights].RemoveBit(origin)
		return Knights
	case p.bitboards[p.activeSide][Bishops].BitIsSet(origin):
		p.bitboards[p.activeSide][Bishops].SetBit(terminus)
		p.bitboards[p.activeSide][Bishops].RemoveBit(origin)
		return Bishops
	case p.bitboards[p.activeSide][Queen].BitIsSet(origin):
		p.bitboards[p.activeSide][Queen].SetBit(terminus)
		p.bitboards[p.activeSide][Queen].RemoveBit(origin)
		return Queen
	case p.bitboards[p.activeSide][King].BitIsSet(origin):
		p.bitboards[p.activeSide][King].SetBit(terminus)
		p.bitboards[p.activeSide][King].RemoveBit(origin)
		return King
	}
	return 0
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

func (p *Position) PrintFen() {
	fmt.Println(p.GetFenString())
}

func (p *Position) convertEnPassanteSqToFenString() string {
	if p.enPassanteSq == 64 {
		return "-"
	}
	return convertIndexToAlgebraic(p.enPassanteSq)
}

func (p *Position) WhiteCanCastleKingSide() bool {
	return p.castlingRights[White].BitIsNotSet(WhiteKingSideCastlingRightsBit)
}

func (p *Position) WhiteCanCastleQueenSide() bool {
	return p.castlingRights[White].BitIsNotSet(WhiteQueenSideCastlingRightsBit)
}

func (p *Position) BlackCanCastleKingSide() bool {
	return p.castlingRights[Black].BitIsNotSet(BlackKingSideCastlingRightsBit)
}

func (p *Position) BlackCanCastleQueenSide() bool {
	return p.castlingRights[Black].BitIsNotSet(BlackQueenSideCastlingRightsBit)
}

func (p *Position) convertCastlingRightsToFenString() string {
	castlingRightsFenString := ""
	if p.castlingRights[White].BitIsNotSet(WhiteKingSideCastlingRightsBit) {
		castlingRightsFenString += "K"
	}
	if p.castlingRights[White].BitIsNotSet(WhiteQueenSideCastlingRightsBit) {
		castlingRightsFenString += "Q"
	}
	if p.castlingRights[Black].BitIsNotSet(BlackKingSideCastlingRightsBit) {
		castlingRightsFenString += "k"
	}
	if p.castlingRights[Black].BitIsNotSet(BlackQueenSideCastlingRightsBit) {
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
	p.activeSide = activeSide
}

func (p *Position) revokeQueenSideCastlingRight() {
	if p.activeSide == White && p.castlingRights[White].BitIsNotSet(WhiteKingSideCastlingRightsBit) {
		p.castlingRights[p.activeSide].SetBit(WhiteQueenSideCastlingRightsBit)
	}
	p.castlingRights[p.activeSide].SetBit(BlackQueenSideCastlingRightsBit)
}

func (p *Position) revokeKingSideCastlingRight() {
	if p.activeSide == White {
		p.castlingRights[p.activeSide].SetBit(WhiteKingSideCastlingRightsBit)
	}
	p.castlingRights[p.activeSide].SetBit(BlackKingSideCastlingRightsBit)
}

func (p *Position) setCastlingRightsFromFen(castlingRights string) {
	p.castlingRights[White].SetBit(WhiteKingSideCastlingRightsBit)
	p.castlingRights[White].SetBit(WhiteQueenSideCastlingRightsBit)
	p.castlingRights[Black].SetBit(BlackKingSideCastlingRightsBit)
	p.castlingRights[Black].SetBit(BlackQueenSideCastlingRightsBit)
	for i := 0; i < len(castlingRights); i++ {
		singleCastlingRight := string(castlingRights[i])
		switch singleCastlingRight {
		case "K":
			p.castlingRights[White].RemoveBit(WhiteKingSideCastlingRightsBit)
		case "Q":
			p.castlingRights[White].RemoveBit(WhiteQueenSideCastlingRightsBit)
		case "k":
			p.castlingRights[Black].RemoveBit(BlackKingSideCastlingRightsBit)
		case "q":
			p.castlingRights[Black].RemoveBit(BlackQueenSideCastlingRightsBit)
		}
	}
}

func (p *Position) setBitboardsFromFen(fenPosition string, activeSide int) {
	var boardIndex int
	for fenStringIndex := 0; fenStringIndex < len(fenPosition); fenStringIndex++ {
		letter := string(fenPosition[fenStringIndex])
		switch letter {
		case "p":
			p.bitboards[Black][Pawns].SetBit(boardIndex)
		case "r":
			p.bitboards[Black][Rooks].SetBit(boardIndex)
		case "n":
			p.bitboards[Black][Knights].SetBit(boardIndex)
		case "b":
			p.bitboards[Black][Bishops].SetBit(boardIndex)
		case "q":
			p.bitboards[Black][Queen].SetBit(boardIndex)
		case "k":
			p.bitboards[Black][King].SetBit(boardIndex)
		case "P":
			p.bitboards[White][Pawns].SetBit(boardIndex)
		case "R":
			p.bitboards[White][Rooks].SetBit(boardIndex)
		case "N":
			p.bitboards[White][Knights].SetBit(boardIndex)
		case "B":
			p.bitboards[White][Bishops].SetBit(boardIndex)
		case "Q":
			p.bitboards[White][Queen].SetBit(boardIndex)
		case "K":
			p.bitboards[White][King].SetBit(boardIndex)
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
		p.bitboards[i][OccupiedSqs].Print()
		fmt.Println("Pawns")
		p.bitboards[i][Pawns].Print()
		fmt.Println("Rooks")
		p.bitboards[i][Rooks].Print()
		fmt.Println("Knights")
		p.bitboards[i][Knights].Print()
		fmt.Println("Bishops")
		p.bitboards[i][Bishops].Print()
		fmt.Println("Queen")
		p.bitboards[i][Queen].Print()
		fmt.Println("King")
		p.bitboards[i][King].Print()
	}

}

func (p *Position) updatedOccupiedSqBitboard(activeSide int) {
	p.bitboards[activeSide][OccupiedSqs].Set(
		p.bitboards[activeSide][Pawns].Value() |
			p.bitboards[activeSide][Rooks].Value() |
			p.bitboards[activeSide][Knights].Value() |
			p.bitboards[activeSide][Bishops].Value() |
			p.bitboards[activeSide][King].Value() |
			p.bitboards[activeSide][Queen].Value())
}

func convertSingleBbRowToFenString(rank int, fenString *string, emptySqs *int, bb [2][]bitboard.Bitboard) {
	for file := int(0); file < 8; file++ {
		convertSingleIndexToString(file, rank, fenString, emptySqs, bb, defaultSingleIndexConversionToFenString)
	}
	if !bb[White][OccupiedSqs].BitIsSet((rank+1)*8-1) && !bb[Black][OccupiedSqs].BitIsSet((rank+1)*8-1) {
		*fenString += strconv.Itoa(*emptySqs)
	}
	notFirstRank := rank != 7
	if notFirstRank {
		*fenString += "/"
	}
	*emptySqs = 0
}

func convertBitboardsToFenString(bb [2][]bitboard.Bitboard) string {
	fenString := ""
	emptySqs := 0
	for rank := int(0); rank < 8; rank++ {
		convertSingleBbRowToFenString(rank, &fenString, &emptySqs, bb)
	}
	return fenString
}

func (p *Position) Print() {
	fmt.Println()
	bb := p.bitboards
	fenString := ""
	emptySqs := 0
	for rank := int(0); rank < 8; rank++ {
		getRowString(rank, &fenString, &emptySqs, bb)
		fenString = ""
	}
	turn := "black"
	if p.activeSide == White {
		turn = "white"
	}
	fmt.Printf("\nmove: %d, turn: %s\ncastling: %s, ep-file: %s\n\n", p.moveCt, turn, p.convertCastlingRightsToFenString(), p.convertEnPassanteSqToFenString())
}

func defaultSingleIndexConversionToCharacter(index int, file int, fenString *string, emptySqs *int, bb [2][]bitboard.Bitboard) {
	*fenString += "."
}

func defaultSingleIndexConversionToFenString(index int, file int, fenString *string, emptySqs *int, bb [2][]bitboard.Bitboard) {
	*emptySqs++
	nextWhiteSqOcc := bb[White][OccupiedSqs].BitIsSet(index + 1)
	nexBlackSqOcc := bb[Black][OccupiedSqs].BitIsSet(index + 1)
	nextSqOcc := nextWhiteSqOcc || nexBlackSqOcc
	notHFile := file != 7
	if nextSqOcc && notHFile {
		*fenString += strconv.Itoa(*emptySqs)
	}
}

func getRowString(rank int, fenString *string, emptySqs *int, bb [2][]bitboard.Bitboard) {
	for file := int(0); file < 8; file++ {
		convertSingleIndexToString(file, rank, fenString, emptySqs, bb, defaultSingleIndexConversionToCharacter)
		*fenString += " "
	}
	if !bb[White][OccupiedSqs].BitIsSet((rank+1)*8-1) && !bb[Black][OccupiedSqs].BitIsSet((rank+1)*8-1) {
		for i := 0; i < *emptySqs; i++ {
			*fenString += "-"
		}
	}
	fmt.Println(*fenString)
	*emptySqs = 0
}

func convertSingleIndexToString(i int, j int, fenString *string, emptySqs *int, bb [2][]bitboard.Bitboard,
	defaultAction func(int, int, *string, *int, [2][]bitboard.Bitboard)) {
	index := int(j*8 + i)
	switch {
	case bb[White][King].BitIsSet(index):
		*fenString += "K"
		*emptySqs = 0
	case bb[White][Queen].BitIsSet(index):
		*fenString += "Q"
		*emptySqs = 0
	case bb[White][Bishops].BitIsSet(index):
		*fenString += "B"
		*emptySqs = 0
	case bb[White][Rooks].BitIsSet(index):
		*fenString += "R"
		*emptySqs = 0
	case bb[White][Knights].BitIsSet(index):
		*fenString += "N"
		*emptySqs = 0
	case bb[White][Pawns].BitIsSet(index):
		*fenString += "P"
		*emptySqs = 0
	case bb[Black][King].BitIsSet(index):
		*fenString += "k"
		*emptySqs = 0
	case bb[Black][Queen].BitIsSet(index):
		*fenString += "q"
		*emptySqs = 0
	case bb[Black][Bishops].BitIsSet(index):
		*fenString += "b"
		*emptySqs = 0
	case bb[Black][Rooks].BitIsSet(index):
		*fenString += "r"
		*emptySqs = 0
	case bb[Black][Knights].BitIsSet(index):
		*fenString += "n"
		*emptySqs = 0
	case bb[Black][Pawns].BitIsSet(index):
		*fenString += "p"
		*emptySqs = 0
	default:
		defaultAction(index, i, fenString, emptySqs, bb)
		// *fenString += ". "
	}
}
