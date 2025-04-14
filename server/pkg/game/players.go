package game

import "fmt"

type Player struct {
	ID   string
	Name string
	Team *Team
	Hand []*Card
}

func (p *Player) String() string {
	return p.Name
}

func (p *Player) PrintHand() {
	fmt.Print(p.Name + "'s hand: ")
	for _, card := range p.Hand {
		fmt.Print(card.String(), " ")
	}
	fmt.Println()
}

type Team struct {
	Name  string
	Cards []*Card
}

func (t *Team) String() string {
	return t.Name
}
