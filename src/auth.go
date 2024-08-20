package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash of the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func Login(cs ChatService, username string, password string) {
	cs.mu.Lock()
	id, found := findValueInMap(cs.users, username)
	if found { // if user exists

		if password == cs.users[id].Password { // if passwords match

			// login
			fmt.Print("You now are Logged in as {}")

		}

	}

	cs.mu.Unlock()

}
func findValueInMap(m map[string]ChatUser, targetValue string) (string, bool) {
	for key, value := range m {
		if value.InternData.Name == targetValue {
			return key, true // Return the key and a boolean indicating the value was found
		}
	}
	return "", false // Return an empty string and false if the value is not found
}
