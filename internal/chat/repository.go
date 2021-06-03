package chat

import (
	"github.com/innovember/real-time-forum/internal/models"
)

type RoomRepository interface {
	InsertRoom(userID1, userID2 int64) (*models.Room, error)
	SelectRoomByUsers(userID1, userID2 int64) (int64, error)
	SelectUsersByRoom(roomID int64) ([]models.User, error)
	SelectAllUsers(userID int64) ([]*models.User, error)
	DeleteRoom(id int64) error
	InsertMessage(roomID int64, msg *models.Message) error
	SelectMessages(roomID int64, lastMessageID int64) ([]models.Message, error)
	SelectLastMessageDate(roomID int64) (int64, error)
	SelectRoomByID(roomID int64) (*models.Room, error)
}

type HubRepository interface {
	NewHub() *models.Hub
	GetHub(roomID int64) (*models.Hub, bool)
	DeleteHub(roomID int64)
	Register(roomID int64, hub *models.Hub)
}
