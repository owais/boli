package cli

import (
	"github.com/owais/boli/server/pkg/game"
)

func Run() {
	io := &CliPlayerIO{}
	game := game.New(io)
	game.Join("Aaqib")
	game.Join("Tariq")
	game.Join("Mudassir")
	game.Join("Irfan")
	game.Join("Irshad")
	game.Join("Owais")
	game.Start()
}

type CliPlayerIO struct {
}

func (c *CliPlayerIO) InputRequest(p *game.Player, ev game.EventType) game.Event {
	return nil
}

func (c *CliPlayerIO) Output(p *game.Player, ev game.Event) {
}
