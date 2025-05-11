package state

type EventType string

const (
	GameJoined EventType = "event:joined"
	GameLeft   EventType = "event:left"
	RoundWon   EventType = "event:round-won"
	GameWon    EventType = "event:game-won"
	Toss       EventType = "event:toss"
	CardDrawn  EventType = "event:card-drawn"
	CardsDealt EventType = "event:dealt"
	CardPlayed EventType = "event:played"
)
