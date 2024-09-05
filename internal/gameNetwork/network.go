package gameNetwork

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

func (c *chat) Start() {
	c.buildServer()
}

var port string = ":8080"

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

func (g *Game) handleConnection(conn net.Conn) {
	in, out := make(chan []byte), make(chan []byte)
	client := newClient(conn, in, out)
	c.clients[client] = true
	client.start(c)
}

func (cl *client) start(c *chat) {
	cl.setUpUsername()
}

func newClient(conn net.Conn, in chan []byte, out chan []byte) *client {
	conn.(*net.TCPConn).SetKeepAlive(true)
	conn.(*net.TCPConn).SetKeepAlivePeriod(15 * time.Second)
	return &client{in: in, out: out, conn: conn}
}

func (c *chat) serve() {
	fmt.Println("Server is running on port ", port)
	for {
		select {
		case msg := <-c.messageQueue:
			log.Printf("Broadcasting message: %s", msg)
			//c.broadcastMessage(msg)
		case client := <-c.exitQueue:
			log.Printf("Processing exit for client: %s", client.name)
			//go c.exitGuide(client)
		}
	}
}

func (cl *client) setUpUsername() (username string) {
	io.WriteString(cl.conn, "Enter your username: ")
	scann := bufio.NewScanner(cl.conn)
	scann.Scan()
	cl.name = scann.Text()
	io.WriteString(cl.conn, fmt.Sprintf("welcome %s\n", cl.name))
	return cl.name
}
