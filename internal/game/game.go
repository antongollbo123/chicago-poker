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
	Deck    *deck.Deck
	Players []*player.Player
	Round   int
	Stage   Stage
}

func NewGame() *Game {
	game := Game{}
	game.Round = 1
	game.Stage = Poker

	deck := deck.NewDeck()
	deck.Shuffle()
	game.Deck = deck
	return &game
}

func (g *Game) CheckHand(p *player.Player) int {

	return 1
}
