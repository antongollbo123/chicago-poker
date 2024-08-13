# chicago-poker
A command line chicago poker game!


# TODO


1. Add AddPlayer method. (Players should be added after a NewGame is made, NewGame should be instantiated with no players)
   
## Networking logic
 Start on networking logic, Server / Client communication via socket. 
 
1. Starting a server, a client being able to join, a client being able to send a message containing name of player, a new player is created using the sent message of the client

2. Add playerConnections map to Game struct, keeping all players individual connections inside the game struct.
Something like this should happen when a game is started:
``` 
            func NewGame(players []*player.Player) *Game {
                game := Game{}
                game.Round = 0
                game.Stage = Poker
                game.Players = players
                game.playerConnections = make(map[string]net.Conn) // Initialize the map
                
                deck := deck.NewDeck()
                deck.Shuffle()
                game.Deck = deck
                return &game
            }
``` 
4. Add addPlayerConnection method to Game:
``` 
   func (g *Game) AddPlayerConnection(playerName string, conn net.Conn) {
    g.playerConnections[playerName] = conn
}
``` 
6. Add broadcastMessage method to Game, enabling broadcasting of game state to several clients, something like this:
``` 
                func (g *Game) BroadcastMessage(message string) {
            for _, conn := range g.playerConnections {
                if conn != nil {
                    conn.Write([]byte(message + "\n"))
                }
            }
        }
``` 

