package main

import (
	"bufio"   // Used to read text line-by-line from client connections
	"fmt"
	"net"     // Provides TCP networking support
	"strings" // Used to trim user input
	"sync"    // Provides sync.Mutex for safe concurrent access to shared data
)

// Shared map of connected clients (key: connection, value: username)
var (
	clients = make(map[net.Conn]string)
	mu      sync.Mutex // Used to lock access to clients during read/write
)

// handleConnection manages an individual client's connection
func handleConnection(conn net.Conn) {
	defer conn.Close() // Ensure we close the connection when the function exits

	reader := bufio.NewReader(conn)

	// Prompt client for a username
	conn.Write([]byte("Username: "))
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username) // Remove newline or extra spaces

	// Prompt for password
	conn.Write([]byte("Password: "))
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Authenticate the user against the database
	if !Authenticate(username, password) {
		conn.Write([]byte("Authentication failed\n"))
		return // Stop handling this client if credentials are invalid
	}

	// Add client to shared list
	mu.Lock()
	clients[conn] = username
	mu.Unlock()

	// Notify others that a new user has joined
	broadcast(fmt.Sprintf("%s joined the chat\n", username), conn)

	// Main loop: listen for messages from this client
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break // Exit loop if the client disconnects or causes an error
		}
		broadcast(fmt.Sprintf("[%s]: %s", username, msg), conn)
	}

	// Cleanup when client disconnects
	mu.Lock()
	delete(clients, conn)
	mu.Unlock()
	broadcast(fmt.Sprintf("%s left the chat\n", username), conn)
}

// broadcast sends a message to all connected clients *except* the sender
func broadcast(message string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for conn := range clients {
		if conn != sender {
			conn.Write([]byte(message))
		}
	}
}

// main starts the server and listens for incoming TCP connections
func main() {
	StartDB() // Setup SQLite DB and default users

	// Start a TCP server listening on port 6000
	listener, err := net.Listen("tcp", ":6000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Chat server started on port 6000")

	// Infinite loop to accept new client connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue // Skip and wait for the next connection
		}
		go handleConnection(conn) // Handle client in a separate goroutine
	}
}
