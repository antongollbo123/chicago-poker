package gameNetwork

import (
	"encoding/json"
	"fmt"
	"net"
)

type Message2 struct {
	Content string `json:"content"`
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		var msg Message2
		err := decoder.Decode(&msg)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		fmt.Printf("Received message: %s\n", msg.Content)

		response := Message2{Content: "Server received: " + msg.Content}
		err = encoder.Encode(response)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
