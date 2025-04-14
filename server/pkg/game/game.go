package game

import "fmt"

type Game struct {
	Deck   *Deck
	team1  *Team
	team2  *Team
	p1     *Player
	p2     *Player
	p3     *Player
	p4     *Player
	p5     *Player
	p6     *Player
	dealer *Player
	caller *Player

	score int

	history []Event
}

func New() *Game {
	team1 := &Team{Name: "Team 1"}
	team2 := &Team{Name: "Team 2"}
	deck := NewDeck()

	return &Game{
		Deck:  deck,
		team1: team1,
		team2: team2,
		p1:    &Player{ID: "1", Name: "Player 1", Team: team1},
		p2:    &Player{ID: "2", Name: "Player 2", Team: team2},
		p3:    &Player{ID: "3", Name: "Player 3", Team: team1},
		p4:    &Player{ID: "4", Name: "Player 4", Team: team2},
		p5:    &Player{ID: "5", Name: "Player 5", Team: team1},
		p6:    &Player{ID: "6", Name: "Player 6", Team: team2},
	}
}

func (g *Game) Start() {
	g.Deck.Shuffle()
	g.PrintState()
	g.Toss()
	g.Deck.Shuffle()
	g.PrintState()
	g.Deal(5)
	g.PrintState()
	g.Deal(3)
	g.PrintState()
}

func (g *Game) PrintState() {
	fmt.Println("Game State:")
	fmt.Println("Dealer: ", g.dealer)
	fmt.Println("Caller: ", g.caller)
	fmt.Print("Deck: ")
	g.Deck.PrintCards()
	for _, player := range []*Player{g.p1, g.p2, g.p3, g.p4, g.p5, g.p6} {
		player.PrintHand()
	}
	fmt.Println("==========================")
}

// Toss deals the cards to each player until someone is dealt a Jack of any suite
func (g *Game) Toss() {
	drawn := []*Card{}

	for {
		for _, player := range []*Player{g.p1, g.p2, g.p3, g.p4, g.p5, g.p6} {
			card := g.Deck.DrawOne()
			drawn = append(drawn, card)
			if card == nil {
				for _, card := range drawn {
					g.Deck.Put(card)
				}
				return
			}
			g.history = append(g.history, Event{Type: "Toss", Player: player, Card: card})
			if card.Value == CardValueJack {
				g.dealer = player
			}
		}
	}

	// g.Deck.Shuffle()
}

// Deal deals cards to each player starting from the player after the dealer in counter-clockwise direction.
// Each player is dealt the number of cards specified by eachPlayer.
func (g *Game) Deal(eachPlayer int) {
	for _, player := range []*Player{g.p1, g.p2, g.p3, g.p4, g.p5, g.p6} {
		cards := g.Deck.DrawN(eachPlayer)
		player.Hand = append(player.Hand, cards...)
		// g.history = append(g.history, Event{Type: "Deal", Player: player, Card: card})
	}
}
