package game

type EventType string

const (
	EventTypeToss     EventType = "Toss"
	EventTypeDraw     EventType = "Draw"
	EventTypeDeal     EventType = "Deal"
	EventTypePlay     EventType = "Play"
	EventTypeWinRound EventType = "WinRound"
	EventTypeWinGame  EventType = "WinGame"
	EventTypeEndGame  EventType = "EndGame"
)

type Event struct {
	Type   string
	Player *Player
	Card   *Card
}
