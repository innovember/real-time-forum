package repository

import (
	"github.com/innovember/real-time-forum/internal/chat"
	"github.com/innovember/real-time-forum/internal/models"
)

type HubRepository struct {
	roomHubs *models.RoomHubs
}

func NewHubRepository(roomHubs *models.RoomHubs) chat.HubRepository {
	return &HubRepository{
		roomHubs: roomHubs,
	}
}

func (hr *HubRepository) NewHub() *models.Hub {
	return &models.Hub{
		Broadcast:  make(chan *models.WsEvent),
		Register:   make(chan *models.Client),
		Unregister: make(chan *models.Client),
		Clients:    make(map[*models.Client]bool),
	}
}

func (hr *HubRepository) GetHub(roomID int64) (*models.Hub, bool) {
	hub, ok := hr.roomHubs.Hubs[roomID]
	if !ok {
		return nil, false
	}
	return hub, true
}

func (hr *HubRepository) DeleteHub(roomID int64) {
	hr.roomHubs.Mu.Lock()
	delete(hr.roomHubs.Hubs, roomID)
	hr.roomHubs.Mu.Unlock()
}

func (hr *HubRepository) Register(roomID int64, hub *models.Hub) {
	hr.roomHubs.Mu.Lock()
	_, ok := hr.roomHubs.Hubs[roomID]
	if !ok {
		hr.roomHubs.Hubs[roomID] = hub
	}
	hr.roomHubs.Mu.Unlock()
}
