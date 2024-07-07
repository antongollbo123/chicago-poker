package deck

import (
	"math/rand"
	"time"

	"github.com/antongollbo123/chicago-poker/pkg/cards"
)

type Deck struct {
	cards []cards.Card
}

func NewDeck() *Deck {
	suits := []cards.Suit{cards.Hearts, cards.Spades, cards.Clubs, cards.Diamonds}
	ranks := []cards.Rank{cards.Two, cards.Three, cards.Four, cards.Five, cards.Six, cards.Seven, cards.Eight, cards.Nine, cards.Ten, cards.Jack, cards.Queen, cards.King, cards.Ace}

	var deck []cards.Card

	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, cards.NewCard(suit, rank))
		}
	}
	return &Deck{cards: deck}
}

func (d *Deck) Shuffle() {
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Draw() (cards.Card, bool) {
	if len(d.cards) == 0 {
		return cards.Card{}, false
	}
	card := d.cards[0]
	d.cards = d.cards[1:]
	return card, true
}
