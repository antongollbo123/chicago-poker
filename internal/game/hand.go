package game

import (
	"reflect"
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

type HandEvaluation struct { // TODO: Rename to Hand ?
	Rank       HandRank
	Score      int
	ScoreCards []cards.Card
	// TODO: Add high card?
}

func EvaluateHand(hand []cards.Card) HandEvaluation {
	rankCounts := make(map[cards.Rank][]cards.Card)
	suitCounts := make(map[cards.Suit]int)

	for _, card := range hand {
		rankCounts[card.Rank] = append(rankCounts[card.Rank], card)
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
		if fullHouseCards, ok := getFullHouse(rankCounts); ok {
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
		if twoPairCards, ok := getTwoPair(rankCounts); ok {
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

func getNOfAKind(rankCounts map[cards.Rank][]cards.Card, n int) ([]cards.Card, bool) {
	for _, cards := range rankCounts {
		if len(cards) == n {
			return cards, true
		}
	}
	return nil, false
}

func getTwoPair(rankCounts map[cards.Rank][]cards.Card) ([]cards.Card, bool) {
	pairs := []cards.Rank{}
	for rank, cards := range rankCounts {
		if len(cards) == 2 {
			pairs = append(pairs, rank)
		}
	}
	if len(pairs) == 2 {
		pairCards := []cards.Card{}
		for _, rank := range pairs {
			pairCards = append(pairCards, rankCounts[rank]...)
		}
		return pairCards, true
	}
	return nil, false
}

func getFullHouse(rankCounts map[cards.Rank][]cards.Card) ([]cards.Card, bool) {
	var threeCards, twoCards []cards.Card
	for _, cards := range rankCounts {
		if len(cards) == 3 {
			threeCards = cards
		} else if len(cards) == 2 {
			twoCards = cards
		}
	}
	if len(threeCards) > 0 && len(twoCards) > 0 {
		return append(threeCards, twoCards...), true
	}
	return nil, false
}

func EvaluateTwoHands(hand1, hand2 []cards.Card) ([]cards.Card, HandEvaluation) {
	hand1Eval := EvaluateHand(hand1)
	hand2Eval := EvaluateHand(hand2)

	// COMPARE RANKS; i.e. PAIR with PAIR, PAIR with TRIPLE etc.
	if hand1Eval.Rank != hand2Eval.Rank {
		if hand1Eval.Rank > hand2Eval.Rank {
			return hand1, hand1Eval
		} else if hand2Eval.Rank > hand1Eval.Rank {
			return hand2, hand2Eval
		}
	}
	// Sort suit must be done before sort cards!
	hand1Eval.ScoreCards = sortCards(hand1Eval.ScoreCards)
	hand2Eval.ScoreCards = sortCards(hand2Eval.ScoreCards)

	hand1 = sortCards(hand1)
	hand2 = sortCards(hand2)

	for i := 0; i < len(hand1Eval.ScoreCards); i++ {
		if hand1Eval.ScoreCards[i].Rank > hand2Eval.ScoreCards[i].Rank {
			return hand1, hand1Eval
		} else if hand1Eval.ScoreCards[i].Rank < hand2Eval.ScoreCards[i].Rank {
			return hand2, hand2Eval
		}
	}
	for i := 0; i < len(hand1); i++ {
		if hand1[i].Rank > hand2[i].Rank {
			return hand1, hand1Eval
		} else if hand1[i].Rank < hand2[i].Rank {
			return hand2, hand2Eval
		}
	}
	winningHand := compareSuit(hand1Eval.ScoreCards, hand2Eval.ScoreCards)

	if reflect.DeepEqual(winningHand, hand1Eval.ScoreCards) {
		return hand1, hand1Eval
	} else if reflect.DeepEqual(winningHand, hand2Eval.ScoreCards) {
		return hand2, hand2Eval
	}
	return nil, HandEvaluation{}
}

func sortCards(cards []cards.Card) []cards.Card {
	cards = sortSuit(cards)
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Rank > cards[j].Rank
	})
	return cards
}

func sortSuit(cards_ []cards.Card) []cards.Card {
	sort.Slice(cards_, func(i, j int) bool {
		return cards.SuitValue(string(cards_[i].Suit)) > cards.SuitValue(string(cards_[j].Suit))
	})
	return cards_

}

func compareSuit(hand1, hand2 []cards.Card) []cards.Card {
	for i := 0; i < len(hand1); i++ {
		if cards.SuitValue(string(hand1[i].Suit)) > cards.SuitValue(string(hand2[i].Suit)) {
			return hand1
		} else if cards.SuitValue(string(hand1[i].Suit)) < cards.SuitValue(string(hand2[i].Suit)) {
			return hand2
		} else {
			continue
		}
	}
	return nil
}

type ConditionFunc func(cards.Card) bool

func filterCards(cardList []cards.Card, condition ConditionFunc) []cards.Card {
	var result []cards.Card
	for _, card := range cardList {
		if condition(card) {

		}
	}
	return result
}
