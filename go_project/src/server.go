package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	ln, _ := net.Listen("tcp", ":8000")
	fmt.Println("Server running on port 8000")
	for {
		conn, _ := ln.Accept()
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	user := conn.RemoteAddr().String()
	fmt.Fprintf(conn, "Welcome %s!\n", user)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fields := strings.Fields(text)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {
		case "add":
			if len(fields) != 3 {
				fmt.Fprintln(conn, "Usage: add <amount> <category>")
				continue
			}
			addFromInput(user, fields[1], fields[2], conn)
		case "report":
			sendReport(user, conn)
		case "export":
			exportCSV(user, conn)
		default:
			fmt.Fprintln(conn, "Unknown command")
		}
	}
}
