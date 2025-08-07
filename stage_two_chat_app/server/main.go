package main

import (
	"bufio" // Used to read text line-by-line from client connections
	"fmt"	// used for printing to the terminal 
	"net"     // Provides TCP networking support
	"strings" // Used to trim user input
	"sync"    // Provides sync.Mutex for safe concurrent access to shared data
)

// Shared map of connected clients (key: connection, value: username)
var (
	clients = make(map[net.Conn]string) // Tracks connected clients by their connection
	mu      sync.Mutex // Used to lock access to clients during read/write
	chatHistory []string
	historyLock sync.Mutex
)



/*
handleMessages continuously listens for messages from a single connected user.
When a message is received, it is broadcast to all other connected clients.
If the client disconnects or an error occurs, it cleans up and notifies others.
*/

func handleMessages(conn net.Conn, reader *bufio.Reader, username string) {
	for {
		// Now read the actual message (blocks until user presses Enter)
		msg, err := reader.ReadString('\n')
		if err != nil {
			// If there's an error (likely user disconnected), break the loop
			break
		}
		finalMsg := fmt.Sprintf("[%s]: %s", username, msg)

		// save message to chat history
		historyLock.Lock()
		chatHistory = append(chatHistory, strings.TrimRight(finalMsg,"\n"))
		historyLock.Unlock()


		// Broadcast the message to all other connected clients
		// Use the actual username in the broadcasted message (not "[You]")
		broadcast(finalMsg, conn)
	}

	// If we get here, the user has disconnected or an error occurred
	// Clean up: remove the client from the shared clients map
	disconnectClient(conn)
}

func disconnectClient(conn net.Conn) {
	mu.Lock()
	if username, ok := clients[conn]; ok {
		delete(clients, conn)
		mu.Unlock()
		
		// Notify other users that this user left the chat
		broadcast(fmt.Sprintf("%s left the chat\n", username), conn)
	} else {
		mu.Unlock()
	}
	conn.Close()
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	var choice string
	
	// Ask if user already has an account (loop until valid yes/no input)
	for{
		conn.Write([]byte("Do you have an account? (yes/no): \n"))
		rawChoice, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		choice = strings.ToLower(strings.TrimSpace(rawChoice))
		if choice == "yes" || choice == "no"{
			break
		}
		conn.Write([]byte("Invalid input. Please type 'yes' or 'no'.\n"))
	}


	if choice == "yes" {
		// Ask for credentials, authenticate, if after 5 authentification attempts
		
		const maxAttempts = 5
		var username, password string
		authenticated := false

		for attempts := 1; attempts <= maxAttempts; attempts++ {
			conn.Write([]byte("Username:\n"))
			u, _ := reader.ReadString('\n')
			username = strings.TrimSpace(u)
			
			if isUserOnline(username) {
				conn.Write([]byte(fmt.Sprintf("This user is already logged in elsewhere. Attempts left: %d\n", maxAttempts-attempts)))
				if attempts == maxAttempts{
					break
				}else{
					continue
				}
			}

			conn.Write([]byte("Password:\n"))
			p, _ := reader.ReadString('\n')
			password = strings.TrimSpace(p)

			if Authenticate(username, password) {
				authenticated = true
				break
			} else {
				conn.Write([]byte(fmt.Sprintf("Authentication failed. Attempts left: %d\n", maxAttempts-attempts)))
			}
		}

		if !authenticated {
			conn.Write([]byte("Too many failed attempts. Connection closing.\n"))
			disconnectClient(conn)
			return
		}

		conn.Write([]byte("Login successful. Welcome to the chat!\n"))

		// actually add the user to the map of users
		mu.Lock()
		clients[conn] = username
		mu.Unlock()
		
		// send info 
		sendWelcomeMessage(conn)
		sendChatHistory(conn)
		
		
		// broadcast to the other users which user joined the connection
		broadcast(fmt.Sprintf("%s joined the chat\n", username), conn)
		handleMessages(conn, reader, username)

		// return

	} else if choice == "no" {
		
		for{
			// Ask for username and check availability
			conn.Write([]byte("Choose a username:\n"))
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)
			
			// checks if the username is an empty string, 
			if username == "" || len(username) < 1{
				conn.Write([]byte("This username is too short. Try a different one\n"))
				continue
			}
			conn.Write([]byte("Choose a password (must be at least 5 characters):\n"))
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)
			
			// If you were able to create a new username and password
			if CreateUser(username, password) {
				conn.Write([]byte("Registration successful. You are now connected to the chat.\n"))
				
				// actually add the user to the map of users
				mu.Lock()
				clients[conn] = username
				mu.Unlock()
				
				// send info 
				sendWelcomeMessage(conn)
				sendChatHistory(conn)
				
				
				// broadcast to the other users which user joined the connection
				broadcast(fmt.Sprintf("%s joined the chat\n", username), conn)
				handleMessages(conn, reader, username)
			}
			
			conn.Write([]byte("Registration failed. Username may already exist, your password may be incorrect, or password is too short (minimum 5 characters).\n"))
		}
	
		
	}
}

func isUserOnline(username string) bool {
	mu.Lock()
	defer mu.Unlock()
	for _, name := range clients {
		if strings.EqualFold(name, username) {
			return true
		}
	}
	return false
}

func sendChatHistory(conn net.Conn){
	historyLock.Lock()
	defer historyLock.Unlock()

	if len(chatHistory) == 0{
		return
	}

	conn.Write([]byte("----- Chat History -----\n"))
	for _, msg := range chatHistory{
		conn.Write([]byte(msg + "\n"))
	}
	conn.Write([]byte("------------------------\n"))
}

func sendWelcomeMessage(conn net.Conn){
	msg := `
Welcome to the chat !

Here are a list of available commands:
	/quit	: Leave the chat and close the connection
	/help	: Show this help message again
	
Enjoy your stay!
`
	conn.Write([]byte(msg))
}


// broadcast sends a message to all connected clients *except* the sender
func broadcast(message string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	// for every client connected, if they aren't the sender of the message,
	// write the message to their terminal
	for conn := range clients {
		if conn != sender {
			conn.Write([]byte(message))
		}
	}
}

// main starts the server and listens for incoming TCP connections
func main() {
	StartDB() // Setup SQLite DB and default users (you need default users to access the db)

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
