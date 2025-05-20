package game

type PlayerIO interface {
	InputRequest(p *Player, ev EventType) Event
	Output(p *Player, ev Event)
}
