package chat

import "github.com/innovember/real-time-forum/internal/models"

type RoomUsecase interface {
	CreateRoom(userID1, userID2 int) (*models.Room, error)
	GetRoomByUsers(userID1, userID2 int) (*models.Room, error)
	GetUsersByRoom(roomID int64) ([]models.User, error)
	GetAllRoomsByUserID(userID int64) ([]models.Room, error)
	DeleteRoom(id int) error
	CreateMessage(roomID int64, msg *models.Message) error
	GetMessages(roomID int64, lastMessageID int64) ([]models.Message, error)
	GetLastMessageDate(roomID int64) (int64, error)
}
