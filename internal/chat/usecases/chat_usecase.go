package usecases

import (
	"github.com/innovember/real-time-forum/internal/chat"
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

func (ru *RoomUsecase) CreateRoom(userID1, userID2 int) (*models.Room, error) {
	room, err := ru.roomRepo.InsertRoom(userID1, userID2)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (ru *RoomUsecase) GetRoomByUsers(userID1, userID2 int) (*models.Room, error) {
	room, err := ru.roomRepo.SelectRoomByUsers(userID1, userID2)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (ru *RoomUsecase) GetUsersByRoom(roomID int) ([]models.User, error) {
	users, err := ru.roomRepo.SelectUsersByRoom(roomID)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ru *RoomUsecase) GetAllRoomsByUserID(userID int64) ([]models.Room, error) {
	rooms, err := ru.roomRepo.SelectAllRoomsByUserID(userID)
	if err != nil {
		return nil, err
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

func (ru *RoomUsecase) CreateMessage(roomID int, msg *models.Message) error {
	err := ru.roomRepo.InsertMessage(roomID, msg)
	if err != nil {
		return err
	}
	return nil
}

func (ru *RoomUsecase) GetMessages(roomID int, lastMessageID int64) ([]models.Message, error) {
	messages, err := ru.roomRepo.SelectMessages(roomID, lastMessageID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (ru *RoomUsecase) GetLastMessageDate(roomID int) (int64, error) {
	lastMessageDate, err := ru.roomRepo.SelectLastMessageDate(roomID)
	if err != nil {
		return 0, err
	}
	return lastMessageDate, nil
}
