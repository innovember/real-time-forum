package models

type Room struct {
	ID              int64    `json:"id"`
	User            *User    `json:"user"`
	LastMessageDate int64    `json:"lastMessageDate"`
	Read            bool     `json:"read"`
	UnreadMsgNumber int64    `json:"unreadMsgNumber"`
	LastMessage     *Message `json:"lastMessage"`
}

type Message struct {
	ID            int64  `json:"id"`
	RoomID        int64  `json:"roomID"`
	Content       string `json:"content"`
	MessageDate   int64  `json:"messageDate"`
	Read          bool   `json:"read"`
	User          *User  `json:"user"`
	IsYourMessage bool   `json:"isYourMessage"`
}
