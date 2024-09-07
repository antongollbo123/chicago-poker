package gameNetwork

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/antongollbo123/chicago-poker/internal/player"
)

type Node struct {
	GameServer  *GameServer            // Pointer to the active game instance
	Connections map[string]*Connection // Map of player IDs to their connections (network sockets)
}

// A Connection struct represents a player's connection (for sending/receiving messages)
type Connection struct {
	PlayerID string // Unique identifier of the player
	// Any fields needed to manage networking for this connection (e.g., WebSocket, TCP conn, etc.)
}

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
	in     chan []byte
	out    chan []byte
	token  chan struct{}
	player *player.Player
}
type GameServer struct {
	Clients      map[*Client]bool
	currentTurn  *Client
	MessageQueue chan []byte  // General message queue
	MoveQueue    chan []byte  // Queue for handling game moves
	ExitQueue    chan *Client // Clients who disconnect or leave
	Game         *Game        // Your poker game logic
}

var port string = ":8080"

func (s *GameServer) BuildServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting the server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Create a new client and handle the connection
		client := &Client{conn: conn}
		go s.handleConnection(client) // Use goroutine for handling connections
	}
}

func (s *GameServer) handleConnection(c *Client) {
	defer c.conn.Close() // Ensure the connection is closed when function exits
	s.Clients[c] = true  // Add the client to the server's client map

	// Setup player username and associate it with the client
	playerName := c.setUpUsername() // Get the player's name
	if playerName == "" {
		fmt.Println("Failed to get player name")
		return
	}
	c.player = &player.Player{Name: playerName} // Associate the client with the player

	if c.player == nil {
		fmt.Println("Failed to initialize player")
		return
	}

	// Add player to the players slice
	s.Game.Players = append(s.Game.Players, c.player)

	// Notify other players of the new player
	s.broadcastMessage([]byte(fmt.Sprintf("%s has joined the game.", playerName)))

	// Check if the game can be started
	if len(s.Game.Players) >= 4 { // Adjust this as needed
		// Notify everyone that the game is starting
		s.broadcastMessage([]byte("Starting the game with 4 players!"))
		s.Game.StartGame() // Start the game
	}

	// Keep listening for messages from the client
	for {
		buffer := make([]byte, 1024) // Buffer to store incoming messages
		_, err := c.conn.Read(buffer)
		if err != nil {
			break // Exit loop if there’s an error
		}
		// Handle the incoming message (this is where you would implement game actions)
		message := string(buffer)
		s.broadcastMessage([]byte(fmt.Sprintf("%s: %s", playerName, message)))
	}

	// Clean up when the connection is closed
	delete(s.Clients, c)
}

func newClient(conn net.Conn, in chan []byte, out chan []byte) *Client {
	conn.(*net.TCPConn).SetKeepAlive(true)
	conn.(*net.TCPConn).SetKeepAlivePeriod(15 * time.Second)
	return &Client{in: in,
		out:   out,
		conn:  conn,
		token: make(chan struct{}), // Token used to signal turns}
	}
}

func (s *GameServer) serve() {
	fmt.Println("Game server is running on port ", port)
	for {
		select {
		case msg := <-s.MoveQueue:
			var move Message // Assuming msg is of type Message
			if err := json.Unmarshal(msg, &move); err != nil {
				log.Printf("Error unmarshaling move: %v", err)
				continue
			}
			err := s.Game.processMove(move.PlayerName, move.MoveType, move.Data)
			if err != nil {
				log.Printf("Error processing move: %v", err)
			}
		case client := <-s.ExitQueue:
			s.removeClient(client)
		}
	}
}

func (s *GameServer) nextTurn() {
	if len(s.Clients) == 0 {
		log.Println("No clients connected")
		return
	}

	// Store the current turn client to avoid confusion
	currentClient := s.currentTurn

	// Iterate through the clients to find the next player
	for client := range s.Clients {
		if client == currentClient {
			continue // Skip the current player
		}

		s.currentTurn = client // Set the new current turn player

		// Create the message to inform the next player
		msg := Message{PlayerName: client.player.Name, MoveType: NextTurn, Data: "It's your turn"}

		// Send message to the new current player
		err := s.sendMessage(client, msg)
		if err != nil {
			log.Printf("Failed to send next turn message: %v", err)
		}

		client.token <- struct{}{} // Signal the client to act
		break                      // Exit after finding the next player
	}
}

func (s *GameServer) removeClient(client *Client) {
	delete(s.Clients, client) // Remove from the game
	fmt.Printf("Player %s has disconnected.\n", client.name)

	// Handle what happens if a player disconnects in the middle of a game
	// You might want to end the game or handle it differently (e.g., skipping turns)
}

func (cl *Client) setUpUsername() string {
	io.WriteString(cl.conn, "Enter your username: ")
	scanner := bufio.NewScanner(cl.conn)
	scanner.Scan()
	cl.name = scanner.Text()
	io.WriteString(cl.conn, fmt.Sprintf("Welcome, %s\n", cl.name))
	return cl.name
}

func (cl *Client) notifyTurn() {
	io.WriteString(cl.conn, "It's your turn\n")
}

func (s *GameServer) sendMessage(client *Client, msg Message) error {
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return err
	}
	_, err = client.conn.Write(jsonMsg)
	if err != nil {
		log.Printf("Error sending message to client: %v", err)
		return err
	}
	return nil
}

// In network.go
func (s *GameServer) broadcastMessage(msg []byte) {
	for c := range s.Clients {
		// Send the message to each connected client
		_, err := c.conn.Write(msg) // Send the message
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
		}
	}
}

func receiveMessage(playerID string) (Message, error) {
	// Read message from the player's network connection (e.g., WebSocket, TCP)
	// This is pseudo-code. You need to adapt it to your network setup.
	return Message{}, nil
}
