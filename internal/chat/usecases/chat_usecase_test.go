package usecases_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	repo "github.com/innovember/real-time-forum/internal/chat/repository"
	usecase "github.com/innovember/real-time-forum/internal/chat/usecases"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/pkg/database"
)

var (
	dbConn *sql.DB
	err    error
)

func setup() *sql.DB {
	dbConn, err = database.GetDBInstance("sqlite3", "../../../database/forum.db")
	if err != nil {
		log.Fatal("dbConn err: ", err)
	}
	return dbConn
}

func setupChatUsecases() (*usecase.RoomUsecase, *usecase.HubUsecase) {
	dbConn := setup()
	hubs := models.NewRoomHubs()
	roomRepository := repo.NewRoomRepository(dbConn)
	hubRepository := repo.NewHubRepository(hubs)
	roomUsecase := usecase.NewRoomUsecase(roomRepository)
	hubUsecase := usecase.NewHubUsecase(hubRepository, roomRepository)
	return roomUsecase, hubUsecase
}
func TestCreateRoom(t *testing.T) {
	roomUsecase, _ := setupChatUsecases()
	room, err := roomUsecase.CreateRoom(1, 2)
	if err != nil {
		t.Error("insert room err ", err)
	}
	fmt.Println(room)
}

func TestGetRoomByUsers(t *testing.T) {
	roomUsecase, _ := setupChatUsecases()
	room, err := roomUsecase.GetRoomByUsers(1, 3)
	if err != nil {
		t.Error("get room err ", err)
	}
	fmt.Println(room)
}

func TestGetUsersByRoom(t *testing.T) {
	roomUsecase, _ := setupChatUsecases()
	users, err := roomUsecase.GetUsersByRoom(1)
	if err != nil {
		t.Error("get users by room err: ", err)
	}
	fmt.Println(users)
}

func TestGetAllRoomsByUserID(t *testing.T) {
	roomUsecase, _ := setupChatUsecases()
	rooms, err := roomUsecase.GetAllRoomsByUserID(1)
	if err != nil {
		t.Error("get all rooms by userID err: ", err)
	}
	fmt.Println(rooms)
}

func TestCreateMessage(t *testing.T) {
	roomUsecase, _ := setupChatUsecases()
	msg := &models.Message{
		RoomID:      1,
		Content:     "hi",
		MessageDate: 1630430430430,
		User:        &models.User{ID: 1},
	}
	err := roomUsecase.CreateMessage(msg)
	if err != nil {
		t.Error("create msg err: ", err)
	}
}

func TestGetMessages(t *testing.T) {
	roomUsecase, _ := setupChatUsecases()
	messages, err := roomUsecase.GetMessages(2, 0)
	if err != nil {
		t.Error("portion msg err: ", err)
	}
	fmt.Println(messages)
}

func TestGetLastMessageDate(t *testing.T) {
	roomUsecase, _ := setupChatUsecases()
	lastMsgDate, err := roomUsecase.GetLastMessageDate(2)
	if err != nil {
		t.Error("msgDate err: ", err)
	}
	fmt.Println(lastMsgDate)
}

func TestGetAllUsers(t *testing.T) {
	roomUsecase, _ := setupChatUsecases()
	users, err := roomUsecase.GetAllUsers(1)
	if err != nil {
		t.Error("all users err: ", err)
	}
	fmt.Println(users)
}

func TestGetRoomByID(t *testing.T) {
	roomUsecase, _ := setupChatUsecases()
	room, err := roomUsecase.GetRoomByID(1)
	if err != nil {
		t.Error("get room by ID err: ", err)
	}
	fmt.Println(room)
}
