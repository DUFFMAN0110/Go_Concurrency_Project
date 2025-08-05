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

/* 
   handleMessages continuously listens for messages from a single connected user.
   When a message is received, it is broadcast to all other connected clients.
   If the client disconnects or an error occurs, it cleans up and notifies others.
*/
func handleMessages(conn net.Conn, reader *bufio.Reader, username string) {
	for {
		// Wait for the user to send a message
		msg, err := reader.ReadString('\n')
		if err != nil {
			// Likely the client disconnected (e.g., Ctrl+C or closed terminal)
			break
		}

		// Send the message to all other connected users
		broadcast(fmt.Sprintf("[%s]: %s", username, msg), conn)
	}

	// Cleanup: remove user and broadcast departure message
	mu.Lock()
	delete(clients, conn)
	mu.Unlock()
	
	broadcast(fmt.Sprintf("%s left the chat\n", username), conn)
}



func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Ask if user already has an account
	conn.Write([]byte("Do you have an account? (yes/no): \n"))
	choice, _ := reader.ReadString('\n')
	choice = strings.ToLower(strings.TrimSpace(choice))

	if choice == "yes" {
		// Ask for credentials
		conn.Write([]byte("Username:\n"))
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		conn.Write([]byte("Password:\n"))
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		if !Authenticate(username, password) {
			conn.Write([]byte("Authentication failed.\n"))
			return
		}
		conn.Write([]byte("Login successful. Welcome to the chat!\n"))

		// Continue to chat loop...
		mu.Lock()
		clients[conn] = username
		mu.Unlock()

		broadcast(fmt.Sprintf("%s joined the chat\n", username), conn)
		handleMessages(conn, reader, username)
		return
	}

	if choice == "no" {
		// Ask for username and check availability
		conn.Write([]byte("Choose a username:\n"))
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		conn.Write([]byte("Choose a password:\n"))
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		if !CreateUser(username, password) {
			conn.Write([]byte("Registration failed. Username may already exist.\n"))
			return
		}
		conn.Write([]byte("Registration successful. You are now connected to the chat.\n"))

		// Continue to chat loop...
		mu.Lock()
		clients[conn] = username
		mu.Unlock()

		broadcast(fmt.Sprintf("%s joined the chat\n", username), conn)
		handleMessages(conn, reader, username)
		return
	}

	// Handle invalid input
	conn.Write([]byte("Invalid input. Please reconnect and enter yes or no.\n"))
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
