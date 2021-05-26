package chat

import "github.com/innovember/real-time-forum/internal/models"

type RoomUsecase interface {
	CreateRoom(userID1, userID2 int) (*models.Room, error)
	GetRoomByUsers(userID1, userID2 int) (*models.Room, error)
	GetUsersByRoom(roomID int) ([]models.User, error)
	GetAllRoomsByUserID(userID int64) ([]models.Room, error)
	DeleteRoom(id int) error
	CreateMessage(roomID int, msg *models.Message) error
	GetMessages(roomID int, lastMessageID int64) ([]models.Message, error)
	GetLastMessageDate(roomID int) (int64, error)
}
