package game

import "fmt"

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
	Hand Hand
}

func (p *Player) String() string {
	return p.Name
}

func (p *Player) PrintHand() {
	fmt.Print(p.Name + "'s hand: ")
	for _, card := range p.Hand.All() {
		fmt.Print(card.String(), " ")
	}
	fmt.Println()
}

func (p *Player) HasSuit(suit CardSuit) bool {
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

func (t *Team) String() string {
	return t.Name
}

func scoreAdd(old, score int) int {
	return old + score
}

func scoreSubtract(old, score int) int {
	return old - score
}

type Players []*Player

func (p Players) Next(current *Player) *Player {
	for i, player := range p {
		if player == current {
			if i+1 < len(p) {
				return p[i+1]
			}
			return p[0]
		}
	}
	return nil
}
