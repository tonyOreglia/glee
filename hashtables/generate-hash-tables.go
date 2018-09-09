// Package hashtables calculates various bitboard lookup tables
// which provide data about specific locations on the chess board.
package hashtables

import (
	"os"
	"strconv"

	"github.com/tonyoreglia/glee/bitboard"
)

// HashTables holds bitoard lookup tables used in move generation
type HashTables struct {
	AfileBb                             uint64
	BfileBb                             uint64
	CfileBb                             uint64
	DfileBb                             uint64
	EfileBb                             uint64
	FfileBb                             uint64
	GfileBb                             uint64
	HfileBb                             uint64
	FourthRankBb                        uint64
	FifthRankBb                         uint64
	SingleIndexBbHash                   [64]uint64
	EnPassantBbHash                     [64]uint64
	AttackedEnPassantPawnLocationBbHash [64]uint64
	NorthArrayBbHash                    [64]uint64
	SouthArrayBbHash                    [64]uint64
	EastArrayBbHash                     [64]uint64
	WestArrayBbHash                     [64]uint64
	NorthEastArrayBbHash                [64]uint64
	NorthWestArrayBbHash                [64]uint64
	SouthEastArrayBbHash                [64]uint64
	SouthWestArrayBbHash                [64]uint64
	KnightAttackBbHash                  [64]uint64
	LegalKingMovesBbHash                [2][64]uint64
	LegalPawnMovesBbHash                [2][64]uint64
}

func CalculateAllLookupBbs() *HashTables {
	hashTables := new(HashTables)
	hashTables.AfileBb = 0x101010101010101
	hashTables.BfileBb = hashTables.AfileBb << 1
	hashTables.CfileBb = hashTables.AfileBb << 2
	hashTables.DfileBb = hashTables.AfileBb << 3
	hashTables.EfileBb = hashTables.AfileBb << 4
	hashTables.FfileBb = hashTables.AfileBb << 5
	hashTables.GfileBb = hashTables.AfileBb << 6
	hashTables.HfileBb = hashTables.AfileBb << 7
	hashTables.FourthRankBb = 0xFF00000000
	hashTables.FifthRankBb = 0xFF000000

	for index := 0; index < 64; index++ {
		hashTables.EnPassantBbHash[index] = uint64(0)
		hashTables.AttackedEnPassantPawnLocationBbHash[index] = uint64(0)
		hashTables.NorthArrayBbHash[index] = uint64(0)
		hashTables.SouthArrayBbHash[index] = uint64(0)
		hashTables.EastArrayBbHash[index] = uint64(0)
		hashTables.WestArrayBbHash[index] = uint64(0)
		hashTables.NorthEastArrayBbHash[index] = uint64(0)
		hashTables.NorthWestArrayBbHash[index] = uint64(0)
		hashTables.SouthEastArrayBbHash[index] = uint64(0)
		hashTables.SouthWestArrayBbHash[index] = uint64(0)
		hashTables.KnightAttackBbHash[index] = uint64(0)
	}
	generateSingleBitLookup(hashTables)
	generateArrayBitboardLookup(hashTables)
	generateEnPassantBitboardLookup(hashTables)
	return hashTables
}

func generateSingleBitLookup(ht *HashTables) {
	for i := uint(0); i < 64; i++ {
		ht.SingleIndexBbHash[i] = uint64(1) << i
	}
}

func generateEnPassantBitboardLookup(ht *HashTables) {
	for i := 0; i < 64; i++ {
		ht.EnPassantBbHash[i] = 0
		ht.AttackedEnPassantPawnLocationBbHash[i] = 0
		switch {
		case i > 23 && i < 32:
			ht.EnPassantBbHash[i] = ht.SingleIndexBbHash[i-8]
		case i > 31 && i < 40:
			ht.EnPassantBbHash[i] = ht.SingleIndexBbHash[i+8]
		case i > 15 && i < 24:
			ht.AttackedEnPassantPawnLocationBbHash[i] = ht.SingleIndexBbHash[i+8]
		case i > 39 && i < 48:
			ht.AttackedEnPassantPawnLocationBbHash[i] = ht.SingleIndexBbHash[i-8]
		}
	}
}

