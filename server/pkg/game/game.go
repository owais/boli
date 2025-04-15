package game

import "fmt"

type Game struct {
	deck    *Deck
	players Players
	dealer  *Player
	score   int

	history []Event
}

func New() *Game {
	team1 := &Team{Name: "Team 1", Allot: scoreAdd}
	team2 := &Team{Name: "Team 2", Allot: scoreSubtract}

	return &Game{
		deck: NewDeck(),
		players: Players{
			{ID: "1", Name: "Player 1", Team: team1},
			{ID: "2", Name: "Player 2", Team: team2},
			{ID: "3", Name: "Player 3", Team: team1},
			{ID: "4", Name: "Player 4", Team: team2},
			{ID: "5", Name: "Player 5", Team: team1},
			{ID: "6", Name: "Player 6", Team: team2},
		},
	}
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
			// g.dealer = g.dealer.Team.oppositePlayer(g.dealer)
			// g.score = 0
		}
	}
}

func (g *Game) collectToDeck(stack []*Card) {
	for _, card := range stack {
		g.deck.Put(card)
	}

	for _, p := range g.players {
		g.deck.Put(p.Hand.Collect()...)
	}
}

func (g *Game) Win(team *Team, score int) {
	// g.winner = team
	// g.history = append(g.history, Event{Type: "WinGame", Player: nil, Card: nil})
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

	stack := []*Card{}

	// play rounds the bid either fails or succeeds (max of 8 rounds)
	currentPlayer := bid.Player
	won := 0
	loss := 0
	lossThreshold := 8 - bid.Rounds + 1

	score := 0
	for i := 0; i < 8; i++ {

		table := &Table{}

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
func (g *Game) PlayRound(table *Table, player *Player, bidder *Player) *Player {
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

func (g *Game) WaitForTrump(player *Player) CardSuit {
	answer := getUserTextInput(player, "Select a trump suit", []string{"hearts", "diamonds", "clubs", "spades"})

	switch answer {
	case "hearts":
		return CardSuitHearts
	case "diamonds":
		return CardSuitDiamonds
	case "clubs":
		return CardSuitClubs
	case "spades":
		return CardSuitSpades
	default:
		fmt.Println("Invalid choice, please try again.")
		return g.WaitForTrump(player)
	}
}

// WaitForMove waits for a player to play a card and returns the move
// If the player does not have a card of the same suit, they can play any card
// If the player has a card of the same suit, they must play it.
// If the player has a trump card, they can play it only if they don't have a card of the same suit.
func (g *Game) WaitForMove(player *Player, suit CardSuit) Move {
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
	return Move{Player: player, Card: player.Hand.Remove(choice)}
}

func (g *Game) DealAndBid() Bid {

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

func (g *Game) WaitForBid(min int, player *Player) Bid {
	// ask each player to bid
	bids := []Bid{}
	for _, p := range g.players {
		num := getUserMinNumberInputOrPass(p, "Enter your bid (0 to pass)", min)
		if num == 0 {
			// player passed
			continue
		}
		bids = append(bids, Bid{Player: p, Rounds: num})
	}

	var maxBid Bid
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

func (g *Game) PrintTable(table *Table, wins, losses int, bid *Bid) {
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
func (g *Game) Toss() *Player {
	drawn := []*Card{}

	for {
		for _, player := range g.players {
			card := g.deck.DrawOne()
			drawn = append(drawn, card)
			// g.history = append(g.history, Event{Type: "Toss", Player: player, Card: card})
			if card.Value == CardValueJack {
				for _, card := range drawn {
					g.deck.Put(card)
				}
				g.deck.Shuffle()
				return player
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
		// g.history = append(g.history, Event{Type: "Deal", Player: player, Card: card})
	}
}

func (g *Game) oppositeTeam(player *Player) *Team {
	return g.players.Next(player).Team
}
