package main

import (
	"github.com/antongollbo123/chicago-poker/internal/gameNetwork"
	"github.com/antongollbo123/chicago-poker/internal/player"
)

func main() {
	// Initialize the GameServer
	gameServer := &gameNetwork.GameServer{
		Clients:      make(map[*gameNetwork.Client]bool),
		MessageQueue: make(chan []byte),
		MoveQueue:    make(chan []byte),
		ExitQueue:    make(chan *gameNetwork.Client),
	}

	// Start the GameServer
	go gameServer.BuildServer()

	// Wait for players to join
	players := make([]*player.Player, 0) // Slice to hold connected players

	// Listen for players and start the game when enough are connected
	go func() {
		for {
			if len(players) >= 4 { // Change this number as needed
				break // Start the game when enough players are connected
			}
			// This is a simple loop to simulate waiting for players
			// In a real application, you would handle this in the handleConnection method
		}

		// Create a new game with the connected players
		newGame := gameNetwork.NewGame(players, gameServer)
		newGame.StartGame() // Start the game
	}()

	// Blocking main goroutine
	select {} // Keeps the main function running
}
