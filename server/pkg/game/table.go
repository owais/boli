package game

import "fmt"

type Move struct {
	Card   *Card
	Player *Player
}

func (m Move) IsZero() bool {
	return m.Card == nil
}

type Table struct {
	moves []Move
	Suit  CardSuit
	Trump CardSuit
}

func (t *Table) Collect() []*Card {
	cards := []*Card{}
	for _, move := range t.moves {
		if move.Card != nil {
			cards = append(cards, move.Card)
		}
	}
	t.moves = []Move{}
	t.Suit = ""
	return cards
}

func (t *Table) Reset() {
	t.moves = []Move{}
	t.Suit = ""
	t.Trump = ""
}

func (t *Table) PrintCards() {
	for _, move := range t.moves {
		if move.Card != nil {
			fmt.Print(move.Card.String())
		}
	}
	fmt.Println()
}

func (t *Table) Add(move Move) {
	t.moves = append(t.moves, move)
	if t.Suit == "" {
		t.Suit = move.Card.Suit
	}
}

func (t *Table) maxOf(suit CardSuit) Move {
	var max Move
	for _, move := range t.moves {
		if move.Card.Suit == suit {
			if max.IsZero() || move.Card.Value > max.Card.Value {
				max = move
			}
		}
	}
	return max
}

func (t *Table) Max() Move {
	max := t.maxOf(t.Suit)
	if t.Trump != "" {
		maxTrump := t.maxOf(t.Trump)
		if !maxTrump.IsZero() {
			max = maxTrump
		}
	}
	return max
}

func (t *Table) Winner() *Player {
	max := t.Max()
	if max.Player == nil {
		fmt.Println("max card not found, printing all moves ")
		for _, m := range t.moves {
			fmt.Println("card: ", m.Card)
			fmt.Println("player: ", m.Player)
			fmt.Println("==========")
		}
	}
	return max.Player
}
