package utils

import (
	"fmt"

	"github.com/antongollbo123/chicago-poker/internal/game"
)

func PrintPlayerNameAndHand(game *game.Game) {
	for _, player := range game.Players {
		fmt.Println("Player Name: ", player.Name, "Player Hand: ", player.Hand)
	}

}
