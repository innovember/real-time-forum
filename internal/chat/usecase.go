package chat

import "github.com/innovember/real-time-forum/internal/models"

type RoomUsecase interface {
	CreateRoom(userID1, userID2 int64) (*models.Room, error)
	GetRoomByUsers(userID1, userID2 int64) (*models.Room, error)
	GetUsersByRoom(roomID int64) ([]models.User, error)
	GetAllRoomsByUserID(userID int64) ([]models.Room, error)
	DeleteRoom(id int64) error
	CreateMessage(roomID int64, msg *models.Message) error
	GetMessages(roomID int64, lastMessageID int64) ([]models.Message, error)
	GetLastMessageDate(roomID int64) (int64, error)
	GetAllUsers(userID int64) ([]*models.User, error)
	GetRoomByID(roomID int64) (*models.Room, error)
}

type Hub interface {
	NewHub() *models.Hub
	GetHub(roomID int64) (*models.Hub, error)
	DeleteHub(roomID int64)
	Register(roomID int64, hub *models.Hub)
	ServeWS()
}
