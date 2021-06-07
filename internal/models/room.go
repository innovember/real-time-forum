package models

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
	HTTPCode      int    `json:"httpCode"`
	State         bool   `json:"state"`
}
