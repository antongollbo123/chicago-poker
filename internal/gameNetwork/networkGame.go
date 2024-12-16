package gameNetwork

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
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
			g.TrickRound(server)
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
		// Send the player's hand
		handMsg := Message{
			PlayerName: player.Name,
			MoveType:   GameUpdate,
			Data:       player.Hand,
		}
		g.notifyServer(server, handMsg)

		// Ask for the player's move
		promptMsg := Message{
			PlayerName: player.Name,
			MoveType:   PokerToss,
			Data:       "Enter the indices of the cards you want to toss, separated by spaces: ",
		}
		g.notifyServer(server, promptMsg)

	}
	bestPlayerIndex, bestHandEvaluation := g.EvaluateHands()
	formattedMsg := fmt.Sprintf("Player %s wins the round with a %v of %v and gets %d points\n",
		g.Players[bestPlayerIndex].Name,
		bestHandEvaluation.Rank,
		bestHandEvaluation.ScoreCards,
		bestHandEvaluation.Score)
	msgBytes := []byte(formattedMsg)
	server.broadcastMessage(msgBytes)
	g.Round++

	if g.Round > 2 {
		g.Stage = Trick
	}
}

func (g *Game) TrickRound(server *GameServer) {
	server.broadcastMessage([]byte("TRICK ROUND!"))
	leadIndex := g.leadIndex

	for trick := 0; trick < 5; trick++ {
		var leadCard cards.Card
		playedCards := make([]cards.Card, len(g.Players))
		indicesToRemove := make([]int, len(g.Players))

		server.broadcastMessage([]byte(fmt.Sprintf("Starting trick %d\n", trick+1)))

		for i := 0; i < len(g.Players); i++ {
			playerIndex := (leadIndex + i) % len(g.Players)
			currentPlayer := g.Players[playerIndex]

			// Notify player of their hand
			handMsg := Message{
				PlayerName: currentPlayer.Name,
				MoveType:   GameUpdate,
				Data:       currentPlayer.Hand,
			}
			g.notifyServer(server, handMsg)

			// Ask for the player's move
			promptMsg := Message{
				PlayerName: currentPlayer.Name,
				MoveType:   TrickPlay,
				Data:       "Enter the index of the card you want to play: ",
			}
			cardIndex := g.notifyServer(server, promptMsg)

			if len(cardIndex) != 1 {
				fmt.Println("Invalid card index. Try again.")
				i-- // Retry for the same player
				continue
			}

			playedCard := currentPlayer.Hand[cardIndex[0]]

			if i == 0 {
				leadCard = playedCard // Set the lead card for the trick
			} else if !isValidTrickMove(currentPlayer, playedCard, leadCard) {
				fmt.Println("Invalid move. You must follow the suit if possible.")
				i-- // Retry for the same player
				continue
			}

			playedCards[playerIndex] = playedCard
			indicesToRemove[playerIndex] = cardIndex[0]

			fmt.Printf("Player %s played %v\n", currentPlayer.Name, playedCard)
		}

		winnerIndex := findWinner(playedCards, leadIndex)
		fmt.Printf("Player %s wins the trick with %v\n", g.Players[winnerIndex].Name, playedCards[winnerIndex])

		// Remove played cards
		for playerIdx, cardIdx := range indicesToRemove {
			if cardIdx >= 0 {
				g.TossCards(playerIdx, []int{cardIdx})
			}
		}

		leadIndex = winnerIndex // Update lead index for the next trick
	}

	// Award points to the player who wins the final trick
	g.Players[leadIndex].Score += TrickWin
	g.Round++
	fmt.Printf("Player %s wins the trick round and gets %d points\n", g.Players[leadIndex].Name, TrickWin)

	g.Stage = Poker // Switch back to Poker round
	g.Deal()
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

func describe(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func (g *Game) processMove(playerName string, moveType MessageType, data interface{}) error {
	log.Printf("Processing move: Player: %s, Type: %s, Data: %v\n", playerName, moveType, data)

	fmt.Print(moveType == "poker_toss")
	fmt.Print(moveType == PokerToss)
	describe(data)

	intIndices, ok := data.([]int)
	fmt.Print(intIndices)
	if !ok {
		return fmt.Errorf("invalid data format")
	}

	playerIndex := g.getPlayerIndex(playerName)
	if playerIndex == -1 {
		return fmt.Errorf("player not found")
	}

	if moveType == "poker_toss" {
		g.TossCards(playerIndex, intIndices)
		newCards := g.Deck.DrawMultiple(len(intIndices))
		g.Players[playerIndex].Hand = append(g.Players[playerIndex].Hand, newCards...)
		log.Printf("Player %s has new hand: %v\n", playerName, g.Players[playerIndex].Hand)
	}
	return nil
}

func (g *Game) processTrickMove(moveType MessageType, data interface{}) []int {
	if moveType == "trick_play" {
		intIndices, _ := data.([]int)
		fmt.Println("intindices in processTrickMove", intIndices)
		return intIndices
	}
	return nil
}

func isValidTrickMove(player *player.Player, playedCard cards.Card, leadCard cards.Card) bool {
	// Check if the player has a card of the lead suit
	if hasSuit(player.Hand, leadCard.Suit) {
		// If the played card matches the lead suit, it's valid
		return playedCard.Suit == leadCard.Suit
	}
	// If the player does not have the lead suit, any card is valid
	return true
}

func (g *Game) notifyServer(server *GameServer, msg Message) []int {

	fmt.Println("I AM IN NOTIFYSERVER, MOVE TYPE: ", msg.MoveType)
	playerConn := server.getPlayerConnection(msg.PlayerName)

	if playerConn == nil {
		fmt.Printf("No connection found for player %s\n", msg.PlayerName)
		return nil
	}

	if msg.MoveType == "poker_toss" {
		reader := bufio.NewReader(playerConn)
		playerConn.Write([]byte("Enter your move: "))
		content, _ := reader.ReadString('\n')
		content = strings.TrimSpace(content)
		msg.MoveType = PokerToss // TODO: Redundant ? --> Remove?
		msg.Data = game.ParseInput(content)
		fmt.Println(msg.Data)
		g.processMove(msg.PlayerName, msg.MoveType, msg.Data)
		return game.ParseInput(content)
	}

	if msg.MoveType == "trick_play" {
		reader := bufio.NewReader(playerConn)
		playerConn.Write([]byte("Enter your move (single card index): "))
		content, _ := reader.ReadString('\n')
		content = strings.TrimSpace(content)

		parsedCardIndex := game.ParseInput(content)
		if len(parsedCardIndex) != 1 {
			fmt.Println("Invalid input. Please enter a single valid card index.")
			return nil
		}

		if parsedCardIndex[0] < 0 || parsedCardIndex[0] >= len(g.Players[g.getPlayerIndex(msg.PlayerName)].Hand) {
			fmt.Println("Index out of range. Try again.")
			return nil
		}

		msg.Data = parsedCardIndex
		return parsedCardIndex
	}

	// Marshal the message to JSON
	jsonData, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error marshaling message for player %s: %v\n", msg.PlayerName, err)
		return nil
	}

	// Send the message over the player's connection
	playerConn.Write([]byte("JSON encoded message: "))
	_, err = playerConn.Write(jsonData)
	if err != nil {
		fmt.Printf("Failed to send message to player %s: %v\n", msg.PlayerName, err)
		return nil
	}
	return nil
}

func (g *Game) getPlayerIndex(playerName string) int {
	for i, p := range g.Players {
		if p.Name == playerName {
			return i
		}
	}
	return -1
}
