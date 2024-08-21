package game

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func StartServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server started on port 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func StartClient() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	initializeClient(conn)
	startChat(conn)
}

func initializeClient(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your application name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	_, err := conn.Write([]byte("NAME:" + name + "\n"))
	if err != nil {
		fmt.Println("Error sending name:", err)
		return
	}

	fmt.Println("Name sent:", name)
}

func startChat(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	go func() {
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading from server:", err)
				return
			}
			fmt.Print("Server: " + msg)
		}
	}()

	for {
		fmt.Print("You: ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		if strings.ToLower(userInput) == "exit" {
			fmt.Println("Exiting chat...")
			return
		}

		_, err := conn.Write([]byte("MESSAGE:" + userInput + "\n"))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection: ", err)
			return
		}
		message = strings.TrimSpace(message)

		if strings.HasPrefix(message, "NAME:") {
			name := strings.TrimPrefix(message, "NAME:")
			fmt.Println(name, " has joined the game!")
		} else if strings.HasPrefix(message, "MESSAGE:") {
			msg := strings.TrimPrefix(message, "MESSAGE:")
			fmt.Println("Received message:", msg)
			// Send a response back to the client
			_, err := conn.Write([]byte("ACK: " + msg + " received\n"))
			if err != nil {
				fmt.Println("Error sending acknowledgment:", err)
				return
			}
		}
	}
}
