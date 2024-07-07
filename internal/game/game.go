package game

import "github.com/antongollbo123/chicago-poker/internal/deck"

type Stage string
type Chicago bool

const (
	Poker Stage = "Poker"
	Trick Stage = "Trick"
)

type Game struct {
	round int
	stage Stage
}

func NewGame() (*Game, *deck.Deck) {
	game := Game{}
	game.round = 1
	game.stage = Poker

	deck := deck.NewDeck()

	return &game, deck
}
