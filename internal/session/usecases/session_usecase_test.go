package usecases_test

import (
	"database/sql"
	"log"
	"testing"

	sessionRepo "github.com/innovember/real-time-forum/internal/session/repository"
	"github.com/innovember/real-time-forum/internal/session/usecases"
	userRepo "github.com/innovember/real-time-forum/internal/user/repository"
	userUsecase "github.com/innovember/real-time-forum/internal/user/usecases"
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

func TestGetByToken(t *testing.T) {
	dbConn := setup()
	token := "f53601f1-79d4-4f54-a444-d316b6ad8200"
	userRepository := userRepo.NewUserDBRepository(dbConn)
	sessionRepository := sessionRepo.NewSessionDBRepository(dbConn)

	userUsecase := userUsecase.NewUserUsecase(userRepository)
	sessionUsecase := usecases.NewSessionUsecase(sessionRepository)
	session, err := sessionUsecase.GetByToken(token)
	if err != nil {
		t.Error("session err ", err)
	}
	user, err := userUsecase.GetByID(session.UserID)
	if err != nil {
		t.Error("user err ", err)
	}
	if user == nil {
		t.Error("user not found")
	}
}
