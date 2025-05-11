package server

import (
	"encoding/json"
)

type command struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}
