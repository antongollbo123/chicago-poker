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

func (g *Game) getHighScore() int {
	highScore := 0
	for _, player := range g.Players {
		if player.Score > highScore {
			highScore = player.Score
		}
	}
	fmt.Println("high score is: ", highScore)
	return highScore
}

func (g *Game) StartGame() {
	g.Deal()
	for g.getHighScore() < 50 {
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
	bestPlayerIndex, bestHandEvaluation := g.EvaluateHands()

	fmt.Printf("Player %s wins the round with a %v of %v and gets %d points\n", g.Players[bestPlayerIndex].Name, bestHandEvaluation.Rank, bestHandEvaluation.ScoreCards, bestHandEvaluation.Score)

	g.Round++

	if g.Round > 2 {
		g.Stage = Trick
	}

}

func (g *Game) TrickRound() {
	fmt.Println("TRICK ROUND!")
	scanner := bufio.NewScanner(os.Stdin)
	// Make copies of hands
	bestPlayerIndex, bestHandEvaluation := g.EvaluateHands()

	leadIndex := 0
	for trick := 0; trick < 5; trick++ {
		var leadCard cards.Card
		playedCards := make([]cards.Card, len(g.Players))
		indicesToRemove := make([][]int, len(g.Players))

		fmt.Printf("Starting trick %d\n", trick+1)
		for i := 0; i < len(g.Players); i++ {
			currentIndex := (leadIndex + i) % len(g.Players)
			player := g.Players[currentIndex]

			fmt.Printf("Player %s, your hand is: %v\n", player.Name, player.Hand)
			fmt.Printf("Enter the index of the card you want to play: ")
			scanner.Scan()
			input := scanner.Text()
			indexToPlay := ParseInput(input)

			if len(indexToPlay) != 1 || indexToPlay[0] < 0 || indexToPlay[0] >= len(player.Hand) {
				fmt.Println("Invalid card index. Try again.")
				i-- // Ask the same player to play again
				continue
			}

			playedCard := player.Hand[indexToPlay[0]]

			if i == 0 {
				leadCard = playedCard
			} else {
				if playedCard.Suit != leadCard.Suit && hasSuit(player.Hand, leadCard.Suit) {
					fmt.Println("You must follow the suit. Try again.")
					i-- // Ask the same player to play again
					continue
				}
			}

			playedCards[currentIndex] = playedCard
			indicesToRemove[currentIndex] = []int{indexToPlay[0]}
			fmt.Printf("Player %s played %v\n", player.Name, playedCard)
		}

		winnerIndex := findWinner(playedCards, leadIndex)
		fmt.Printf("Player %s wins the trick with %v\n", g.Players[winnerIndex].Name, playedCards[winnerIndex])

		// Remove played cards from players' hands using TossCards
		for i := range indicesToRemove {
			if len(indicesToRemove[i]) > 0 {
				g.TossCards(i, indicesToRemove[i])
			}
		}

		leadIndex = winnerIndex
	}
	// Award points to the player who wins the last trick
	g.Players[leadIndex].Score += TrickWin
	g.Round++
	fmt.Printf("Player %s wins the trick round and gets %d points\n", g.Players[leadIndex].Name, TrickWin)
	g.Stage = Poker
	g.Deal()
	fmt.Printf("Player %s wins the final poker round with a %v of %v and gets %d points\n", g.Players[bestPlayerIndex].Name, bestHandEvaluation.Rank, bestHandEvaluation.ScoreCards, bestHandEvaluation.Score)
}

// Check if a player has a card of a given suit
func hasSuit(hand []cards.Card, suit cards.Suit) bool {
	for _, card := range hand {
		if card.Suit == suit {
			return true
		}
	}
	return false
}

// Find the winner of the current trick
func findWinner(playedCards []cards.Card, leadIndex int) int {
	leadSuit := playedCards[leadIndex].Suit
	highestCard := playedCards[leadIndex]
	highestIndex := leadIndex

	for i, card := range playedCards {
		if card.Suit == leadSuit && card.Rank > highestCard.Rank {
			highestCard = card
			highestIndex = i
		}
	}

	return highestIndex
}

func (g *Game) EvaluateHands() (int, HandEvaluation) {
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
		return bestPlayerIndex, bestHandEvaluation
	}
	return 0, HandEvaluation{}
}
