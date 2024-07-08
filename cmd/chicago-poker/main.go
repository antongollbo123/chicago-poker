package main

import (
	"github.com/antongollbo123/chicago-poker/internal/game"
	"github.com/antongollbo123/chicago-poker/internal/player"
)

func main() {
	p1 := player.NewPlayer("Anton")
	p2 := player.NewPlayer("Nora")
	gamet := game.NewGame([]*player.Player{p1, p2})

	gamet.StartGame()
}
