package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var db *sql.DB // Global DB instance used throughout server

// StartDB initializes the SQLite database and adds default users
func StartDB() {
	var err error

	// Open (or create) a SQLite database file
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Create the users table if it doesn't already exist
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	fmt.Println("Database initialized.")

	// Add default users if missing
	createUser("alice", "password123")
	createUser("bob", "secure456")
}

// createUser inserts a new user into the database if they don't already exist
func createUser(username, password string) {
	var exists int

	// Check if username already exists
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&exists)
	if err != nil {
		log.Println("User check error:", err)
		return
	}

	if exists > 0 {
		return // Skip creating the user if they already exist
	}

	// Insert user
	stmt, err := db.Prepare("INSERT INTO users(username, password) VALUES (?, ?)")
	if err != nil {
		log.Println("Statement preparation failed:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, password)
	if err != nil {
		log.Println("Insert failed:", err)
	} else {
		fmt.Printf("User '%s' created.\n", username)
	}
}
