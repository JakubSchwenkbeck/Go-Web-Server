package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Message structure to represent a chat message
type Message struct {
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Message    string    `json:"message"`
	TimeStamp  time.Time `json:"timestamp"`
}

// ChatUser represents a user in the chat application
type ChatUser struct {
	InternData  User   // Internal user data
	DisplayName string `json:"display_name"`
	Password    string `json:"-"`
}

// User represents internal user data
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ChatService handles user management and message sending
type ChatService struct {
	users    map[string]ChatUser // Map of user IDs to users
	messages []Message           // Slice to store messages
	mu       sync.Mutex          // Mutex to handle concurrent access
}

// NewChatService creates a new ChatService
func NewChatService() *ChatService {
	return &ChatService{
		users:    make(map[string]ChatUser),
		messages: []Message{},
	}
}

// RegisterUser registers a new user in the chat application
func (cs *ChatService) RegisterUser(id, name, displayName, password string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.users[id] = ChatUser{
		InternData: User{
			ID:   id,
			Name: name,
		},
		DisplayName: displayName,
		Password:    password,
	}
}

// SendMessage sends a message from one user to another
func (cs *ChatService) SendMessage(senderID, receiverID, message string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, ok := cs.users[senderID]; !ok {
		fmt.Println("Sender not found")
		return
	}
	if _, ok := cs.users[receiverID]; !ok {
		fmt.Println("Receiver not found")
		return
	}

	msg := Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    message,
		TimeStamp:  time.Now(),
	}

	cs.messages = append(cs.messages, msg)
	fmt.Printf("Message from %s to %s: %s\n", cs.users[senderID].DisplayName, cs.users[receiverID].DisplayName, message)
}

// GetMessagesForUser retrieves all messages sent to a specific user
func (cs *ChatService) GetMessagesForUser(userID string) []Message {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	var userMessages []Message
	for _, msg := range cs.messages {
		if msg.ReceiverID == userID {
			userMessages = append(userMessages, msg)
		}
	}
	return userMessages
}

// HTTP Handlers

func (cs *ChatService) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	cs.SendMessage(msg.SenderID, msg.ReceiverID, msg.Message)
	w.WriteHeader(http.StatusCreated)
}

func (cs *ChatService) GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	messages := cs.GetMessagesForUser(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func main() {
	chatService := NewChatService()
	r := mux.NewRouter()

	// Register routes for the chat functionality
	r.HandleFunc("/messages", chatService.SendMessageHandler).Methods("POST")
	r.HandleFunc("/messages/{id}", chatService.GetMessagesHandler).Methods("GET")

	// Example: Register some users
	chatService.RegisterUser("1", "Jakub", "Jakub123", "password123")
	chatService.RegisterUser("2", "Alice", "Alice456", "password456")

	port := "8080"
	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
