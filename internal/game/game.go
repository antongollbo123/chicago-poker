package game

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/antongollbo123/chicago-poker/internal/deck"
	"github.com/antongollbo123/chicago-poker/internal/player"
)

type Stage string
type Chicago bool

const (
	Poker Stage = "Poker"
	Trick Stage = "Trick"
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
	game.Stage = Poker
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

	// Deal new cards from the deck
	newCards := g.Deck.DrawMultiple(len(indicesToRemove))
	g.Players[playerIndex].Hand = append(g.Players[playerIndex].Hand, newCards...)
}

func (g *Game) StartGame() {

	for g.Stage != Trick {
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
		fmt.Printf("Player %s has new hand: %v\n", player.Name, player.Hand)
	}
	g.EvaluateHands()

	g.Round++

	if g.Round > 2 {
		g.Stage = Trick
	}

}

func (g *Game) TrickRound() {
	//TODO: Add logic for trick round, exit condition should be if all players cards are empty
	fmt.Println("Trick round not yet implemented.")
}

func (g *Game) EvaluateHands() {
	bestScore := -1
	bestPlayerIndex := -1

	for i, player := range g.Players {
		handEvaluation := EvaluateHand(player.Hand)
		fmt.Printf("Player %d has a %v with a score of %d\n", i+1, handEvaluation.Rank, handEvaluation.Score)
		if handEvaluation.Score > bestScore {
			bestScore = handEvaluation.Score
			bestPlayerIndex = i
		}
	}

	if bestPlayerIndex != -1 {
		g.Players[bestPlayerIndex].Score += bestScore
		fmt.Printf("Player %d wins the round and gets %d points\n", bestPlayerIndex+1, bestScore)
	}
}
