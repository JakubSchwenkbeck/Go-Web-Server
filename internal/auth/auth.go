package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given password using bcrypt and returns the hash as a string.
// Returns an error if hashing fails.
func HashPassword(password string) (string, error) {
	// Hash the password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err // Return error if hashing fails
	}
	return string(hash), nil // Return the hashed password
}

// Login authenticates a user by username and password.
// If successful, returns a hashed token; otherwise, returns "invalid".
func Login(cs ChatService, username string, password string) string {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	id, found := findValueInMap(cs.users, username)
	var token = "invalid"
	if found && password == cs.users[id].Password {
		// Successful login
		fmt.Print("Logged in as " + username)
		token, _ = HashPassword(password) // Create token
	} else {
		// Login failed
		fmt.Print("Login failed for " + username)
	}

	return token
}

// findValueInMap looks for a user with a specific name in the map.
// Returns the user ID and a boolean indicating if found.
func findValueInMap(m map[string]ChatUser, targetValue string) (string, bool) {
	for key, value := range m {
		if value.InternData.Name == targetValue {
			return key, true // User found
		}
	}
	return "", false // User not found
}
