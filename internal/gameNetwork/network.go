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
	s.Game.AddPlayer(c.player, s)
	s.broadcastMessage([]byte(fmt.Sprintf("%s has joined the game.", playerName)))

	if len(s.Game.Players) == 2 {
		s.broadcastMessage([]byte("Starting the game with 2 players!"))
		go s.Game.StartGame(s)
	}

	for {

	}
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

func (s *GameServer) getPlayerConnection(playerName string) net.Conn {
	for client := range s.Clients {
		if client.name == playerName {
			return client.conn
		}
	}
	return nil
}

// TODO: Add format message function
