package game

import (
	"fmt"
	"sync"

	"github.com/owais/boli/server/pkg/comms"
	"github.com/owais/boli/server/pkg/game/core"
	"github.com/owais/boli/server/pkg/game/state"
)

type Game struct {
	mu *sync.Mutex

	// TODO: separate state and game. store all data in state
	// state   state.State
	deck    *core.Deck
	players core.Players
	dealer  *core.Player
	score   int

	// communication channel with the players
	comms comms.Comms
}

func New() *Game {
	team1 := &core.Team{Name: "Team 1", Allot: scoreAdd}
	team2 := &core.Team{Name: "Team 2", Allot: scoreSubtract}

	g := &Game{
		deck: core.NewDeck(),
		players: core.Players{
			core.NewPlayer("1", "Player 1", team1),
			core.NewPlayer("2", "Player 2", team2),
			core.NewPlayer("3", "Player 3", team1),
			core.NewPlayer("4", "Player 4", team2),
			core.NewPlayer("5", "Player 5", team1),
			core.NewPlayer("6", "Player 6", team2),
		},
	}

	g.comms = comms.NewCli(g.players)

	return g
}

func scoreAdd(old, score int) int {
	return old + score
}

func scoreSubtract(old, score int) int {
	return old - score
}

func (g *Game) Start() {
	g.dealer = g.Toss()

	for {
		g.PrintState()
		setScore := g.PlaySet()
		g.score += setScore

		// TODO: implement light suffle and do it here

		// if threshold is hit, end the game
		if g.score >= 50 {
			g.Win(g.oppositeTeam(g.dealer), g.score)
			return
		}

		// TODO: if score was zero and now went to positive, make the next player dealer

		// if score goes below 0, make the next player the dealer
		if g.score <= 0 {
			g.score = 0 - g.score
			g.dealer = g.players.Next(g.dealer)
		}
	}
}

func (g *Game) collectToDeck(stack []*core.Card) {
	for _, card := range stack {
		g.deck.Put(card)
	}

	for _, p := range g.players {
		g.deck.Put(p.Hand.Collect()...)
	}
}

func (g *Game) Win(team *core.Team, score int) {
	// g.winner = team
	g.score = score
	fmt.Println("Game Over! Winner: ", team.Name)
}

func (g *Game) PlaySet() int {
	// move bid to every round
	bid := g.DealAndBid()
	if bid.IsZero() {
		// reduce the score by 1 and end the round
		// g.score -= 1
		return -1
	}

	stack := []*core.Card{}

	// play rounds the bid either fails or succeeds (max of 8 rounds)
	currentPlayer := bid.Player
	won := 0
	loss := 0
	lossThreshold := 8 - bid.Rounds + 1

	score := 0
	for i := 0; i < 8; i++ {

		table := &core.Table{}

		g.PlayRound(table, currentPlayer, bid.Player)
		winner := table.Winner()
		stack = append(stack, table.Collect()...)

		// winner starts the next round
		currentPlayer = winner

		fmt.Println("Round winner: ", winner)

		if winner.Team == bid.Player.Team {
			won++
		} else {
			loss++
		}

		g.PrintTable(table, won, loss, &bid)
		if won >= bid.Rounds {
			score = bid.Rounds
			fmt.Println("Bidder won the Set")
			break
		} else if loss >= lossThreshold {
			score = bid.Rounds * 2
			fmt.Println("Bidder lost the Set")
			break
		}

	}

	g.collectToDeck(stack)
	return score

}

// PlayRound plays a round of the game and returns the winner
func (g *Game) PlayRound(table *core.Table, player *core.Player, bidder *core.Player) *core.Player {
	for i := 0; i < len(g.players); i++ {
		// If suit has been set, i.e, first card has been played
		if table.Suit != "" {
			// .. and if trump has not been set and the player does not have a card of the same suit
			if table.Trump == "" && !player.HasSuit(table.Suit) {
				// .. then ask the bidder to select the trump suit
				table.Trump = g.WaitForTrump(bidder)
			}
		}

		move := g.WaitForMove(player, table.Suit)
		table.Add(move)

		player = g.players.Next(player)
	}
	return table.Winner()
}

