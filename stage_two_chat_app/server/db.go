package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"golang.org/x/crypto/bcrypt"
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

	// Add default user if missing users
	CreateUser("z", "z12345678")
}

/*
   CreateUser attempts to register a new user in the database with the given username and password.
   It first checks if the username is already taken. If not, it inserts the new user.
   Returns true if registration succeeds, false if the username exists or an error occurs.
*/

func CreateUser(username, password string) bool {

	username = strings.ToLower(strings.TrimSpace(username)) // make the username lowercase

	if username == "" || len(username) < 1{
		return false
	}
	//Checks if the username already exists
	if len(password) < 5{
		return false
	}

	var exists int

	// Run a query to count how many users exist with the given username
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&exists)

	if err != nil {
		// If there's an error during query execution (e.g., DB not reachable), log it and return false
		fmt.Println("Error checking user existence:", err)
		return false
	}

	if exists > 0 {
		// If count > 0, the username is already taken — registration should fail
		return false
	}

	// Hash the password securely with bcrypt
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Print("Error hasing password:", err)
	}

	// Prepare a SQL statement for inserting the new user
	stmt, err := db.Prepare("INSERT INTO users(username, password) VALUES (?, ?)")

	if err != nil {
		// If preparing the statement fails, return false
		fmt.Println("Error preparing insert statement:", err)
		return false
	}

	defer stmt.Close() // Ensure we close the statement no matter what (best practice)

	// Execute the INSERT statement with the provided username and password
	_, err = stmt.Exec(username, string(hashed))

	if err != nil {
		// If executing the insert fails (e.g., DB constraint error), return false
		fmt.Println("Error inserting new user:", err)
		return false
	}

	// Registration successful, user created
	return true
}
