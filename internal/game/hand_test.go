package game

import (
	"testing"

	"github.com/antongollbo123/chicago-poker/pkg/cards"
)

func TestEvaluateHand(t *testing.T) {
	tests := []struct {
		name     string
		hand     []cards.Card
		expected HandEvaluation
	}{
		{
			name: "HighCard",
			hand: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Two},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Five},
				{Suit: cards.Hearts, Rank: cards.Seven},
				{Suit: cards.Hearts, Rank: cards.Nine},
			},
			expected: HandEvaluation{Rank: HighCard, Score: 0},
		},
		{
			name: "Pair",
			hand: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Two},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Hearts, Rank: cards.Eight},
				{Suit: cards.Hearts, Rank: cards.Four},
			},
			expected: HandEvaluation{Rank: Pair, Score: 1},
		},
		{
			name: "TwoPair",
			hand: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Two},
				{Suit: cards.Spades, Rank: cards.Two},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Hearts, Rank: cards.Three},
				{Suit: cards.Hearts, Rank: cards.Four},
			},
			expected: HandEvaluation{Rank: TwoPair, Score: 2},
		},
		{
			name: "Triple",
			hand: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Hearts, Rank: cards.Eight},
				{Suit: cards.Hearts, Rank: cards.Four},
			},
			expected: HandEvaluation{Rank: Triple, Score: 3},
		},
		{
			name: "Straight",
			hand: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Five},
				{Suit: cards.Spades, Rank: cards.Six},
				{Suit: cards.Hearts, Rank: cards.Seven},
				{Suit: cards.Spades, Rank: cards.Eight},
				{Suit: cards.Hearts, Rank: cards.Nine},
			},
			expected: HandEvaluation{Rank: Straight, Score: 4},
		},
		{
			name: "Flush",
			hand: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Two},
				{Suit: cards.Hearts, Rank: cards.Four},
				{Suit: cards.Hearts, Rank: cards.Six},
				{Suit: cards.Hearts, Rank: cards.Eight},
				{Suit: cards.Hearts, Rank: cards.Ten},
			},
			expected: HandEvaluation{Rank: Flush, Score: 5},
		},
		{
			name: "FullHouse",
			hand: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Hearts, Rank: cards.Four},
				{Suit: cards.Spades, Rank: cards.Four},
			},
			expected: HandEvaluation{Rank: FullHouse, Score: 6},
		},
		{
			name: "FourOfAKind",
			hand: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Four},
				{Suit: cards.Spades, Rank: cards.Four},
				{Suit: cards.Diamonds, Rank: cards.Four},
				{Suit: cards.Clubs, Rank: cards.Four},
				{Suit: cards.Hearts, Rank: cards.Two},
			},
			expected: HandEvaluation{Rank: FourOfAKind, Score: 7},
		},
		{
			name: "StraightFlush",
			hand: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Five},
				{Suit: cards.Hearts, Rank: cards.Six},
				{Suit: cards.Hearts, Rank: cards.Seven},
				{Suit: cards.Hearts, Rank: cards.Eight},
				{Suit: cards.Hearts, Rank: cards.Nine},
			},
			expected: HandEvaluation{Rank: StraightFlush, Score: 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EvaluateHand(tt.hand)
			if result != tt.expected {
				t.Errorf("EvaluateHand(%v) = %v, want %v", tt.hand, result, tt.expected)
			}
		})
	}
}

func TestIsStraight(t *testing.T) {
	tests := []struct {
		hand     []cards.Card
		expected bool
	}{
		{
			hand: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Five},
				{Suit: cards.Spades, Rank: cards.Six},
				{Suit: cards.Hearts, Rank: cards.Seven},
				{Suit: cards.Spades, Rank: cards.Eight},
				{Suit: cards.Hearts, Rank: cards.Nine},
			},
			expected: true,
		},
		{
			hand: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Five},
				{Suit: cards.Spades, Rank: cards.Six},
				{Suit: cards.Hearts, Rank: cards.Seven},
				{Suit: cards.Spades, Rank: cards.Eight},
				{Suit: cards.Hearts, Rank: cards.Ten},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := isStraight(tt.hand)
			if result != tt.expected {
				t.Errorf("isStraight(%v) = %v, want %v", tt.hand, result, tt.expected)
			}
		})
	}
}

func TestHasNOfAKind(t *testing.T) {
	rankCounts := map[cards.Rank]int{
		cards.Two:   1,
		cards.Three: 2,
		cards.Four:  1,
		cards.Eight: 1,
	}

	tests := []struct {
		n        int
		expected bool
	}{
		{n: 2, expected: true},
		{n: 3, expected: false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := hasNOfAKind(rankCounts, tt.n)
			if result != tt.expected {
				t.Errorf("hasNOfAKind(%v, %d) = %v, want %v", rankCounts, tt.n, result, tt.expected)
			}
		})
	}
}

func TestHasTwoPair(t *testing.T) {
	tests := []struct {
		rankCounts map[cards.Rank]int
		expected   bool
	}{
		{
			rankCounts: map[cards.Rank]int{
				cards.Two:   2,
				cards.Three: 2,
				cards.Four:  1,
			},
			expected: true,
		},
		{
			rankCounts: map[cards.Rank]int{
				cards.Two:   2,
				cards.Three: 1,
				cards.Four:  1,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := hasTwoPair(tt.rankCounts)
			if result != tt.expected {
				t.Errorf("hasTwoPair(%v) = %v, want %v", tt.rankCounts, result, tt.expected)
			}
		})
	}
}

func TestHasFullHouse(t *testing.T) {
	tests := []struct {
		rankCounts map[cards.Rank]int
		expected   bool
	}{
		{
			rankCounts: map[cards.Rank]int{
				cards.Two:   2,
				cards.Three: 3,
			},
			expected: true,
		},
		{
			rankCounts: map[cards.Rank]int{
				cards.Two:   1,
				cards.Three: 3,
				cards.Four:  1,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := hasFullHouse(tt.rankCounts)
			if result != tt.expected {
				t.Errorf("hasFullHouse(%v) = %v, want %v", tt.rankCounts, result, tt.expected)
			}
		})
	}
}