func generateArrayBitboardLookup(ht *HashTables) {
	for index := 0; index < 64; index++ {
		northOfIndex := index
		for northOfIndex > 7 {
			ht.NorthArrayBbHash[index] |= ht.SingleIndexBbHash[northOfIndex-8]
			northOfIndex -= 8
		}
		southOfIndex := index
		for southOfIndex < 56 {
			ht.SouthArrayBbHash[index] |= ht.SingleIndexBbHash[southOfIndex+8]
			southOfIndex += 8
		}
		eastOfIndex := index
		for (ht.SingleIndexBbHash[eastOfIndex])&ht.HfileBb == 0 {
			ht.EastArrayBbHash[index] |= ht.SingleIndexBbHash[eastOfIndex+1]
			eastOfIndex++
		}
		westOfIndex := index
		for (ht.SingleIndexBbHash[westOfIndex])&ht.AfileBb == 0 {
			ht.WestArrayBbHash[index] |= ht.SingleIndexBbHash[westOfIndex-1]
			westOfIndex--
		}
		northEastOfIndex := index
		for (ht.SingleIndexBbHash[northEastOfIndex])&ht.HfileBb == 0 && northEastOfIndex > 7 {
			ht.NorthEastArrayBbHash[index] |= ht.SingleIndexBbHash[northEastOfIndex-7]
			northEastOfIndex -= 7
		}
		northWestOfIndex := index
		for (ht.SingleIndexBbHash[northWestOfIndex])&ht.AfileBb == 0 && northWestOfIndex > 8 {
			ht.NorthWestArrayBbHash[index] |= ht.SingleIndexBbHash[northWestOfIndex-9]
			northWestOfIndex -= 9
		}
		southEastOfIndex := index
		for (ht.SingleIndexBbHash[southEastOfIndex])&ht.HfileBb == 0 && southEastOfIndex < 56 {
			ht.SouthEastArrayBbHash[index] |= ht.SingleIndexBbHash[southEastOfIndex+9]
			southEastOfIndex += 9
		}
		southWestOfIndex := index
		for (ht.SingleIndexBbHash[southWestOfIndex])&ht.AfileBb == 0 && southWestOfIndex < 56 {
			ht.SouthWestArrayBbHash[index] |= ht.SingleIndexBbHash[southWestOfIndex+7]
			southWestOfIndex += 7
		}
		if (index+10) < 64 && ht.SingleIndexBbHash[index]&(ht.HfileBb|ht.GfileBb) == 0 {
			ht.KnightAttackBbHash[index] |= ht.SingleIndexBbHash[index+10]
		}
		if (index+6) < 64 && ht.SingleIndexBbHash[index]&(ht.AfileBb|ht.BfileBb) == 0 {
			ht.KnightAttackBbHash[index] |= ht.SingleIndexBbHash[index+6]
		}
		if (index+17) < 64 && ht.SingleIndexBbHash[index]&ht.HfileBb == 0 {
			ht.KnightAttackBbHash[index] |= ht.SingleIndexBbHash[index+17]
		}
		if (index+15) < 64 && ht.SingleIndexBbHash[index]&ht.AfileBb == 0 {
			ht.KnightAttackBbHash[index] |= ht.SingleIndexBbHash[index+15]
		}
		if (index-10) >= 0 && ht.SingleIndexBbHash[index]&(ht.AfileBb|ht.BfileBb) == 0 {
			ht.KnightAttackBbHash[index] |= ht.SingleIndexBbHash[index-10]
		}
		if (index-6) >= 0 && ht.SingleIndexBbHash[index]&(ht.HfileBb|ht.GfileBb) == 0 {
			ht.KnightAttackBbHash[index] |= ht.SingleIndexBbHash[index-6]
		}
		if (index-17) >= 0 && ht.SingleIndexBbHash[index]&ht.AfileBb == 0 {
			ht.KnightAttackBbHash[index] |= ht.SingleIndexBbHash[index-17]
		}
		if (index-15) >= 0 && ht.SingleIndexBbHash[index]&ht.HfileBb == 0 {
			ht.KnightAttackBbHash[index] |= ht.SingleIndexBbHash[index-15]
		}

		ht.LegalKingMovesBbHash[0][index] = 0
		if index != 63 {
			ht.LegalKingMovesBbHash[0][index] |= ht.SingleIndexBbHash[index+1]
		}
		if index != 0 {
			ht.LegalKingMovesBbHash[0][index] |= ht.SingleIndexBbHash[index-1]
		}
		if index <= 55 {
			ht.LegalKingMovesBbHash[0][index] |= ht.SingleIndexBbHash[index+8]
		}
		if index >= 8 {
			ht.LegalKingMovesBbHash[0][index] |= ht.SingleIndexBbHash[index-8]
		}
		if index <= 54 {
			ht.LegalKingMovesBbHash[0][index] |= ht.SingleIndexBbHash[index+9]
		}
		if index >= 9 {
			ht.LegalKingMovesBbHash[0][index] |= ht.SingleIndexBbHash[index-9]
		}
		if index <= 56 {
			ht.LegalKingMovesBbHash[0][index] |= ht.SingleIndexBbHash[index+7]
		}
		if index >= 7 {
			ht.LegalKingMovesBbHash[0][index] |= ht.SingleIndexBbHash[index-7]
		}
		ht.LegalKingMovesBbHash[1][index] = ht.LegalKingMovesBbHash[0][index]
		if ht.SingleIndexBbHash[index]&ht.AfileBb != 0 {
			ht.LegalKingMovesBbHash[0][index] &= ^ht.HfileBb
			ht.LegalKingMovesBbHash[1][index] &= ^ht.HfileBb
		}
		if ht.SingleIndexBbHash[index]&ht.HfileBb != 0 {
			ht.LegalKingMovesBbHash[0][index] &= ^ht.AfileBb
			ht.LegalKingMovesBbHash[1][index] &= ^ht.AfileBb
		}
		if index == 4 {
			ht.LegalKingMovesBbHash[0][4] |= ht.SingleIndexBbHash[2] | ht.SingleIndexBbHash[6]
		}
		if index == 60 {
			ht.LegalKingMovesBbHash[1][60] |= ht.SingleIndexBbHash[62] | ht.SingleIndexBbHash[58]
		}
		ht.LegalPawnMovesBbHash[0][index] = 0
		ht.LegalPawnMovesBbHash[1][index] = 0
		if index <= 56 {
			ht.LegalPawnMovesBbHash[0][index] |= ht.SingleIndexBbHash[index+7]
		}
		if index <= 55 {
			ht.LegalPawnMovesBbHash[0][index] |= ht.SingleIndexBbHash[index+8]
		}
		if index <= 54 {
			ht.LegalPawnMovesBbHash[0][index] |= ht.SingleIndexBbHash[index+9]
		}
		if index >= 7 {
			ht.LegalPawnMovesBbHash[1][index] |= ht.SingleIndexBbHash[index-7]
		}
		if index >= 8 {
			ht.LegalPawnMovesBbHash[1][index] |= ht.SingleIndexBbHash[index-8]
		}
		if index >= 9 {
			ht.LegalPawnMovesBbHash[1][index] |= ht.SingleIndexBbHash[index-9]
		}

		if index > 7 && index < 16 {
			ht.LegalPawnMovesBbHash[0][index] |= ht.SingleIndexBbHash[index+16]
		}
		if index > 47 && index < 56 {
			ht.LegalPawnMovesBbHash[1][index] |= ht.SingleIndexBbHash[index-16]
		}
		if ht.SingleIndexBbHash[index]&ht.AfileBb != 0 {
			ht.LegalPawnMovesBbHash[0][index] &= ^ht.HfileBb
			ht.LegalPawnMovesBbHash[1][index] &= ^ht.HfileBb

		}
		if ht.SingleIndexBbHash[index]&ht.HfileBb != 0 {
			ht.LegalPawnMovesBbHash[0][index] &= ^ht.AfileBb
			ht.LegalPawnMovesBbHash[1][index] &= ^ht.AfileBb
		}
	}
}

