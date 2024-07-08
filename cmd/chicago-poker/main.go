package main

import (
	"fmt"

	"github.com/antongollbo123/chicago-poker/internal/game"
	"github.com/antongollbo123/chicago-poker/internal/player"
)

func main() {
	gamet := game.NewGame()
	p1 := player.NewPlayer("Anton")
	p2 := player.NewPlayer("Nora")

	gamet.Players = append(gamet.Players, p1, p2)
	for _, player := range gamet.Players {
		for i := 0; i < 5; i++ {
			card, ok := gamet.Deck.Draw()
			if !ok {
				fmt.Println("No more cards")
				break
			}

			player.Hand = append(player.Hand, card)
		}
		fmt.Println(player.Name, player.Score, player.Hand)
		fmt.Println("Player 1 score: ", game.EvaluateHand(p1.Hand))
		fmt.Println("Player 2 score: ", game.EvaluateHand(p2.Hand))
	}
}
