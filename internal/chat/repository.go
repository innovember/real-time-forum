package chat

import (
	"github.com/innovember/real-time-forum/internal/models"
)

type RoomRepository interface {
	InsertRoom(userID1, userID2 int) (*models.Room, error)
	SelectRoomByUsers(userID1, userID2 int) (*models.Room, error)
	SelectUsersByRoom(roomID int) ([]models.User, error)
	SelectAllRoomsByUserID(userID int64) ([]models.Room, error)
	DeleteRoom(id int) error
	InsertMessage(roomID int, msg *models.Message) error
	SelectMessages(roomID int, lastMessageID int64) ([]models.Message, error)
	SelectLastMessageDate(roomID int) (int64, error)
}
