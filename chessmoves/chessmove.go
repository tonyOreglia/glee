package chessmoves

const MoveListSize = 100

type Move struct {
	origin      uint8
	destination uint8
	promotion   uint64
}

func NewMovesList() []Move {
	return make([]Move, MoveListSize)
}
