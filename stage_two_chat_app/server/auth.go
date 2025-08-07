package main

import (
	"strings"

	"golang.org/x/crypto/bcrypt" // For password hashing

)
// Authenticate takes a username and password and checks the DB for their validity
func Authenticate(username, password string) bool {
	var dbPassword string

	// Switch the username to lowercase for case-insensitive comparison 
	username = strings.ToLower(username)

	// Query the password stored for this username
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPassword)
	if err != nil {
		return false // User doesn't exist or query failed
	}
	
	// Use bcrypt to compare the provided password with the hashed one
	err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	return err == nil // true if the password matches, false otherwise
}
