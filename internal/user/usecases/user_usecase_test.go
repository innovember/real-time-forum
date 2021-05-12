package usecases_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/user/repository"
	"github.com/innovember/real-time-forum/internal/user/usecases"
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

func TestCreateUser(t *testing.T) {
	dbConn := setup()
	user := models.User{
		Nickname:  "erha",
		Email:     "erha@gmail.com",
		Password:  "reactdev",
		FirstName: "Erha",
		LastName:  "D2",
		Age:       20,
		Gender:    "male",
		Status:    "offline",
	}
	userRepo := repository.NewUserDBRepository(dbConn)
	userUCase := usecases.NewUserUsecase(userRepo)
	if err := userUCase.Create(&user); err != nil {
		t.Error(err)
	}
}

func TestCheckPassword(t *testing.T) {
	dbConn := setup()
	user := models.InputUserSignIn{
		Nickname: "erha",
		Password: "reactdev",
	}
	userRepo := repository.NewUserDBRepository(dbConn)
	userUCase := usecases.NewUserUsecase(userRepo)
	if err := userUCase.CheckPassword(&user); err != nil {
		t.Error(err)
	}
}

func TestGetUserByEmailOrNickname(t *testing.T) {
	dbConn := setup()
	user := models.InputUserSignIn{
		Nickname: "erha",
		Password: "reactdev",
	}
	userRepo := repository.NewUserDBRepository(dbConn)
	userUCase := usecases.NewUserUsecase(userRepo)
	u, err := userUCase.GetByEmailOrNickname(user.Nickname)
	fmt.Println(u)
	if err != nil {
		t.Error(err)
	}
}
