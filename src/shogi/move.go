package shogi

type Move struct {
	drop   bool
	figPos int
	pos    int
}

func MoveFromString(desc []rune) *Move {
	created := new(Move)
	if desc[0] == '-' {
		created.drop = true
		created.figPos = pieceStrToInt(rune(desc[1]))
	} else {
		created.drop = false
		created.figPos = posStrToInt(desc[0:2])
	}
	created.pos = posStrToInt(desc[2:4])
	return created
}

func (move *Move) String() string {
	out := ""
	if move.drop {
		out += "-"
		out += string(pieceIntToStr(move.figPos))
		out += posIntToStr(move.pos)
	} else {
		out += posIntToStr(move.figPos)
		out += posIntToStr(move.pos)
	}
	return out
}
