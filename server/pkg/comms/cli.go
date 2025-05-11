package comms

import (
	"fmt"

	"github.com/owais/boli/server/pkg/game/core"
	"github.com/owais/boli/server/pkg/game/state"
)

type Cli struct {
	players core.Players
}

func NewCli(players core.Players) *Cli {
	return &Cli{
		players: players,
	}
}

func (c *Cli) Send(p core.Player, event state.Event) error {
	fmt.Println(event)
	return nil
}

func (c *Cli) Receive(p core.Player, cmd state.Command) error {
	println("Enter message:")
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return err
	}

	return nil

	// return Message{Type: "cli", Content: input}, nil
	// return &state.Command{}, nil
}