func (g *Game) WaitForTrump(player *core.Player) core.Suit {
	answer := getUserTextInput(player, "Select a trump suit", []string{"hearts", "diamonds", "clubs", "spades"})

	switch answer {
	case "hearts":
		return core.SuitHearts
	case "diamonds":
		return core.SuitDiamonds
	case "clubs":
		return core.SuitClubs
	case "spades":
		return core.SuitSpades
	default:
		fmt.Println("Invalid choice, please try again.")
		return g.WaitForTrump(player)
	}
}

// WaitForMove waits for a player to play a card and returns the move
// If the player does not have a card of the same suit, they can play any card
// If the player has a card of the same suit, they must play it.
// If the player has a trump card, they can play it only if they don't have a card of the same suit.
func (g *Game) WaitForMove(player *core.Player, suit core.Suit) core.Move {
	choices := []string{}

	for _, card := range player.Hand.All() {
		if card.Suit == suit {
			choices = append(choices, card.Id)
		}
	}
	if len(choices) == 0 {
		for _, card := range player.Hand.All() {
			choices = append(choices, card.Id)
		}
	}

	choice := getUserTextInput(player, "Play a card", choices)

	// remove the card from player's hand and return it in a Move
	return core.Move{Player: player, Card: player.Hand.Remove(choice)}
}

func (g *Game) DealAndBid() core.Bid {

	// deal five cards to each player and wait for any player to bid at least 6
	g.Deal(5)
	g.PrintHands()
	bid := g.WaitForBid(6, nil)

	// deal three more cards to each player and wait for any player to bid at least 6
	g.Deal(3)
	g.PrintHands()

	bid2 := g.WaitForBid(7, nil)

	if bid2.Rounds > bid.Rounds {
		bid = bid2
	}

	if bid.IsZero() {
		// wait for the dealer to bid at least 5
		bid = g.WaitForBid(5, g.dealer)
	}

	return bid
}

func (g *Game) WaitForBid(min int, player *core.Player) core.Bid {
	// ask each player to bid
	bids := []core.Bid{}
	for _, p := range g.players {
		bid := state.CommandPlaceBid{ValidateMin: min}

		err := g.comms.Receive(p, &bid)
		if err != nil {
			// pass validation error back to user
			fmt.Println("Error receiving bid: ", err)
			continue
		}
		// num := getUserMinNumberInputOrPass(&p, "Enter your bid (0 to pass)", min)
		if bid.Bid == 0 {
			// player passed
			continue
		}
		bids = append(bids, core.Bid{Player: &p, Rounds: bid.Bid})
	}

	var maxBid core.Bid
	for _, b := range bids {
		if b.Rounds > maxBid.Rounds {
			maxBid = b
		}
	}
	return maxBid
}

func (g *Game) PrintState() {
	fmt.Println()
	fmt.Println("===========================")
	fmt.Println("Overall Score: ", g.score)
	fmt.Printf("Dealer: %s (%s)\n", g.dealer, g.dealer.Team)
	fmt.Println("===========================")
	fmt.Println()
}

func (g *Game) PrintTable(table *core.Table, wins, losses int, bid *core.Bid) {
	fmt.Println()
	fmt.Println("===========================")
	fmt.Printf("Bid: %d by %s (%s)\n", bid.Rounds, bid.Player, bid.Player.Team)
	fmt.Println("Round Score: ", fmt.Sprintf("%d:%d", wins, losses))
	fmt.Println("Trump: ", table.Trump)
	g.PrintHands()
	fmt.Println("===========================")
	fmt.Println()
}

func (g *Game) PrintHands() {
	for _, player := range g.players {
		player.PrintHand()
	}
}

// Toss deals the cards to each player until someone is dealt a Jack of any suite
func (g *Game) Toss() *core.Player {
	// TODO: separate this into steps so players get to see the toss in UI
	drawn := []*core.Card{}
	for {
		for _, player := range g.players {
			card := g.deck.DrawOne()
			drawn = append(drawn, card)
			if card.Value == core.CardValueJack {
				for _, card := range drawn {
					g.deck.Put(card)
				}
				g.deck.Shuffle()
				return &player
			}
		}
	}
}

// Deal deals cards to each player starting from the player after the dealer in counter-clockwise direction.
// Each player is dealt the number of cards specified by eachPlayer.
func (g *Game) Deal(eachPlayer int) {
	for _, player := range g.players {
		cards := g.deck.DrawN(eachPlayer)
		player.Hand.Add(cards...)
	}
}

func (g *Game) oppositeTeam(player *core.Player) *core.Team {
	return g.players.Next(player).Team
}
