package main

type Message struct {
	SenderID   string
	ReceiverID string // Messages will be send over the IDs

	Message string

	TimeStamp string
}
type ChatUser struct {
	InternData User // intern (maybe hidden?) ID and intern name

	DisplayName string // maybe Displayname

	role string // maybe for hierachy useful
}
