package game

type Bid struct {
	Player *Player
	Rounds int
}

func (b Bid) IsZero() bool {
	return b.Rounds == 0
}

type Round struct {
	Winner *Team
	Bid    *Bid
}
