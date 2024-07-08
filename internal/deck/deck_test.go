package deck

import (
	"fmt"
	"testing"

	"github.com/antongollbo123/chicago-poker/pkg/cards"
)

func TestNewDeck(t *testing.T) {
	deck := NewDeck()

	if deck.NumCards != 52 {
		t.Errorf("expected 52 cards, got %d", deck.NumCards)
	}
	fmt.Println(deck.NumCards, deck.cards)
	cardSet := make(map[cards.Card]bool)
	for _, card := range deck.cards {
		if cardSet[card] {
			t.Errorf("duplicate card found: %v", card)
		}
		cardSet[card] = true
	}

	if len(cardSet) != 52 {
		t.Errorf("expected 52 unique cards, got %d", len(cardSet))
	}

}

func TestShuffle(t *testing.T) {
	deck := NewDeck()
	initialOrder := make([]cards.Card, 52)
	copy(initialOrder, deck.cards)

	deck.Shuffle()

	sameOrder := true
	for i, card := range initialOrder {
		if deck.cards[i] != card {
			sameOrder = false
			break
		}
	}

	if sameOrder {
		t.Errorf("expected shuffled deck to have a different order than the initial deck")
	}
}

func TestDraw(t *testing.T) {
	deck := NewDeck()
	initialNumCards := deck.NumCards

	card, ok := deck.Draw()
	if !ok {
		t.Errorf("expected to draw a card, but draw failed")
	}

	if deck.NumCards != initialNumCards-1 {
		t.Errorf("expected %d cards, got %d", initialNumCards-1, deck.NumCards)
	}

	if card == (cards.Card{}) {
		t.Errorf("expected a valid card, got an empty card")
	}
}

func TestDrawMultiple(t *testing.T) {
	deck := NewDeck()
	initialNumCards := deck.NumCards

	numToDraw := 5
	drawnCards := deck.DrawMultiple(numToDraw)

	if len(drawnCards) != numToDraw {
		t.Errorf("expected to draw %d cards, but got %d", numToDraw, len(drawnCards))
	}

	if deck.NumCards != initialNumCards-numToDraw {
		t.Errorf("expected %d cards left in the deck, got %d", initialNumCards-numToDraw, deck.NumCards)
	}
}

func TestDrawMultipleNotEnoughCards(t *testing.T) {
	deck := NewDeck()
	numToDraw := 60 // More than the number of cards in the deck

	drawnCards := deck.DrawMultiple(numToDraw)

	if len(drawnCards) != 52 {
		t.Errorf("expected to draw 52 cards, but got %d", len(drawnCards))
	}

	if deck.NumCards != 0 {
		t.Errorf("expected 0 cards left in the deck, got %d", deck.NumCards)
	}
}
