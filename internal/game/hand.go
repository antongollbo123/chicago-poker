package game

import (
	"sort"

	"github.com/antongollbo123/chicago-poker/pkg/cards"
)

type HandRank int

const (
	HighCard HandRank = iota
	Pair
	TwoPair
	Triple
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
)

type HandEvaluation struct {
	Rank  HandRank
	Score int
}

func EvaluateHand(hand []cards.Card) HandEvaluation {
	rankCounts := make(map[cards.Rank]int)
	suitCounts := make(map[cards.Suit]int)

	for _, card := range hand {
		rankCounts[card.Rank]++
		suitCounts[card.Suit]++
	}

	isFlush := len(suitCounts) == 1
	isStraight := isStraight(hand)
	switch {
	case isStraight && isFlush:
		return HandEvaluation{Rank: StraightFlush, Score: 8}
	case hasNOfAKind(rankCounts, 4):
		return HandEvaluation{Rank: FourOfAKind, Score: 7}
	case hasFullHouse(rankCounts):
		return HandEvaluation{Rank: FullHouse, Score: 6}
	case isFlush:
		return HandEvaluation{Rank: Flush, Score: 5}
	case isStraight:
		return HandEvaluation{Rank: Straight, Score: 4}
	case hasNOfAKind(rankCounts, 3):
		return HandEvaluation{Rank: Triple, Score: 3}
	case hasTwoPair(rankCounts):
		return HandEvaluation{Rank: TwoPair, Score: 2}
	case hasNOfAKind(rankCounts, 2):
		return HandEvaluation{Rank: Pair, Score: 1}
	default:
		return HandEvaluation{Rank: HighCard, Score: 0}
	}
}

func isStraight(hand []cards.Card) bool {
	if len(hand) < 5 {
		return false
	}
	ranks := []int{}
	rankMap := map[cards.Rank]int{
		cards.Two: 2, cards.Three: 3, cards.Four: 4, cards.Five: 5, cards.Six: 6, cards.Seven: 7,
		cards.Eight: 8, cards.Nine: 9, cards.Ten: 10, cards.Jack: 11, cards.Queen: 12, cards.King: 13, cards.Ace: 14,
	}
	for _, card := range hand {
		ranks = append(ranks, rankMap[card.Rank])
	}
	sort.Ints(ranks)
	for i := 0; i < len(ranks)-1; i++ {
		if ranks[i+1] != ranks[i]+1 {
			return false
		}
	}
	return true
}

func hasNOfAKind(rankCounts map[cards.Rank]int, n int) bool {
	for _, count := range rankCounts {
		if count == n {
			return true
		}
	}
	return false
}

func hasTwoPair(rankCounts map[cards.Rank]int) bool {
	pairCount := 0
	for _, count := range rankCounts {
		if count == 2 {
			pairCount++
		}
	}
	return pairCount == 2
}

func hasFullHouse(rankCounts map[cards.Rank]int) bool {
	hasThree := false
	hasTwo := false
	for _, count := range rankCounts {
		if count == 3 {
			hasThree = true
		}
		if count == 2 {
			hasTwo = true
		}
	}
	return hasThree && hasTwo
}
