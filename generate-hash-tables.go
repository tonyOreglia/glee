package main

import (
	"fmt"

	"github.com/tonyoreglia/glee/bitboard"
)

var AfileBb uint64 = 0x101010101010101
var BfileBb uint64 = AfileBb << 1
var CfileBb uint64 = AfileBb << 2
var DfileBb uint64 = AfileBb << 3
var EfileBb uint64 = AfileBb << 4
var FfileBb uint64 = AfileBb << 5
var GfileBb uint64 = AfileBb << 6
var HfileBb uint64 = AfileBb << 7
var FourthRankBb uint64 = 0xFF00000000
var FifthRankBb uint64 = 0xFF000000

var SingleIndexBbHash [64]uint64
var EnPassantBbHash [64]uint64
var AttackedEnPassantPawnLocationBbHash [64]uint64
var NorthArrayBbHash [64]uint64
var SouthArrayBbHash [64]uint64
var EastArrayBbHash [64]uint64
var WestArrayBbHash [64]uint64
var NorthEastArrayBbHash [64]uint64
var NorthWestArrayBbHash [64]uint64
var SouthEastArrayBbHash [64]uint64
var SouthWestArrayBbHash [64]uint64
var KnightAttackBbHash [64]uint64
var LegalKingMovesBbHash [2][64]uint64
var LegalPawnMovesBbHash [2][64]uint64

func calculateAllLookupBbs() {
	for index := 0; index < 64; index++ {
		EnPassantBbHash[index] = uint64(0)
		AttackedEnPassantPawnLocationBbHash[index] = uint64(0)
		NorthArrayBbHash[index] = uint64(0)
		SouthArrayBbHash[index] = uint64(0)
		EastArrayBbHash[index] = uint64(0)
		WestArrayBbHash[index] = uint64(0)
		NorthEastArrayBbHash[index] = uint64(0)
		NorthWestArrayBbHash[index] = uint64(0)
		SouthEastArrayBbHash[index] = uint64(0)
		SouthWestArrayBbHash[index] = uint64(0)
		KnightAttackBbHash[index] = uint64(0)
	}
	generateSingleBitLookup()
	generateArrayBitboardLookup()
	generateEnPassantBitboardLookup()
	printAllBitboards()
}

func generateSingleBitLookup() {
	for i := uint(0); i < 64; i++ {
		SingleIndexBbHash[i] = uint64(1) << i
	}
}

func generateEnPassantBitboardLookup() {
	for i := 0; i < 64; i++ {
		EnPassantBbHash[i] = 0
		AttackedEnPassantPawnLocationBbHash[i] = 0
		switch {
		case i > 23 && i < 32:
			EnPassantBbHash[i] = SingleIndexBbHash[i-8]
		case i > 31 && i < 40:
			EnPassantBbHash[i] = SingleIndexBbHash[i+8]
		case i > 15 && i < 24:
			AttackedEnPassantPawnLocationBbHash[i] = SingleIndexBbHash[i+8]
		case i > 39 && i < 48:
			AttackedEnPassantPawnLocationBbHash[i] = SingleIndexBbHash[i-8]
		}
	}
}

