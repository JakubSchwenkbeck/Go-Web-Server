package main

import "time"

type Message struct {
	SenderID   string
	ReceiverID string // Messages will be send over the IDs

	Message   string
	TimeStamp time.Time
}
type ChatUser struct {
	InternData User // intern (maybe hidden?) ID and intern name

	Password string //

	HashPassword string // maybe for hierachy useful
}

// Defining JSON User strcuture with id and name
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
