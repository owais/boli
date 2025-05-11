package state

import (
	"encoding/json"

	"github.com/owais/boli/server/pkg/game/core"
)

type Event struct {
	Type      EventType       `json:"name"`
	EventData json.RawMessage `json:"eventData"`
}

func (e *Event) Apply() error {
	return nil
}

type EventGameJoined struct {
	Player core.Player `json:"player"`
}

func (e *EventGameJoined) Reduce(s State) (State, error) {
	if err := s.Players.Add(e.Player); err != nil {
		return s, err
	}
	return s, nil
}

type EventGameLeft struct {
}

func (e *EventGameLeft) Reduce(s State) (State, error) {
	return s, nil
}

type EventCardDrawn struct {
	Player core.Player `json:"player"`
	Card   core.Card   `json:"card"`
}

func (e *EventCardDrawn) Reduce(s State) (State, error) {
	card := s.Deck.DrawOne()
	e.Player.Hand.Add(card)
	return s, nil
}

type EventToss struct {
	Dealer core.Player `json:"dealer"`
}

func (e *EventToss) Reduce(s State) (State, error) {
	s.Dealer = e.Dealer
	return s, nil
}

type EventCardDraw struct {
}

func (e *EventCardDraw) Reduce(s State) (State, error) {
	return s, nil
}

type EventCardDealt struct {
}

func (e *EventCardDealt) Reduce(s State) (State, error) {
	return s, nil
}

type EventCardPlayed struct {
}

func (e *EventCardPlayed) Reduce(s State) (State, error) {
	return s, nil
}

type EventRoundWon struct {
}

func (e *EventRoundWon) Reduce(s State) (State, error) {
	return s, nil
}

type EventGameWon struct {
}

func (e *EventGameWon) Reduce(s State) (State, error) {
	return s, nil
}
