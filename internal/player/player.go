package player

import "github.com/antongollbo123/chicago-poker/pkg/cards"

type Player struct {
	Name  string
	Hand  []cards.Card
	Score int
}

func newPlayer(name string) *Player {
	return &Player{Name: name, Hand: []cards.Card{}, Score: 0}
}
