package main

// Authenticate takes a username and password and checks the DB
func Authenticate(username, password string) bool {
	var dbPassword string

	// Query the password stored for this username
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&dbPassword)
	if err != nil {
		return false // User doesn't exist or query failed
	}

	// NOTE: Plaintext password check — not secure!
	return dbPassword == password
}