// PrintAllBitboards prints all hash lookups to a local file
func PrintAllBitboards(ht *HashTables) {
	f, _ := os.Create("hash-tables.txt")
	defer f.Close()
	f.Write([]byte("\tA FILE"))
	printBitBoard(f, ht.AfileBb)
	f.Write([]byte("\tB FILE:"))
	printBitBoard(f, ht.BfileBb)
	f.Write([]byte("\tC FILE:"))
	printBitBoard(f, ht.CfileBb)
	f.Write([]byte("\tD FILE:"))
	printBitBoard(f, ht.DfileBb)
	f.Write([]byte("\tE FILE:"))
	printBitBoard(f, ht.EfileBb)
	f.Write([]byte("\tF FILE:"))
	printBitBoard(f, ht.FfileBb)
	f.Write([]byte("\tG FILE:"))
	printBitBoard(f, ht.GfileBb)
	f.Write([]byte("\tH FILE:"))
	printBitBoard(f, ht.HfileBb)
	for i := 0; i < 64; i++ {
		f.Write([]byte("\tBITBOARD LOOKUP:"))
		printBitBoard(f, ht.SingleIndexBbHash[i])
		f.Write([]byte("\tNORTH:"))
		printBitBoard(f, ht.NorthArrayBbHash[i])
		f.Write([]byte("\tSOUTH:"))
		printBitBoard(f, ht.SouthArrayBbHash[i])
		f.Write([]byte("\tEAST:"))
		printBitBoard(f, ht.EastArrayBbHash[i])
		f.Write([]byte("\tWEST:"))
		printBitBoard(f, ht.WestArrayBbHash[i])
		f.Write([]byte("\tNORTH EAST:"))
		printBitBoard(f, ht.NorthEastArrayBbHash[i])
		f.Write([]byte("\tNORTH WEST:"))
		printBitBoard(f, ht.NorthWestArrayBbHash[i])
		f.Write([]byte("\tSOUTH EAST:"))
		printBitBoard(f, ht.SouthEastArrayBbHash[i])
		f.Write([]byte("\tSOUTH WEST:"))
		printBitBoard(f, ht.SouthWestArrayBbHash[i])
		f.Write([]byte("\tKNIGHT ATTACK:"))
		printBitBoard(f, ht.KnightAttackBbHash[i])
		f.Write([]byte("\tKING MOVES:"))
		printBitBoard(f, ht.LegalKingMovesBbHash[0][i])
		printBitBoard(f, ht.LegalKingMovesBbHash[1][i])
		f.Write([]byte("\tPAWN MOVES:"))
		printBitBoard(f, ht.LegalPawnMovesBbHash[0][i])
		printBitBoard(f, ht.LegalPawnMovesBbHash[1][i])
		f.Write([]byte("\tEN PASSANT BB LOOKUP BY PAWN DESTINATION: "))
		printBitBoard(f, ht.EnPassantBbHash[i])
		f.Write([]byte("\tATTACKED PAWN LOCATION BY EN PASSANT CAPTURE: "))
		printBitBoard(f, ht.AttackedEnPassantPawnLocationBbHash[i])
	}
}

func printBitBoard(f *os.File, bb uint64) {
	bitboard, _ := bitboard.NewBitboard(bb)
	nl := []byte("\n")
	f.Write(nl)
	for i := 0; i < 64; i++ {
		b := []byte(strconv.Itoa(int(bitboard.GetBitValue(i))))
		f.Write(b)
		if ((i + 1) % 8) == 0 {
			f.Write(nl)
		}
	}
	f.Write(nl)
	f.Write(nl)
}
