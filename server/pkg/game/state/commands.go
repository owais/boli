package state

import (
	"encoding/json"
	"fmt"
)

type CommandName string

type Reducer interface {
	Reduce(State) (State, error)
}

const (
	JoinGame  CommandName = "game:join"
	LeaveGame CommandName = "game:leave"
	PlayCard  CommandName = "game:play-card"
	PlaceBid  CommandName = "game:place-bid"
)

type CommandValidator interface {
	Validate(state *State, cmd Command) error
}

type Command interface {
	Decode([]byte) error
	Validate(state *State) error
}

type Cmd struct {
	Name   CommandName     `json:"name"`
	Player string          `json:"player"`
	Data   json.RawMessage `json:"payload"`

	Validator CommandValidator
	// reducer Reducer
}

func (c *Cmd) Decode(data []byte) error {
	cmd := &Cmd{}
	if err := json.Unmarshal(data, cmd); err != nil {
		return fmt.Errorf("error decoding event: %w", err)
	}

	reducer, err := getReducerForEventType(cmd.Name)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(cmd.Data, reducer); err != nil {
		return fmt.Errorf("error decoding event data: %w", err)
	}
	// cmd.reducer = reducer
	return nil
}

func (c *Cmd) Validate(s *State) error {
	return nil
}

func getReducerForEventType(n CommandName) (Command, error) {
	switch n {
	case JoinGame:
		return &CommandJoinGame{}, nil
	case LeaveGame:
		return &CommandLeaveGame{}, nil
	case PlayCard:
		return &CommandPlayCard{}, nil
	case PlaceBid:
		return &CommandPlaceBid{}, nil
	default:
		return nil, fmt.Errorf("unknown command: %s", n)
	}
}

type CommandJoinGame struct {
	Cmd
}

type CommandLeaveGame struct {
	Cmd
}

type CommandPlayCard struct {
	Cmd
	Card string `json:"card"`
}

type CommandPlaceBid struct {
	Cmd
	Bid int `json:"bid"`

	// minimum bid allowed
	ValidateMin int
}

func (c *CommandPlaceBid) Validate(s *State) error {
	if c.Bid < 0 {
		return fmt.Errorf("bid must be non-negative")
	}
	if c.Bid < c.ValidateMin {
		return fmt.Errorf("bid must be greater than %d", c.ValidateMin)
	}
	return nil
}
