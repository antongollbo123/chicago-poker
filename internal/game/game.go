package game

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/antongollbo123/chicago-poker/internal/deck"
	"github.com/antongollbo123/chicago-poker/internal/player"
	"github.com/antongollbo123/chicago-poker/pkg/cards"
)

type Stage string
type Chicago bool

const (
	Poker Stage = "Poker"
	Trick Stage = "Trick"
)

const (
	TrickWin int = 3
)

type Game struct {
	Deck    *deck.Deck
	Players []*player.Player
	Round   int
	Stage   Stage
}

func NewGame(players []*player.Player) *Game {
	game := Game{}
	game.Round = 0
	game.Stage = Trick // Change back to poker
	game.Players = players
	deck := deck.NewDeck()
	deck.Shuffle()
	game.Deck = deck
	game.Deal()
	return &game
}

func (g *Game) Deal() {
	for _, player := range g.Players {
		cards := g.Deck.DrawMultiple(5)
		player.Hand = cards
	}
}

func (g *Game) TossCards(playerIndex int, indicesToRemove []int) {
	if playerIndex < 0 || playerIndex >= len(g.Players) {
		// TODO: Handle invalid player index
		return
	}

	// Sort indicesToRemove in descending order to safely remove cards from hand slice
	sort.Sort(sort.Reverse(sort.IntSlice(indicesToRemove)))

	// Remove cards from player's hand based on indicesToRemove
	for _, idx := range indicesToRemove {
		if idx >= 0 && idx < len(g.Players[playerIndex].Hand) {
			g.Players[playerIndex].Hand = append(g.Players[playerIndex].Hand[:idx], g.Players[playerIndex].Hand[idx+1:]...)
		}
	}
}

func (g *Game) StartGame() {

	for {
		switch g.Stage {
		case Poker:
			g.PokerRound()
		case Trick:
			g.TrickRound()
		}

	}

}

func (g *Game) PokerRound() {
	fmt.Println("Starting Poker Round: ", g.Round+1)
	scanner := bufio.NewScanner(os.Stdin)
	for i, player := range g.Players {
		fmt.Printf("Player %s, your hand is: %v\n", player.Name, player.Hand)
		fmt.Printf("Enter the indices of the cards you want to toss, separated by spaces: ")
		scanner.Scan()
		input := scanner.Text()
		indicesToRemove := ParseInput(input)
		fmt.Printf("Player %s is tossing cards: %v\n", player.Name, indicesToRemove)
		g.TossCards(i, indicesToRemove)
		// Deal new cards from the deck
		newCards := g.Deck.DrawMultiple(len(indicesToRemove))
		g.Players[i].Hand = append(g.Players[i].Hand, newCards...)

		fmt.Printf("Player %s has new hand: %v\n", player.Name, player.Hand)
	}
	g.EvaluateHands()

	g.Round++

	if g.Round > 2 {
		g.Stage = Trick
	}

}

func (g *Game) TrickRound() {
	fmt.Println("TRICK ROUND!")
	scanner := bufio.NewScanner(os.Stdin)

	var handCopies []cards.Card

	var playedCard cards.Card

	for _, player := range g.Players {
		handCopies = append(handCopies, player.Hand...)

	}
	for i, player := range g.Players {
		fmt.Printf("Player %s, your hand is: %v\n", player.Name, player.Hand)
		fmt.Printf("Enter the index of the card you want to play: ")
		scanner.Scan()
		input := scanner.Text()
		// TODO: Create copies of players hands ?
		//trickHand := g.Players[i].Hand
		indexToPlay := ParseInput(input)
		fmt.Println(playedCard == (cards.Card{}))
		if indexToPlay[0] < len(g.Players[i].Hand) && i == 0 {
			playedCard := g.Players[i].Hand[indexToPlay[0]]
			g.TossCards(i, indexToPlay)
			fmt.Println(playedCard, g.Players[i].Hand)
		} else {

			lastSuit := playedCard.Suit
			condition := func(c cards.Card) bool {
				return c.Suit == lastSuit
			}

			filteredCards := filterCards(g.Players[i].Hand, condition)
			fmt.Println(filteredCards)

		}
	}
}

func (g *Game) EvaluateHands() {
	bestScore := -1
	bestPlayerIndex := -1
	var bestHandEvaluation HandEvaluation
	for i := 0; i < len(g.Players)-1; i++ {
		for j := i + 1; j < len(g.Players); j++ {
			winningHand, winningHandEvaluation := EvaluateTwoHands(g.Players[i].Hand, g.Players[j].Hand)
			if winningHandEvaluation.Score > bestScore {
				bestScore = winningHandEvaluation.Score
				bestHandEvaluation = winningHandEvaluation
				if reflect.DeepEqual(winningHand, g.Players[i].Hand) {
					bestPlayerIndex = i
				} else if reflect.DeepEqual(winningHand, g.Players[j].Hand) {
					bestPlayerIndex = j
				}
			}
		}
	}

	if bestPlayerIndex != -1 {
		g.Players[bestPlayerIndex].Score += bestScore
		fmt.Printf("Player %s wins the round with a %v of %v and gets %d points\n", g.Players[bestPlayerIndex].Name, bestHandEvaluation.Rank, bestHandEvaluation.ScoreCards, bestScore)
	}

}
