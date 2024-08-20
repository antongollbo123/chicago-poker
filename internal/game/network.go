package game

import (
    "fmt"
	"net"
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

    message := "Hello, Server!"
    _, err = conn.Write([]byte(message))
    if err != nil {
        fmt.Println("Error sending message:", err)
        return
    }

    fmt.Println("Message sent:", message)
}


func handleConnection(conn net.Conn) {
    defer conn.Close()

    buf := make([]byte, 1024)
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading from connection:", err)
        return
    }

    fmt.Println("Received:", string(buf[:n]))
}



