package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

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
func (cs *ChatService) RegisterUser(id, name, password string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.users[id] = ChatUser{
		InternData: User{
			ID:   id,
			Name: name,
		},

		Password: password,
	}
	fmt.Print("Successfully Registered " + name + " \n")
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
	fmt.Printf("Message from %s to %s: %s\n", cs.users[senderID].InternData.Name, cs.users[receiverID].InternData.Name, message)
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

func ChatAppMain(r mux.Router, port string) {
	chatService := NewChatService()
	// Register routes for the chat functionality
	r.HandleFunc("/messages", chatService.SendMessageHandler).Methods("POST")
	r.HandleFunc("/messages/{id}", chatService.GetMessagesHandler).Methods("GET")

	// Example: Register some users
	chatService.RegisterUser("1", "Jakub", "password123")
	chatService.RegisterUser("2", "Marie", "password456")

	chatService.SendMessage("1", "2", "Hello there!")

	//log.Fatal(http.ListenAndServe(":"+port, r))
}
