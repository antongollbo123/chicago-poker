package cards

import "fmt"

type Suit string
type Rank int

const (
	Hearts   Suit = "♥"
	Spades   Suit = "♠"
	Diamonds Suit = "♦"
	Clubs    Suit = "♣"
)

const (
	Two   Rank = 2
	Three Rank = 3
	Four  Rank = 4
	Five  Rank = 5
	Six   Rank = 6
	Seven Rank = 7
	Eight Rank = 8
	Nine  Rank = 9
	Ten   Rank = 10
	Jack  Rank = 11
	Queen Rank = 12
	King  Rank = 13
	Ace   Rank = 14
)

type Card struct {
	Suit Suit
	Rank Rank
}

// NewCard creates a new card with the given suit and rank
func NewCard(suit Suit, rank Rank) Card {
	return Card{Suit: suit, Rank: rank}
}

func RankToString(r Rank) string {
	switch r {
	case 11:
		return "Jack"
	case 12:
		return "Queen"
	case 13:
		return "King"
	case 14:
		return "Ace"
	default:
		return fmt.Sprint(r)
	}
}

func SuitValue(suit string) int {
	switch suit {
	case "♥":
		return 4
	case "♠":
		return 3
	case "♦":
		return 2
	case "♣":
		return 1
	default:
		return 0
	}
}
