package comms

import (
	"github.com/owais/boli/server/pkg/game/core"
	"github.com/owais/boli/server/pkg/game/state"
)

type Comms interface {
	Send(p core.Player, event state.Event) error
	Receive(p core.Player, cmd state.Command) error
}
