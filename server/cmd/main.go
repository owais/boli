package main

import (
	"github.com/owais/boli/server/pkg/game"
)

func main() {
	game := game.New()
	game.Start()
}
