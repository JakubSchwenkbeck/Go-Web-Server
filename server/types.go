package main

import "time"

// Message represents a chat message exchanged between users.
// It contains information about the sender, receiver, message content, and timestamp.
type Message struct {
	SenderID   string    // ID of the user sending the message
	ReceiverID string    // ID of the user receiving the message
	Message    string    // The content of the message
	TimeStamp  time.Time // The time when the message was sent
}

// ChatUser represents a user in the chat system with associated internal and security details.
// It includes the user's internal data, password, and hashed password for authentication and security purposes.
type ChatUser struct {
	InternData   User   // Internal user data such as ID and name
	Password     string // Password for user authentication (should be securely hashed in real scenarios)
	HashPassword string // Hashed password for secure comparison and authentication (useful for hierarchical systems)
}

// User represents a basic user structure with an ID and name.
// It is used for identifying users and can be serialized to/from JSON format.
type User struct {
	ID   string `json:"id"`   // Unique identifier for the user
	Name string `json:"name"` // Name of the user
}
