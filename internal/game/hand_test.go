package game

import (
	"reflect"
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
			if result.Rank != tt.expected.Rank && result.Score != tt.expected.Score {
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
			_, result := isStraight(tt.hand)
			if result != tt.expected {
				t.Errorf("isStraight(%v) = %v, want %v", tt.hand, result, tt.expected)
			}
		})
	}
}

func TestHasNOfAKind(t *testing.T) {
	rankCounts := map[cards.Rank][]cards.Card{
		cards.Two: {
			{Suit: cards.Hearts, Rank: cards.Two},
		},
		cards.Three: {
			{Suit: cards.Diamonds, Rank: cards.Three},
			{Suit: cards.Spades, Rank: cards.Three},
		},
		cards.Four: {
			{Suit: cards.Diamonds, Rank: cards.Four},
		},
		cards.Eight: {
			{Suit: cards.Diamonds, Rank: cards.Eight},
		},
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
			_, result := getNOfAKind(rankCounts, tt.n)
			if result != tt.expected {
				t.Errorf("hasNOfAKind(%v, %d) = %v, want %v", rankCounts, tt.n, result, tt.expected)
			}
		})
	}
}

func TestHasTwoPair(t *testing.T) {
	tests := []struct {
		rankCounts map[cards.Rank][]cards.Card
		expected   bool
	}{
		{
			rankCounts: map[cards.Rank][]cards.Card{
				cards.Two: {
					{Suit: cards.Hearts, Rank: cards.Two},
					{Suit: cards.Spades, Rank: cards.Two},
				},
				cards.Three: {
					{Suit: cards.Diamonds, Rank: cards.Three},
					{Suit: cards.Spades, Rank: cards.Three},
				},
			},
			expected: true,
		},
		{
			rankCounts: map[cards.Rank][]cards.Card{
				cards.Two: {
					{Suit: cards.Hearts, Rank: cards.Two},
				},
				cards.Three: {
					{Suit: cards.Hearts, Rank: cards.Three},
					{Suit: cards.Spades, Rank: cards.Three},
				},
				cards.Four: {
					{Suit: cards.Spades, Rank: cards.Four},
					{Suit: cards.Diamonds, Rank: cards.Four},
				},
			},
			expected: true,
		},
		{
			rankCounts: map[cards.Rank][]cards.Card{
				cards.Two: {
					{Suit: cards.Hearts, Rank: cards.Two},
				},
				cards.Three: {
					{Suit: cards.Hearts, Rank: cards.Three},
				},
				cards.Four: {
					{Suit: cards.Spades, Rank: cards.Four},
					{Suit: cards.Diamonds, Rank: cards.Four},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			_, result := getTwoPair(tt.rankCounts)
			if result != tt.expected {
				t.Errorf("hasTwoPair(%v) = %v, want %v", tt.rankCounts, result, tt.expected)
			}
		})
	}
}

func TestHasFullHouse(t *testing.T) {
	tests := []struct {
		rankCounts map[cards.Rank][]cards.Card
		expected   bool
	}{
		{
			rankCounts: map[cards.Rank][]cards.Card{
				cards.Two: {
					{Suit: cards.Hearts, Rank: cards.Two},
					{Suit: cards.Spades, Rank: cards.Two},
				},
				cards.Three: {
					{Suit: cards.Diamonds, Rank: cards.Three},
					{Suit: cards.Spades, Rank: cards.Three},
					{Suit: cards.Hearts, Rank: cards.Three},
				},
			},
			expected: true,
		},
		{
			rankCounts: map[cards.Rank][]cards.Card{
				cards.Two: {
					{Suit: cards.Hearts, Rank: cards.Two},
				},
				cards.Three: {
					{Suit: cards.Hearts, Rank: cards.Three},
					{Suit: cards.Spades, Rank: cards.Three},
					{Suit: cards.Diamonds, Rank: cards.Three},
				},
				cards.Four: {
					{Suit: cards.Spades, Rank: cards.Four},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			_, result := getFullHouse(tt.rankCounts)
			if result != tt.expected {
				t.Errorf("hasFullHouse(%v) = %v, want %v", tt.rankCounts, result, tt.expected)
			}
		})
	}
}

func TestEvaluteTwoHands(t *testing.T) {
	tests := []struct {
		hand1    []cards.Card
		hand2    []cards.Card
		expected []cards.Card
	}{
		{

			hand1: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Two},
				{Suit: cards.Spades, Rank: cards.Two},
				{Suit: cards.Diamonds, Rank: cards.Four},
				{Suit: cards.Spades, Rank: cards.Five},
				{Suit: cards.Hearts, Rank: cards.Six},
			},
			hand2: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Two},
				{Suit: cards.Spades, Rank: cards.Two},
				{Suit: cards.Diamonds, Rank: cards.Four},
				{Suit: cards.Spades, Rank: cards.Five},
				{Suit: cards.Hearts, Rank: cards.Seven},
			},
			expected: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Two},
				{Suit: cards.Spades, Rank: cards.Two},
				{Suit: cards.Diamonds, Rank: cards.Four},
				{Suit: cards.Spades, Rank: cards.Five},
				{Suit: cards.Hearts, Rank: cards.Seven},
			},
		},
		{
			hand1: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Two},
				{Suit: cards.Spades, Rank: cards.Two},
				{Suit: cards.Diamonds, Rank: cards.Four},
				{Suit: cards.Spades, Rank: cards.Five},
				{Suit: cards.Hearts, Rank: cards.Eight},
			},
			hand2: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Diamonds, Rank: cards.Four},
				{Suit: cards.Spades, Rank: cards.Five},
				{Suit: cards.Hearts, Rank: cards.Seven},
			},
			expected: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Diamonds, Rank: cards.Four},
				{Suit: cards.Spades, Rank: cards.Five},
				{Suit: cards.Hearts, Rank: cards.Seven},
			},
		},
		{
			hand1: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Diamonds, Rank: cards.Four},
				{Suit: cards.Spades, Rank: cards.Five},
				{Suit: cards.Hearts, Rank: cards.Eight},
			},
			hand2: []cards.Card{
				{Suit: cards.Diamonds, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Diamonds, Rank: cards.Four},
				{Suit: cards.Spades, Rank: cards.Five},
				{Suit: cards.Hearts, Rank: cards.Eight},
			},
			expected: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Diamonds, Rank: cards.Four},
				{Suit: cards.Spades, Rank: cards.Five},
				{Suit: cards.Hearts, Rank: cards.Eight},
			},
		},
		{
			hand1: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Diamonds, Rank: cards.Five},
				{Suit: cards.Clubs, Rank: cards.Four},
				{Suit: cards.Hearts, Rank: cards.Four},
			},
			hand2: []cards.Card{
				{Suit: cards.Diamonds, Rank: cards.Four},
				{Suit: cards.Spades, Rank: cards.Four},
				{Suit: cards.Clubs, Rank: cards.Five},
				{Suit: cards.Clubs, Rank: cards.Three},
				{Suit: cards.Diamonds, Rank: cards.Three},
			},
			expected: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Three},
				{Suit: cards.Spades, Rank: cards.Three},
				{Suit: cards.Diamonds, Rank: cards.Five},
				{Suit: cards.Clubs, Rank: cards.Four},
				{Suit: cards.Hearts, Rank: cards.Four},
			},
		},
		{
			hand1: []cards.Card{
				{Suit: cards.Hearts, Rank: cards.Six},
				{Suit: cards.Spades, Rank: cards.Six},
				{Suit: cards.Diamonds, Rank: cards.Five},
				{Suit: cards.Clubs, Rank: cards.Three},
				{Suit: cards.Hearts, Rank: cards.Four},
			},
			hand2: []cards.Card{
				{Suit: cards.Diamonds, Rank: cards.Ten},
				{Suit: cards.Spades, Rank: cards.Ten},
				{Suit: cards.Clubs, Rank: cards.Five},
				{Suit: cards.Clubs, Rank: cards.Three},
				{Suit: cards.Diamonds, Rank: cards.Four},
			},
			expected: []cards.Card{
				{Suit: cards.Diamonds, Rank: cards.Ten},
				{Suit: cards.Spades, Rank: cards.Ten},
				{Suit: cards.Clubs, Rank: cards.Five},
				{Suit: cards.Clubs, Rank: cards.Three},
				{Suit: cards.Diamonds, Rank: cards.Four},
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			winningHand, winningHandEvaluation := EvaluateTwoHands(tt.hand1, tt.hand2)
			expectedSorted := sortCards(sortSuit(tt.expected))
			if !reflect.DeepEqual(winningHand, tt.expected) {
				t.Errorf("EvaluateTwoHands(%v) = %v, want %v", winningHand, winningHandEvaluation, expectedSorted)
			}
		})
	}
}
