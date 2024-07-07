package main

import (
	"fmt"

	"github.com/antongollbo123/chicago-poker/internal/deck"
)

func main() {

	d := deck.NewDeck()
	fmt.Println(d)
	d.Shuffle()
	fmt.Println("shuffled:", d)

	for i := 0; i < 5; i++ {

		card, ok := d.Draw()

		if !ok {
			fmt.Println("No more cards")
			break
		}
		fmt.Printf("%s of %s\n", card.Rank, card.Suit)
	}
}
