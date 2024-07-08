package game

import (
	"sort"

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

func (g *Game) TossCards(playerIndex int, indicesToRemove []int) {
	if playerIndex < 0 || playerIndex >= len(g.Players) {
		// Handle invalid player index
		return
	}

	// Sort indicesToRemove in descending order to safely remove cards from hand slice
	sort.Sort(sort.Reverse(sort.IntSlice(indicesToRemove)))

	// Remove cards from player's hand based on indicesToRemove
	for _, idx := range indicesToRemove {
		if idx >= 0 && idx < len(g.Players[playerIndex].Hand) {
			g.Players[playerIndex].Hand = append(g.Players[playerIndex].Hand[:idx], g.Players[playerIndex].Hand[idx+1:]...)
		}
	}

	// Deal new cards from the deck
	newCards := g.Deck.DrawMultiple(len(indicesToRemove))
	g.Players[playerIndex].Hand = append(g.Players[playerIndex].Hand, newCards...)
}
