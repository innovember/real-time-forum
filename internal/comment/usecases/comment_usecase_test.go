package usecases_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/innovember/real-time-forum/internal/comment/repository"
	"github.com/innovember/real-time-forum/internal/comment/usecases"
	"github.com/innovember/real-time-forum/internal/models"
	userRepo "github.com/innovember/real-time-forum/internal/user/repository"
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

func TestCreateComment(t *testing.T) {
	dbConn := setup()
	comment := models.Comment{
		AuthorID: 1,
		PostID:   1,
		Content:  "some random text",
	}
	userRepo := userRepo.NewUserDBRepository(dbConn)
	commentRepo := repository.NewCommentDBRepository(dbConn, userRepo)
	commentUcase := usecases.NewCommentUsecase(commentRepo)
	if err := commentUcase.Create(&comment); err != nil {
		t.Error(err)
	}
}
