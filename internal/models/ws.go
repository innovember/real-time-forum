package models

type WsResponse struct {
	RoomID   int64  `json:"roomID"`
	Content  string `json:"content"`
	HTTPCode int    `json:"httpCode,omitempty"`
	State    *bool  `json:"state,omitempty"`
}

var WsEventType = struct {
	Message          string
	WsError          string
	WsClosed         string
	PingMessage      string
	PongMessage      string
	TypingInReceiver string
	TypingInSender   string
}{
	Message:          "Message",
	WsError:          "WsError",
	WsClosed:         "WsClosed",
	PingMessage:      "PingMessage",
	PongMessage:      "PongMessage",
	TypingInReceiver: "TypingInReceiver",
	TypingInSender:   "TypingInSender",
}

type WsEvent struct {
	EventType          string   `json:"eventType,omitempty"`
	Content            string   `json:"content,omitempty"`
	RecipientID        int64    `json:"recipientID,omitempty"`
	TypingInReceiverID int64    `json:"typingInReceiverID,omitempty"`
	TypingInUserID     int64    `json:"typingInUserID,omitempty"`
	RoomID             int64    `json:"roomID,omitempty"`
	Message            *Message `json:"message,omitempty"`
	State              bool     `json:"state"`
}
