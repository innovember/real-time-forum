package usecases_test

import (
	"database/sql"
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

func init() {
	dbConn, err = database.GetDBInstance("sqlite3", "../../../database/forum.db")
	if err != nil {
		log.Fatal("dbConn err: ", err)
	}
}

func TestCreateUseCaseReturnsPass(t *testing.T) {
	defer dbConn.Close()
	user := models.User{
		Nickname:  "erha",
		Email:     "erha@gmail.com",
		Password:  "reactdev",
		FirstName: "Erha",
		LastName:  "D2",
		Age:       20,
		Gender:    "male",
	}
	userRepo := repository.NewUserDBRepository(dbConn)
	userUCase := usecases.NewUserUsecase(userRepo)
	if err := userUCase.Create(&user); err != nil {
		t.Error(err)
	}
}
