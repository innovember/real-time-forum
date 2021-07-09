package usecases

import (
	"sort"

	"github.com/innovember/real-time-forum/internal/chat"
	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/helpers"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/user"
)

type RoomUsecase struct {
	roomRepo chat.RoomRepository
	userRepo user.UserRepository
}

func NewRoomUsecase(roomRepo chat.RoomRepository,
	userRepo user.UserRepository) *RoomUsecase {
	return &RoomUsecase{
		roomRepo: roomRepo,
		userRepo: userRepo,
	}
}

func (ru *RoomUsecase) CreateRoom(userID1, userID2 int64) (*models.Room, error) {
	room, err := ru.roomRepo.InsertRoom(userID1, userID2)
	if err != nil {
		return nil, err
	}
	room.Read = true
	return room, nil
}

func (ru *RoomUsecase) GetRoomByUsers(userID1, userID2 int64) (*models.Room, error) {
	roomID, err := ru.roomRepo.SelectRoomByUsers(userID1, userID2)
	if err != nil {
		return nil, err
	}
	room, err := ru.roomRepo.SelectRoomByID(roomID)
	if err != nil {
		return nil, err
	}
	user, err := ru.userRepo.SelectByID(userID2)
	if err != nil {
		return nil, err
	}
	room.User = user
	room.UnreadMsgNumber, err = ru.GetUnReadMessages(room.ID)
	if err != nil && err != consts.ErrNoData {
		return nil, err
	}
	if room.UnreadMsgNumber != 0 {
		room.Read = false
	} else {
		room.Read = true
	}
	room.LastMessageDate, err = ru.GetLastMessageDate(room.ID)
	if err != nil && err != consts.ErrNoData {
		return nil, err
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
		if err != nil {
			switch err {
			case consts.ErrNoData:
				continue
			default:
				return nil, err
			}
		}
		room.LastMessageDate, err = ru.GetLastMessageDate(room.ID)
		if err != nil && err != consts.ErrNoData {
			return nil, err
		}
		room.UnreadMsgNumber, err = ru.GetUnReadMessages(room.ID)
		if err != nil && err != consts.ErrNoData {
			return nil, err
		}
		if room.UnreadMsgNumber != 0 {
			room.Read = false
		} else {
			room.Read = true
		}
		rooms = append(rooms, room)
	}
	sort.SliceStable(rooms, func(i, j int) bool {
		return rooms[i].LastMessageDate > rooms[j].LastMessageDate
	})
	return rooms, nil
}

func (ru *RoomUsecase) DeleteRoom(id int64) error {
	err := ru.roomRepo.DeleteRoom(id)
	if err != nil {
		return err
	}
	return nil
}

func (ru *RoomUsecase) CreateMessage(msg *models.Message) (*models.Message, error) {
	message, err := ru.roomRepo.InsertMessage(msg)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (ru *RoomUsecase) GetMessages(roomID, lastMessageID, userID int64) ([]models.Message, error) {
	messages, err := ru.roomRepo.SelectMessages(roomID, lastMessageID, userID)
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

func (ru *RoomUsecase) GetAllUsers(userID int64) ([]*models.User, error) {
	users, err := ru.roomRepo.SelectAllUsers(userID)
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (ru *RoomUsecase) GetRoomByID(roomID int64) (*models.Room, error) {
	room, err := ru.roomRepo.SelectRoomByID(roomID)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (ru *RoomUsecase) GetUnReadMessages(roomID int64) (int64, error) {
	unreadMsgNumber, err := ru.roomRepo.SelectUnReadMessages(roomID)
	if err != nil {
		return 0, err
	}
	return unreadMsgNumber, nil
}

func (ru *RoomUsecase) UpdateMessageStatus(roomID, messageID int64) error {
	err := ru.roomRepo.UpdateMessageStatus(roomID, messageID)
	if err != nil {
		return err
	}
	return nil
}

func (ru *RoomUsecase) UpdateMessagesStatusForReceiver(roomID, userID int64) error {
	users, err := ru.GetUsersByRoom(roomID)
	if err != nil {
		return err
	}
	authorID := helpers.SelectSecondUser(users, userID)
	err = ru.roomRepo.UpdateMessagesStatusForReceiver(roomID, authorID)
	if err != nil {
		return err
	}
	return nil
}
