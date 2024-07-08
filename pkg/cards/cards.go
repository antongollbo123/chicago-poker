package cards

type Suit string
type Rank string

const (
	Spades   Suit = "♠"
	Hearts   Suit = "♥"
	Diamonds Suit = "♦"
	Clubs    Suit = "♣"
)

const (
	Two   Rank = "2"
	Three Rank = "3"
	Four  Rank = "4"
	Five  Rank = "5"
	Six   Rank = "6"
	Seven Rank = "7"
	Eight Rank = "8"
	Nine  Rank = "9"
	Ten   Rank = "10"
	Jack  Rank = "Jack"
	Queen Rank = "Queen"
	King  Rank = "King"
	Ace   Rank = "Ace"
)

type Card struct {
	Suit Suit
	Rank Rank
}

// NewCard creates a new card with the given suit and rank
func NewCard(suit Suit, rank Rank) Card {
	return Card{Suit: suit, Rank: rank}
}
