package gameNetwork

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/antongollbo123/chicago-poker/internal/deck"
	"github.com/antongollbo123/chicago-poker/internal/game"
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
	Deck      *deck.Deck
	Players   []*player.Player
	Round     int
	Stage     Stage
	leadIndex int
}

func NewGame(players []*player.Player) *Game {
	game := Game{
		Players: players,
		Round:   0,
		Stage:   Poker,
	}

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

func (g *Game) StartGame(server *GameServer) {
	g.Deal()
	for g.getHighScore() < 50 {
		switch g.Stage {
		case Poker:
			g.PokerRound(server)
		case Trick:
			g.TrickRound()
		}

	}

}

func (g *Game) AddPlayer(player *player.Player, server *GameServer) {
	if g == nil {
		fmt.Println("Error: Game instance is nil.")
		return
	}
	print("PLAYER NAME INSIDE: ", player.Name)
	g.Players = append(g.Players, player)
	fmt.Printf("Player %s has been added to the game.\n", player.Name)

	// Optionally notify the server or other players about the new player
	msg := Message{
		PlayerName: player.Name,
		MoveType:   PlayerJoined,
		Data:       fmt.Sprintf("%s has joined the game!", player.Name),
	}
	g.notifyServer(server, msg)
}

func (g *Game) PokerRound(server *GameServer) {
	fmt.Println("Starting Poker Round: ", g.Round+1)

	for _, player := range g.Players {
		// Show player's hand
		msg := Message{
			PlayerName: player.Name,
			MoveType:   GameUpdate,
			Data:       player.Hand,
		}
		g.notifyServer(server, msg) // Notify server to send to player

		// Prompt for cards to toss
		msg = Message{
			PlayerName: player.Name,
			MoveType:   NextTurn,
			Data:       "Enter the indices of the cards you want to toss, separated by spaces: ",
		}
		g.notifyServer(server, msg)

		input := waitForInput(player.Name)

		log.Printf("Input inside pokerRound before sending: %v\n", input)

		g.notifyServer(server, input)
	}

	// Rest of the method remains unchanged...
}

func waitForInput(playerName string) Message {
	fmt.Printf("%s, please enter your move: ", playerName)

	var input string = "0 1 2" // Hardcoded for testing
	/*
		if err != nil {
			log.Printf("Error reading input: %v\n", err)
			return Message{}
		}
	*/
	// Split input into a slice of indices (assuming space-separated)
	indices := strings.Fields(input)

	// Create the message struct
	msg := Message{
		PlayerName: playerName,
		MoveType:   PokerToss, // Set the correct move type
		Data:       indices,   // Indices should be sent as a slice of strings
	}

	fmt.Printf("%s, your move ", indices)

	return msg
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
			indexToPlay := game.ParseInput(input)

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

func (g *Game) EvaluateHands() (int, game.HandEvaluation) {
	bestScore := -1
	bestPlayerIndex := -1
	var bestHandEvaluation game.HandEvaluation
	for i := 0; i < len(g.Players)-1; i++ {
		for j := i + 1; j < len(g.Players); j++ {
			winningHand, winningHandEvaluation := game.EvaluateTwoHands(g.Players[i].Hand, g.Players[j].Hand)
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
	return 0, game.HandEvaluation{}
}

func (g *Game) processMove(playerName string, moveType MessageType, data interface{}) error {
	var playerIndex int = -1

	// Find the player making the move
	for i, p := range g.Players {
		if p.Name == playerName {
			playerIndex = i
			break
		}
	}

	if playerIndex == -1 {
		return fmt.Errorf("Player not found")
	}

	switch moveType {
	case PokerToss:
		indicesToRemove, ok := data.([]int)
		if !ok {
			return fmt.Errorf("Invalid data format for poker_toss")
		}

		g.TossCards(playerIndex, indicesToRemove)
		newCards := g.Deck.DrawMultiple(len(indicesToRemove))
		g.Players[playerIndex].Hand = append(g.Players[playerIndex].Hand, newCards...)
		fmt.Printf("Player %s has new hand: %v\n", g.Players[playerIndex].Name, g.Players[playerIndex].Hand)

	case TrickPlay:
		indexToPlay, ok := data.(int)
		if !ok {
			return fmt.Errorf("Invalid data format for trick_play")
		}

		player := g.Players[playerIndex]
		if indexToPlay < 0 || indexToPlay >= len(player.Hand) {
			return fmt.Errorf("Invalid card index")
		}

		playedCard := player.Hand[indexToPlay]

		if g.Stage == Trick {
			leadCard := g.Players[g.leadIndex].Hand[0] // Lead card of the round
			if !isValidTrickMove(player, playedCard, leadCard) {
				return fmt.Errorf("Invalid trick move: must follow suit")
			}
		}

		g.TossCards(playerIndex, []int{indexToPlay})
		fmt.Printf("Player %s played %v\n", g.Players[playerIndex].Name, playedCard)

	default:
		return fmt.Errorf("Unknown move type: %s", moveType)
	}

	return nil
}

func isValidTrickMove(player *player.Player, playedCard cards.Card, leadCard cards.Card) bool {
	// If the played card matches the lead suit, it's valid
	if playedCard.Suit == leadCard.Suit {
		return true
	}
	// Check if the player has a card of the lead suit
	return hasSuit(player.Hand, leadCard.Suit)
}

func (g *Game) notifyServer(server *GameServer, msg Message) {
	// Notify server to send message to the specific player
	fmt.Print("Inside notify")

	// Assuming server has a method to get the connection of a player
	playerConn := server.getPlayerConnection(msg.PlayerName)
	if playerConn == nil {
		fmt.Printf("No connection found for player %s\n", msg.PlayerName)
		return
	}

	// Marshal the message to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error marshaling message for player %s: %v\n", msg.PlayerName, err)
		return
	}
	print(jsonData)
	// Send the message over the player's connection
	_, err = playerConn.Write(jsonData)
	if err != nil {
		fmt.Printf("Failed to send message to player %s: %v\n", msg.PlayerName, err)
	}
}
