package models

import (
	"database/sql"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	ID              int64 `json:"id"`
	User            *User `json:"user"`
	LastMessageDate int64 `json:"lastMessageDate"`
}

type Message struct {
	ID            int64  `json:"id"`
	RoomID        int64  `json:"roomID"`
	Content       string `json:"content"`
	MessageDate   int64  `json:"messageDate"`
	User          *User  `json:"user"`
	IsYourMessage bool   `json:"isYourMessage"`
}

type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan []*Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	MU sync.Mutex
}

type Client struct {
	UserID int64

	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []*Message

	DB *sql.DB

	MU sync.Mutex
}

