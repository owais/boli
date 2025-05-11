package state

import "github.com/owais/boli/server/pkg/game/core"

type State struct {
	Players core.Players
	Teams   []core.Team
	Deck    core.Deck
	Dealer  core.Player
	Score   int
}

func (s *State) Apply(event Event) error {
	return nil
	/*
		reducer, err := event.Decode()
		if err != nil {
			return *s, err
		}

		newState, err := event.reducer.Reduce(*s)
		if err != nil {
			return *s, err
		}
		return newState, nil
	*/
}
