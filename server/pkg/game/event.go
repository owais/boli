package game

import "time"

type EventType string

const (
	EventTypeToss    EventType = "Toss"
	EventTypeBid     EventType = "Bid"
	EventTypeMaxBid  EventType = "MaxWon"
	EventTypeDraw    EventType = "Draw"
	EventTypeDeal    EventType = "Deal"
	EventTypePlay    EventType = "Play"
	EventTypeSetEnd  EventType = "SetEnd"
	EventTypeGameEnd EventType = "GameEnd"
)

type Event interface {
	Type() EventType
}

type event struct {
	EventType EventType `json:"type"`
	Ts        time.Time `json:"ts"`
}

func (e event) Type() EventType {
	return e.EventType
}

type CardDealtEvent struct {
	event
	Player Player `json:"player"`
	Cards  []Card `json:"card"`
}

func NewCardDealtEvent(player Player, cards Cards) CardDealtEvent {
	return CardDealtEvent{
		event: event{
			EventType: EventTypeDeal,
			Ts:        time.Now(),
		},
		Player: player,
		Cards:  cards.Values(),
	}
}

type BidEvent struct {
	event
	Player Player `json:"player"`
	Bid    int    `json:"rounds"`
}

func NewBidEvent(bid Bid) BidEvent {
	return BidEvent{
		event: event{
			EventType: EventTypeBid,
			Ts:        time.Now(),
		},
		Player: *bid.Player,
		Bid:    bid.Rounds,
	}
}

type MaxBidEvent struct {
	event
	Player Player `json:"player"`
	Bid    int    `json:"rounds"`
}

func NewMaxBidEvent(bid Bid) MaxBidEvent {
	return MaxBidEvent{
		event: event{
			EventType: EventTypeMaxBid,
			Ts:        time.Now(),
		},
		Player: *bid.Player,
		Bid:    bid.Rounds,
	}
}

type SetEndEvent struct {
	event
	Winner Team `json:"team"`
	Score  int  `json:"score"`
}

func NewSetEndEvent(winner Team, score int) SetEndEvent {
	return SetEndEvent{
		event: event{
			EventType: EventTypeSetEnd,
			Ts:        time.Now(),
		},
		Winner: winner,
		Score:  score,
	}
}

type GameEndEvent struct {
	event
	Winner Team `json:"team"`
	Score  int  `json:"score"`
}

func NewGameEndEvent(winner Team, score int) GameEndEvent {
	return GameEndEvent{
		event: event{
			EventType: EventTypeGameEnd,
			Ts:        time.Now(),
		},
		Winner: winner,
		Score:  score,
	}
}

type TrumpSelectedEvent struct {
	event
	Player Player   `json:"player"`
	Suit   CardSuit `json:"suit"`
}

func NewTrumpSelectedEvent(player Player, suit CardSuit) TrumpSelectedEvent {
	return TrumpSelectedEvent{
		event: event{
			EventType: EventTypePlay,
			Ts:        time.Now(),
		},
		Player: player,
		Suit:   suit,
	}
}

type RoundEndEvent struct {
	event
	Winner Player `json:"player"`
}

func NewRoundEndEvent(winner Player) RoundEndEvent {
	return RoundEndEvent{
		event: event{
			EventType: EventTypeDraw,
			Ts:        time.Now(),
		},
		Winner: winner,
	}
}
