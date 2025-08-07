// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"net"
// 	"os"
// )

// // main connects to the chat server and handles input/output
// func main() {
// 	// Establish connection to the TCP server at localhost:6000
// 	conn, err := net.Dial("tcp", "localhost:6000")
// 	if err != nil {
// 		panic(err) // Fail fast if server is unreachable
// 	}
// 	defer conn.Close()

// 	// Goroutine: Listen for incoming messages from server
// 	go func() {
// 		scanner := bufio.NewScanner(conn)
// 		for scanner.Scan() {
// 			fmt.Println(scanner.Text()) // Display server messages
// 		}
// 	}()

// 	// Read user input from terminal and send to server
// 	input := bufio.NewScanner(os.Stdin)
// 	for input.Scan() {
// 		text := input.Text()
// 		if text == "/quit" {
// 			break // Exit the loop and close connection
// 		}
// 		fmt.Fprintf(conn, "%s", text+"\n") // Send user input to server
// 	}
// }


package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	
	"golang.org/x/term"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:6000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Put terminal in raw mode so we can read input byte-by-byte
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)


	done := make(chan struct{}) // Channel to signal connection is closed
	inputChan := make(chan byte)
	// Input buffer to store what the user is typing 
	var inputBuffer []byte
	buf := make([]byte,1) // Used to read one byte at a time from stdin


	// Start listening to server messages
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			incoming := scanner.Text()
			
			// Clear current line: \r moves to beginning, \033[K clears the line
			fmt.Print("\r\033[K") 
			fmt.Println(incoming)
			
			// Reprint the user's prompt and any partially typed input
			fmt.Print("[You]: ")
			fmt.Print(string(inputBuffer))
		}

		// If the server connection is closed, show this message client-side
		fmt.Println("\nYou have been disconnected from the server.")
		close(done)
	}()

	go func(){
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0{
				close(inputChan)
				return
			}
			inputChan <- buf[0]
		}
	}()

	// User input handling
	fmt.Print("[You]: ")

	// Read input one character at a time
	for {
		select{
		case <-done:
			return // exits client back to the original terminal
		case char, ok := <-inputChan:
			if !ok{
				return
			}
			
			switch char {
			case 13: // Enter key (carriage return)
				fmt.Print("\r\n") // Go to next line
				message := string(inputBuffer)
				trimmed := strings.TrimSpace(message)

				// handle the local commands
				if trimmed == "/help"{
					helpMsg := "Here are a list of available commands:\n\t/quit     : Leave the chat and close the connection\n\t/help     : Show this help message again\n\nEnjoy your stay!\n\n"
					fmt.Print(helpMsg)
				}else if trimmed == "/quit" { // if the message is '/quit' close terminal
					fmt.Println("Exiting chat...")
					return 
				}else{
					// Send input (everything is treated as characters/strings) to server
					conn.Write(append(inputBuffer, '\n')) 
				}

				// Reset buffer to reset input
				inputBuffer = inputBuffer[:0]         
				fmt.Print("[You]: ")

			case 127: // Backspace
				if len(inputBuffer) > 0 {
					inputBuffer = inputBuffer[:len(inputBuffer)-1]
					fmt.Print("\b \b") // Erase last character in terminal
				}
			default:
				inputBuffer = append(inputBuffer, char)
				fmt.Printf("%c", char) // Echo typed character
			}
		}
	}
}
