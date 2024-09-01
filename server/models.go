package main

import (
	"fmt"
	"time"
)

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

// CreateMessage initializes a new Message with the provided details.
// It sets the TimeStamp field to the current time.
func CreateMessage(senderID, receiverID, messageContent string) Message {
	return Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    messageContent,
		TimeStamp:  time.Now(), // Set the timestamp to the current time
	}
}

// DisplayMessage prints the details of a Message to the console in a readable format.
func (m Message) DisplayMessage() {
	fmt.Printf("From: %s\nTo: %s\nMessage: %s\nSent at: %s\n",
		m.SenderID, m.ReceiverID, m.Message, m.TimeStamp.Format(time.RFC1123))
}

// CreateChatUser initializes a new ChatUser with the given user data and passwords.
// The hashed password should be provided after hashing the plain password.
func CreateChatUser(user User, password, hashedPassword string) ChatUser {
	return ChatUser{
		InternData:   user,
		Password:     password,
		HashPassword: hashedPassword,
	}
}

// DisplayUser prints the details of a ChatUser to the console.
func (u ChatUser) DisplayUser() {
	fmt.Printf("User ID: %s\nUser Name: %s\n", u.InternData.ID, u.InternData.Name)
}
