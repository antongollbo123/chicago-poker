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

func (hr HandRank) String() string {
	switch hr {
	case HighCard:
		return "High Card"
	case Pair:
		return "Pair"
	case TwoPair:
		return "Two Pair"
	case Triple:
		return "Triple"
	case Straight:
		return "Straight"
	case Flush:
		return "Flush"
	case FullHouse:
		return "Full House"
	case FourOfAKind:
		return "Four of a Kind"
	case StraightFlush:
		return "Straight Flush"
	default:
		return "Unknown Hand Rank"
	}
}

type HandEvaluation struct {
	Rank       HandRank
	Score      int
	ScoreCards []cards.Card
}

func EvaluateHand(hand []cards.Card) HandEvaluation {
	rankCounts := make(map[cards.Rank]int)
	suitCounts := make(map[cards.Suit]int)

	for _, card := range hand {
		rankCounts[card.Rank]++
		suitCounts[card.Suit]++
	}

	isFlush := len(suitCounts) == 1
	straightCards, isStraight := isStraight(hand)
	switch {
	case isStraight && isFlush:
		return HandEvaluation{Rank: StraightFlush, Score: 8, ScoreCards: straightCards}
	default:
		if nOfAKindCards, ok := getNOfAKind(rankCounts, 4); ok {
			return HandEvaluation{Rank: FourOfAKind, Score: 7, ScoreCards: nOfAKindCards}
		}
		if fullHouseCards, ok := getFullHouse(rankCounts, hand); ok {
			return HandEvaluation{Rank: FullHouse, Score: 6, ScoreCards: fullHouseCards}
		}
		if isFlush {
			return HandEvaluation{Rank: Flush, Score: 5, ScoreCards: hand}
		}
		if isStraight {
			return HandEvaluation{Rank: Straight, Score: 4, ScoreCards: straightCards}
		}
		if nOfAKindCards, ok := getNOfAKind(rankCounts, 3); ok {
			return HandEvaluation{Rank: Triple, Score: 3, ScoreCards: nOfAKindCards}
		}
		if twoPairCards, ok := getTwoPair(rankCounts, hand); ok {
			return HandEvaluation{Rank: TwoPair, Score: 2, ScoreCards: twoPairCards}
		}
		if nOfAKindCards, ok := getNOfAKind(rankCounts, 2); ok {
			return HandEvaluation{Rank: Pair, Score: 1, ScoreCards: nOfAKindCards}
		}
	}
	return HandEvaluation{Rank: HighCard, Score: 0, ScoreCards: hand}
}

func isStraight(hand []cards.Card) ([]cards.Card, bool) {
	if len(hand) < 5 {
		return nil, false
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
			return nil, false
		}
	}
	return hand, true
}

func getNOfAKind(rankCounts map[cards.Rank]int, n int) ([]cards.Card, bool) {
	for rank, count := range rankCounts {
		if count == n {
			return []cards.Card{{Rank: rank}}, true
		}
	}
	return nil, false
}

func getTwoPair(rankCounts map[cards.Rank]int, hand []cards.Card) ([]cards.Card, bool) {
	pairs := []cards.Rank{}
	for rank, count := range rankCounts {
		if count == 2 {
			pairs = append(pairs, rank)
		}
	}
	if len(pairs) == 2 {
		pairCards := []cards.Card{}
		for _, card := range hand {
			if card.Rank == pairs[0] || card.Rank == pairs[1] {
				pairCards = append(pairCards, card)
			}
		}
		return pairCards, true
	}
	return nil, false
}

func getFullHouse(rankCounts map[cards.Rank]int, hand []cards.Card) ([]cards.Card, bool) {
	var threeRank, twoRank cards.Rank
	hasThree := false
	hasTwo := false
	for rank, count := range rankCounts {
		if count == 3 {
			threeRank = rank
			hasThree = true
		}
		if count == 2 {
			twoRank = rank
			hasTwo = true
		}
	}
	if hasThree && hasTwo {
		fullHouseCards := []cards.Card{}
		for _, card := range hand {
			if card.Rank == threeRank || card.Rank == twoRank {
				fullHouseCards = append(fullHouseCards, card)
			}
		}
		return fullHouseCards, true
	}
	return nil, false
}

func EvaluateTwoHands(hand1, hand2 []cards.Card) HandEvaluation {
	hand1Eval := EvaluateHand(hand1)
	hand2Eval := EvaluateHand(hand2)

	if hand1Eval.Rank != hand2Eval.Rank {
		if hand1Eval.Rank > hand2Eval.Rank {
			return hand1Eval
		}
		return hand2Eval
	}

	switch hand1Eval.Rank {
	case Pair, Triple:
		if hand1Eval.ScoreCards[0] == hand2Eval.ScoreCards[0] {
			compareHighCards(hand1, hand2)
		}

	}
	return HandEvaluation{}
}

func compareHighCards(hand1, hand2 []cards.Card) []cards.Card {
	sort.Slice(hand1, func(i, j int) bool { return hand1[i].Rank > hand1[j].Rank })
	sort.Slice(hand2, func(i, j int) bool { return hand2[i].Rank > hand2[j].Rank })

	for i := 0; i < len(hand1); i++ {
		if hand1[i].Rank > hand2[i].Rank {
			return hand1
		} else if hand2[i].Rank > hand1[i].Rank {
			return hand2
		}
	}

	return hand1
}

// TODO: Add functionality to evaluate two equal hands
