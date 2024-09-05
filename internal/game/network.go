package game

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type client struct {
	name  string
	conn  net.Conn
	in    chan []byte
	out   chan []byte
	token chan struct{}
}
type chat struct {
	clients      map[*client]bool
	messageQueue chan []byte
	exitQueue    chan *client
}

var port string = ":8082"

func NewChat() *chat {
	return &chat{
		clients:      make(map[*client]bool),
		messageQueue: make(chan []byte),
		exitQueue:    make(chan *client),
	}
}

func (c *chat) Start() {
	c.buildServer()
}

func (c *chat) buildServer() {
	server, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("could not start chat: %v", err)
	}
	go c.serve()
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatalf("connection err: %v", err)
			continue
		}
		go c.handleConnection(conn)
	}
}

func (c *chat) exitGuide(client *client) {
	// Send a message to all other clients notifying them the client has left
	leaveMessage := fmt.Sprintf("%s has left the chat.\n", client.name)
	log.Printf("Sending exit message: %s", leaveMessage) // Log exit message
	c.broadcastMessage([]byte(leaveMessage))             // This should trigger a broadcast

	delete(c.clients, client)
	close(client.in)
	close(client.out)
	defer client.conn.Close()
	log.Printf("Client %s has disconnected.", client.name)
}

func (c *chat) serve() {
	fmt.Println("Server is running on port ", port)
	for {
		select {
		case msg := <-c.messageQueue:
			log.Printf("Broadcasting message: %s", msg)
			c.broadcastMessage(msg)
		case client := <-c.exitQueue:
			log.Printf("Processing exit for client: %s", client.name)
			go c.exitGuide(client)
		}
	}
}

func (c *chat) broadcastMessage(msg []byte) {
	for client := range c.clients {
		select {
		case client.out <- msg:
			log.Printf("Sent message to client: %s", client.name)
		default:
			log.Printf("Failed to send message to client %s: channel full or closed", client.name)
		}
	}
}

func (c *chat) handleConnection(conn net.Conn) {
	in, out := make(chan []byte), make(chan []byte)
	client := newClient(conn, in, out)
	c.clients[client] = true
	client.start(c)
}

func newClient(conn net.Conn, in chan []byte, out chan []byte) *client {
	conn.(*net.TCPConn).SetKeepAlive(true)
	conn.(*net.TCPConn).SetKeepAlivePeriod(15 * time.Second)
	return &client{in: in, out: out, conn: conn}
}

func (cl *client) start(c *chat) {
	cl.setUpUsername()
	go cl.receiveMessages(c)
	go cl.sendMessages()
}

func (cl *client) receiveMessages(c *chat) {
	defer func() {
		c.exitQueue <- cl
	}()

	scanner := bufio.NewScanner(cl.conn)
	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Printf("Error reading from client %s: %v", cl.name, scanner.Err())
			break
		}
		msg := fmt.Sprintf("%s: %s\n", cl.name, scanner.Text())
		c.messageQueue <- []byte(msg)
	}
	log.Printf("Client %s has disconnected.", cl.name)
}

func (cl *client) sendMessages() {
	for msg := range cl.out {
		_, err := cl.conn.Write(msg)
		if err != nil {
			log.Printf("Failed to send message to client %s: %v", cl.name, err)
			return
		}
	}
}

func (cl *client) setUpUsername() {
	io.WriteString(cl.conn, "Enter your username: ")
	scann := bufio.NewScanner(cl.conn)
	scann.Scan()
	cl.name = scann.Text()
	io.WriteString(cl.conn, fmt.Sprintf("welcome %s\n", cl.name))
}
