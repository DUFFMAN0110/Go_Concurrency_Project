package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// main connects to the chat server and handles input/output
func main() {
	// Establish connection to the TCP server at localhost:6000
	conn, err := net.Dial("tcp", "localhost:6000")
	if err != nil {
		panic(err) // Fail fast if server is unreachable
	}
	defer conn.Close()

	// Goroutine: Listen for incoming messages from server
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text()) // Display server messages
		}
	}()

	// Read user input from terminal and send to server
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		text := input.Text()
		if text == "/quit" {
			break // Exit the loop and close connection
		}
		fmt.Fprintf(conn, text+"\n") // Send user input to server
	}
}