func generateArrayBitboardLookup() {
	for index := 0; index < 64; index++ {
		northOfIndex := index
		for northOfIndex > 7 {
			NorthArrayBbHash[index] |= SingleIndexBbHash[northOfIndex-8]
			northOfIndex -= 8
		}
		southOfIndex := index
		for southOfIndex < 56 {
			SouthArrayBbHash[index] |= SingleIndexBbHash[southOfIndex+8]
			southOfIndex += 8
		}
		eastOfIndex := index
		for (SingleIndexBbHash[eastOfIndex])&HfileBb == 0 {
			EastArrayBbHash[index] |= SingleIndexBbHash[eastOfIndex+1]
			eastOfIndex++
		}
		westOfIndex := index
		for (SingleIndexBbHash[westOfIndex])&AfileBb == 0 {
			WestArrayBbHash[index] |= SingleIndexBbHash[westOfIndex-1]
			westOfIndex--
		}
		northEastOfIndex := index
		for (SingleIndexBbHash[northEastOfIndex])&HfileBb == 0 && northEastOfIndex > 7 {
			NorthEastArrayBbHash[index] |= SingleIndexBbHash[northEastOfIndex-7]
			northEastOfIndex -= 7
		}
		northWestOfIndex := index
		for (SingleIndexBbHash[northWestOfIndex])&AfileBb == 0 && northWestOfIndex > 8 {
			NorthWestArrayBbHash[index] |= SingleIndexBbHash[northWestOfIndex-9]
			northWestOfIndex -= 9
		}
		southEastOfIndex := index
		for (SingleIndexBbHash[southEastOfIndex])&HfileBb == 0 && southEastOfIndex < 56 {
			SouthEastArrayBbHash[index] |= SingleIndexBbHash[southEastOfIndex+9]
			southEastOfIndex += 9
		}
		southWestOfIndex := index
		for (SingleIndexBbHash[southWestOfIndex])&AfileBb == 0 && southWestOfIndex < 56 {
			SouthWestArrayBbHash[index] |= SingleIndexBbHash[southWestOfIndex+7]
			southWestOfIndex += 7
		}
		if (index+10) < 64 && SingleIndexBbHash[index]&(HfileBb|GfileBb) == 0 {
			KnightAttackBbHash[index] |= SingleIndexBbHash[index+10]
		}
		if (index+6) < 64 && SingleIndexBbHash[index]&(AfileBb|BfileBb) == 0 {
			KnightAttackBbHash[index] |= SingleIndexBbHash[index+6]
		}
		if (index+17) < 64 && SingleIndexBbHash[index]&HfileBb == 0 {
			KnightAttackBbHash[index] |= SingleIndexBbHash[index+17]
		}
		if (index+15) < 64 && SingleIndexBbHash[index]&AfileBb == 0 {
			KnightAttackBbHash[index] |= SingleIndexBbHash[index+15]
		}
		if (index-10) >= 0 && SingleIndexBbHash[index]&(AfileBb|BfileBb) == 0 {
			KnightAttackBbHash[index] |= SingleIndexBbHash[index-10]
		}
		if (index-6) >= 0 && SingleIndexBbHash[index]&(HfileBb|GfileBb) == 0 {
			KnightAttackBbHash[index] |= SingleIndexBbHash[index-6]
		}
		if (index-17) >= 0 && SingleIndexBbHash[index]&AfileBb == 0 {
			KnightAttackBbHash[index] |= SingleIndexBbHash[index-17]
		}
		if (index-15) >= 0 && SingleIndexBbHash[index]&HfileBb == 0 {
			KnightAttackBbHash[index] |= SingleIndexBbHash[index-15]
		}

		LegalKingMovesBbHash[0][index] = 0
		if index != 63 {
			LegalKingMovesBbHash[0][index] |= SingleIndexBbHash[index+1]
		}
		if index != 0 {
			LegalKingMovesBbHash[0][index] |= SingleIndexBbHash[index-1]
		}
		if index <= 55 {
			LegalKingMovesBbHash[0][index] |= SingleIndexBbHash[index+8]
		}
		if index >= 8 {
			LegalKingMovesBbHash[0][index] |= SingleIndexBbHash[index-8]
		}
		if index <= 54 {
			LegalKingMovesBbHash[0][index] |= SingleIndexBbHash[index+9]
		}
		if index >= 9 {
			LegalKingMovesBbHash[0][index] |= SingleIndexBbHash[index-9]
		}
		if index <= 56 {
			LegalKingMovesBbHash[0][index] |= SingleIndexBbHash[index+7]
		}
		if index >= 7 {
			LegalKingMovesBbHash[0][index] |= SingleIndexBbHash[index-7]
		}
		LegalKingMovesBbHash[1][index] = LegalKingMovesBbHash[0][index]
		if SingleIndexBbHash[index]&AfileBb != 0 {
			LegalKingMovesBbHash[0][index] &= ^HfileBb
			LegalKingMovesBbHash[1][index] &= ^HfileBb
		}
		if SingleIndexBbHash[index]&HfileBb != 0 {
			LegalKingMovesBbHash[0][index] &= ^AfileBb
			LegalKingMovesBbHash[1][index] &= ^AfileBb
		}
		if index == 4 {
			LegalKingMovesBbHash[0][4] |= SingleIndexBbHash[2] | SingleIndexBbHash[6]
		}
		if index == 60 {
			LegalKingMovesBbHash[1][60] |= SingleIndexBbHash[62] | SingleIndexBbHash[58]
		}

		// PAWN MOVE
		LegalPawnMovesBbHash[0][index] = 0
		LegalPawnMovesBbHash[1][index] = 0

		if index <= 56 {
			LegalPawnMovesBbHash[0][index] |= SingleIndexBbHash[index+7]
		}
		if index <= 55 {
			LegalPawnMovesBbHash[0][index] |= SingleIndexBbHash[index+8]
		}
		if index <= 54 {
			LegalPawnMovesBbHash[0][index] |= SingleIndexBbHash[index+9]
		}
		if index >= 7 {
			LegalPawnMovesBbHash[1][index] |= SingleIndexBbHash[index-7]
		}
		if index >= 8 {
			LegalPawnMovesBbHash[1][index] |= SingleIndexBbHash[index-8]
		}
		if index >= 9 {
			LegalPawnMovesBbHash[1][index] |= SingleIndexBbHash[index-9]
		}

		if index > 7 && index < 16 {
			LegalPawnMovesBbHash[0][index] |= SingleIndexBbHash[index+16]
		}
		if index > 47 && index < 56 {
			LegalPawnMovesBbHash[1][index] |= SingleIndexBbHash[index-16]
		}
		if SingleIndexBbHash[index]&AfileBb != 0 {
			LegalPawnMovesBbHash[0][index] &= ^HfileBb
			LegalPawnMovesBbHash[1][index] &= ^HfileBb

		}
		if SingleIndexBbHash[index]&HfileBb != 0 {
			LegalPawnMovesBbHash[0][index] &= ^AfileBb
			LegalPawnMovesBbHash[1][index] &= ^AfileBb
		}
	}
}

