package models

import "sync"

type Hub struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan *Message

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	Mu sync.Mutex
}

type RoomHubs struct {
	Hubs map[int64]*Hub
	Mu   *sync.Mutex
}

func NewRoomHubs() *RoomHubs {
	return &RoomHubs{
		Hubs: make(map[int64]*Hub),
		Mu:   &sync.Mutex{},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.Unregister:
			h.Mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.Mu.Unlock()
		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					h.Mu.Lock()
					close(client.Send)
					delete(h.Clients, client)
					h.Mu.Unlock()
				}
			}
		}
	}
}
