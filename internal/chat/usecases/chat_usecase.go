package usecases

import (
	"github.com/innovember/real-time-forum/internal/chat"
	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/models"
)

type RoomUsecase struct {
	roomRepo chat.RoomRepository
}

func NewRoomUsecase(roomRepo chat.RoomRepository) *RoomUsecase {
	return &RoomUsecase{
		roomRepo: roomRepo,
	}
}

func (ru *RoomUsecase) CreateRoom(userID1, userID2 int64) (*models.Room, error) {
	room, err := ru.roomRepo.InsertRoom(userID1, userID2)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (ru *RoomUsecase) GetRoomByUsers(userID1, userID2 int64) (int64, error) {
	room, err := ru.roomRepo.SelectRoomByUsers(userID1, userID2)
	if err != nil {
		return 0, err
	}
	return room, nil
}

func (ru *RoomUsecase) GetUsersByRoom(roomID int64) ([]models.User, error) {
	users, err := ru.roomRepo.SelectUsersByRoom(roomID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ru *RoomUsecase) GetAllRoomsByUserID(userID int64) ([]models.Room, error) {
	var rooms []models.Room
	users, err := ru.roomRepo.SelectAllUsers(userID)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		var room models.Room
		room.User = user
		room.ID, err = ru.roomRepo.SelectRoomByUsers(userID, user.ID)
		if err != nil && err != consts.ErrNoData {
			return nil, err
		}
		room.LastMessageDate, err = ru.roomRepo.SelectLastMessageDate(room.ID)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}

func (ru *RoomUsecase) DeleteRoom(id int) error {
	err := ru.roomRepo.DeleteRoom(id)
	if err != nil {
		return err
	}
	return nil
}

func (ru *RoomUsecase) CreateMessage(roomID int64, msg *models.Message) error {
	err := ru.roomRepo.InsertMessage(roomID, msg)
	if err != nil {
		return err
	}
	return nil
}

func (ru *RoomUsecase) GetMessages(roomID int64, lastMessageID int64) ([]models.Message, error) {
	messages, err := ru.roomRepo.SelectMessages(roomID, lastMessageID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (ru *RoomUsecase) GetLastMessageDate(roomID int64) (int64, error) {
	lastMessageDate, err := ru.roomRepo.SelectLastMessageDate(roomID)
	if err != nil {
		return 0, err
	}
	return lastMessageDate, nil
}