func printAllBitboards() {
	fmt.Println("\tA FILE")
	printBitBoard(AfileBb)
	fmt.Println("\tB FILE:")
	printBitBoard(BfileBb)
	fmt.Println("\tC FILE:")
	printBitBoard(CfileBb)
	fmt.Println("\tD FILE:")
	printBitBoard(DfileBb)
	fmt.Println("\tE FILE:")
	printBitBoard(EfileBb)
	fmt.Println("\tF FILE:")
	printBitBoard(FfileBb)
	fmt.Println("\tG FILE:")
	printBitBoard(GfileBb)
	fmt.Println("\tH FILE:")
	printBitBoard(HfileBb)
	for i := 0; i < 64; i++ {
		fmt.Println("\tBITBOARD LOOKUP:")
		printBitBoard(SingleIndexBbHash[i])
		fmt.Println("\tNORTH:")
		printBitBoard(NorthArrayBbHash[i])
		fmt.Println("\tSOUTH:")
		printBitBoard(SouthArrayBbHash[i])
		fmt.Println("\tEAST:")
		printBitBoard(EastArrayBbHash[i])
		fmt.Println("\tWEST:")
		printBitBoard(WestArrayBbHash[i])
		fmt.Println("\tNORTH EAST:")
		printBitBoard(NorthEastArrayBbHash[i])
		fmt.Println("\tNORTH WEST:")
		printBitBoard(NorthWestArrayBbHash[i])
		fmt.Println("\tSOUTH EAST:")
		printBitBoard(SouthEastArrayBbHash[i])
		fmt.Println("\tSOUTH WEST:")
		printBitBoard(SouthWestArrayBbHash[i])
		fmt.Println("\tKNIGHT ATTACK:")
		printBitBoard(KnightAttackBbHash[i])
		fmt.Println("\tKING MOVES:")
		printBitBoard(LegalKingMovesBbHash[0][i])
		printBitBoard(LegalKingMovesBbHash[1][i])
		fmt.Println("\tPAWN MOVES:")
		printBitBoard(LegalPawnMovesBbHash[0][i])
		printBitBoard(LegalPawnMovesBbHash[1][i])
		fmt.Println("\tEN PASSANT BB LOOKUP BY PAWN DESTINATION: ", i)
		printBitBoard(EnPassantBbHash[i])
		fmt.Println("\tATTACKED PAWN LOCATION BY EN PASSANT CAPTURE: ", i)
		printBitBoard(AttackedEnPassantPawnLocationBbHash[i])
	}
}

func printBitBoard(bb uint64) {
	bitboard, _ := bitboard.NewBitboard(bb)
	for i := 0; i < 64; i++ {
		fmt.Print(bitboard.GetBitValue(i))
		if ((i + 1) % 8) == 0 {
			fmt.Println("")
		}
	}
	fmt.Println("")
	fmt.Println("")
}
