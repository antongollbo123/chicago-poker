package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/antongollbo123/chicago-poker/internal/game"
)

func main() {
	mode := flag.String("mode", "", "start as 'server' or 'client'")
	flag.Parse()

	if *mode == "server" {
		game.StartServer()
	} else if *mode == "client" {
		game.StartClient()
	} else {
		fmt.Println("Please specify a mode: -mode=server or -mode=client")
		os.Exit(1)
	}
}

/*p1 := player.NewPlayer("Anton")
p2 := player.NewPlayer("Nora")
p3 := player.NewPlayer("Niklas")
gamet := game.NewGame([]*player.Player{p1, p2, p3})

gamet.StartGame()
*/
