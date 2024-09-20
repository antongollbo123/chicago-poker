package gameNetwork

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/antongollbo123/chicago-poker/internal/player"
)

type MessageType string

const (
	PokerToss    MessageType = "poker_toss"
	TrickPlay    MessageType = "trick_play"
	GameUpdate   MessageType = "game_update"
	NextTurn     MessageType = "next_turn"
	PlayerJoined MessageType = "player_joined"
)

type Message struct {
	PlayerName string      `json:"player_name"`
	MoveType   MessageType `json:"move_type"`
	Data       interface{} `json:"data"`
}

type Client struct {
	name   string
	conn   net.Conn
	player *player.Player
}

type GameServer struct {
	Clients map[*Client]bool
	Game    *Game
}

func (s *GameServer) BuildServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 8080...")
	s.Game = NewGame([]*player.Player{})

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		client := &Client{conn: conn}
		go s.handleConnection(client)
	}
}

func (s *GameServer) handleConnection(c *Client) {
	defer c.conn.Close()
	s.Clients[c] = true
	print("connection added", c.conn)
	playerName := c.setUpUsername()
	if playerName == "" {
		fmt.Println("Failed to get player name")
		return
	}
	c.player = &player.Player{Name: playerName}
	print("HERE IS THE NAME:", c.player.Name)
	// Notify the game to add the player
	s.Game.AddPlayer(c.player, s)

	// Notify other players of the new player
	s.broadcastMessage([]byte(fmt.Sprintf("%s has joined the game.", playerName)))

	// Check if the game can be started
	if len(s.Game.Players) >= 2 {
		s.broadcastMessage([]byte("Starting the game with 2 players!"))
		s.Game.StartGame(s)
	}

	for {
		buffer := make([]byte, 1024)
		n, err := c.conn.Read(buffer)
		if err != nil {
			break
		}

		var msg Message
		print("HERE IS THE ATTEMPTED MESSAGE:")
		if err := json.Unmarshal(buffer[:n], &msg); err != nil {
			fmt.Println("Error parsing message:", err)
			continue
		}

		// Instead of broadcasting all messages, process moves from the player
		err = s.Game.processMove(msg.PlayerName, msg.MoveType, msg.Data)
		if err != nil {
			fmt.Printf("Error processing move from player %s: %v\n", msg.PlayerName, err)
			continue
		}
	}
	delete(s.Clients, c)
}

func (c *Client) setUpUsername() string {

	if c.conn == nil {
		fmt.Println("Error: Connection is nil for client.")
		return ""
	}
	io.WriteString(c.conn, "Enter your username: ")
	scanner := bufio.NewScanner(c.conn)
	scanner.Scan()
	c.name = scanner.Text()
	//io.WriteString(c.conn, fmt.Sprintf("Welcome, %s\n", c.name))
	return c.name
}

func (s *GameServer) broadcastMessage(msg []byte) {
	for c := range s.Clients {
		_, err := c.conn.Write(msg)
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
		}
	}
}

func (s *GameServer) sendMessageToPlayer(playerName string, msg Message) error {
	for client := range s.Clients { // Iterate over all connected clients
		if client.player.Name == playerName { // Match by player name
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				return err
			}
			_, err = client.conn.Write(jsonMsg) // Send the message
			return err
		}
	}
	return fmt.Errorf("Client for player %s not found", playerName)
}

// TODO: Add format message function
