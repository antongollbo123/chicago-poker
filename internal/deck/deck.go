package deck

import (
	"math/rand"
	"time"

	"github.com/antongollbo123/chicago-poker/pkg/cards"
)

type Deck struct {
	cards    []cards.Card
	NumCards int
}

func NewDeck() *Deck {
	suits := []cards.Suit{cards.Hearts, cards.Spades, cards.Clubs, cards.Diamonds}
	ranks := []cards.Rank{cards.Two, cards.Three, cards.Four, cards.Five, cards.Six, cards.Seven, cards.Eight, cards.Nine, cards.Ten, cards.Jack, cards.Queen, cards.King, cards.Ace}

	var deck []cards.Card

	cardCount := 0

	for _, suit := range suits {
		for _, rank := range ranks {
			deck = append(deck, cards.NewCard(suit, rank))
			cardCount++
		}
	}
	return &Deck{cards: deck, NumCards: cardCount}
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
	d.cards = d.cards[1:] // TODO: Add rest deck?
	d.NumCards--
	return card, true
}

func (d *Deck) DrawMultiple(numCards int) []cards.Card {
	drawnCards := make([]cards.Card, 0, numCards)
	for i := 0; i < numCards; i++ {
		card, ok := d.Draw()
		if !ok {
			break // Handle case where there are not enough cards left in the deck
		}
		drawnCards = append(drawnCards, card)
	}
	return drawnCards
}
