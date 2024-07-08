package game

import (
	"github.com/antongollbo123/chicago-poker/internal/deck"
	"github.com/antongollbo123/chicago-poker/internal/player"
)

type Stage string
type Chicago bool

const (
	Poker Stage = "Poker"
	Trick Stage = "Trick"
)

type Game struct {
	deck    *deck.Deck
	players []*player.Player
	round   int
	stage   Stage
}

func NewGame() (*Game, *deck.Deck) {
	game := Game{}
	game.round = 1
	game.stage = Poker

	deck := deck.NewDeck()

	return &game, deck
}
