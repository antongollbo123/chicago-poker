package player

import "github.com/antongollbo123/chicago-poker/pkg/cards"

type Player struct {
	Name  string
	Hand  []cards.Card // TODO: Change to max 5 slots in array
	Score int
}

func NewPlayer(name string) *Player {
	return &Player{Name: name, Hand: []cards.Card{}, Score: 0}
}
