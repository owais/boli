package game

import (
	"fmt"
	mathrand "math/rand"
	"strings"
	"time"
)

var rand = mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

type CardSuit string

const (
	CardSuitHearts   CardSuit = "H"
	CardSuitDiamonds CardSuit = "D"
	CardSuitClubs    CardSuit = "C"
	CardSuitSpades   CardSuit = "S"
)

var CardSuitLabels = map[CardSuit]string{
	CardSuitHearts:   "♥️",
	CardSuitDiamonds: "♦️",
	CardSuitClubs:    "♣️",
	CardSuitSpades:   "♠️",
}

type CardValue int

const (
	CardValueThree CardValue = 3
	CardValueFour  CardValue = 4
	CardValueFive  CardValue = 5
	CardValueSix   CardValue = 6
	CardValueSeven CardValue = 7
	CardValueEight CardValue = 8
	CardValueNine  CardValue = 9
	CardValueTen   CardValue = 10
	CardValueJack  CardValue = 11
	CardValueQueen CardValue = 12
	CardValueKing  CardValue = 13
	CardValueAce   CardValue = 14
)

var cardValueLabels = map[CardValue]string{
	CardValueThree: "3",
	CardValueFour:  "4",
	CardValueFive:  "5",
	CardValueSix:   "6",
	CardValueSeven: "7",
	CardValueEight: "8",
	CardValueNine:  "9",
	CardValueTen:   "10",
	CardValueJack:  "J",
	CardValueQueen: "Q",
	CardValueKing:  "K",
	CardValueAce:   "A",
}

type Card struct {
	Id    string
	Suit  CardSuit
	Value CardValue
	label string
}

func (c Card) String() string {
	return string(c.label) + CardSuitLabels[c.Suit]
}

type Deck struct {
	cards []*Card
}

func NewDeck() *Deck {
	cards := []*Card{}
	for _, suit := range []CardSuit{CardSuitHearts, CardSuitDiamonds, CardSuitClubs, CardSuitSpades} {
		for value := CardValueThree; value <= CardValueAce; value++ {
			cards = append(cards, &Card{
				Id:   strings.ToLower(fmt.Sprintf("%d%s", value, suit)),
				Suit: suit, Value: value, label: cardValueLabels[value]},
			)
		}
	}
	return &Deck{cards: cards}
}

func (d *Deck) Put(cards ...*Card) {
	d.cards = append(d.cards, cards...)
}

func (d *Deck) Shuffle() {
	for i := range d.cards {
		j := rand.Intn(i + 1)
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	}
}

// Deck prints all cards in the deck printing 4 cards per line
func (d *Deck) PrintCards() {
	for _, card := range d.cards {
		fmt.Print(card.String() + " ")
	}
	fmt.Println()
}

func (d *Deck) Size() int {
	return len(d.cards)
}

func (d *Deck) DrawOne() *Card {
	if len(d.cards) == 0 {
		return nil
	}
	card := d.cards[len(d.cards)-1]
	d.cards = d.cards[:len(d.cards)-1]
	return card
}

func (d *Deck) DrawN(count int) []*Card {
	if count <= 0 || len(d.cards) == 0 {
		return nil
	}

	if count > len(d.cards) {
		count = len(d.cards)
	}

	drawnCards := make([]*Card, count)
	for i := 0; i < count; i++ {
		card := d.cards[len(d.cards)-1]
		d.cards = d.cards[:len(d.cards)-1]
		drawnCards[i] = card
	}

	return drawnCards
}
