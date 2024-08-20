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

func Login(cs ChatService, username string, password string) string {
	cs.mu.Lock()
	id, found := findValueInMap(cs.users, username)
	var token = "invalid"
	if found { // if user exists

		if password == cs.users[id].Password { // if passwords match

			// login
			fmt.Print("You now are Logged in as " + username)
			token, _ = HashPassword(password)
		} else {
			fmt.Print("Login to " + username + " was not successfull!")

		}

	} else {
		fmt.Print("User " + username + " does not exist yet. Please register are try with another name!")

	}

	cs.mu.Unlock()
	return token

}
func findValueInMap(m map[string]ChatUser, targetValue string) (string, bool) {
	for key, value := range m {
		if value.InternData.Name == targetValue {
			return key, true // Return the key and a boolean indicating the value was found
		}
	}
	return "", false // Return an empty string and false if the value is not found
}
