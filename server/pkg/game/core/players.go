package core

import (
	"fmt"
)

type Hand struct {
	cards []*Card
}

func (h *Hand) All() []*Card {
	return h.cards
}

func (h *Hand) Add(cards ...*Card) {
	h.cards = append(h.cards, cards...)
}

func (h *Hand) Remove(id string) *Card {
	for i, card := range h.cards {
		if card.Id == id {
			h.cards = append(h.cards[:i], h.cards[i+1:]...)
			return card
		}
	}
	return nil
}

func (h *Hand) Collect() []*Card {
	cards := h.cards
	h.cards = []*Card{}
	return cards
}

type Player struct {
	ID   string
	Name string
	Team *Team
	Hand *Hand
}

func NewPlayer(id, name string, team *Team) Player {
	return Player{
		ID:   id,
		Name: name,
		Team: team,
		Hand: &Hand{},
	}
}

func (p Player) String() string {
	return p.Name
}

func (p Player) PrintHand() {
	fmt.Print(p.Name + "'s hand: ")
	for _, card := range p.Hand.All() {
		fmt.Print(card.String(), " ")
	}
	fmt.Println()
}

func (p Player) HasSuit(suit Suit) bool {
	for _, card := range p.Hand.All() {
		if card.Suit == suit {
			return true
		}
	}
	return false
}

type Team struct {
	Name  string
	Cards []*Card
	Allot func(old, score int) int
}

func (t Team) String() string {
	return t.Name
}

type Players []Player

func (p *Players) Add(player Player) error {
	if len(*p) >= 6 {
		return fmt.Errorf("cannot add more than 6 players")
	}
	*p = append(*p, player)
	return nil
}

/*
func (p Players) Add(player Player) (Players, error) {
	if len(p) >= 6 {
		return p, fmt.Errorf("cannot add more than 6 players")
	}
	p = append(p, player)
	return p, nil
}
*/

func (p *Players) Next(current *Player) *Player {
	var next *Player
	for i, player := range *p {
		if player.ID == current.ID {
			// if current is the last player, return the first player, otherwise current + 1
			if i+1 < len(*p) {
				next = &(*p)[i+1]
			} else {
				next = &(*p)[0]
			}
			break
		}
	}
	return next
}
