package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	UserID int64
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan *WsEvent
	Mu     sync.Mutex
}
