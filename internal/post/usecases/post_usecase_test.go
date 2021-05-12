package usecases_test

import (
	"database/sql"
	"log"
	"testing"

	categoryRepo "github.com/innovember/real-time-forum/internal/category/repository"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/post/repository"
	"github.com/innovember/real-time-forum/internal/post/usecases"
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

func TestCreatePost(t *testing.T) {
	dbConn := setup()
	post := models.Post{
		AuthorID: 1,
		Title:    "smth new",
		Content:  "some random text",
	}
	categories := []string{"new", "random"}
	userRepo := userRepo.NewUserDBRepository(dbConn)
	categoryRepository := categoryRepo.NewCategoryDBRepository(dbConn)
	postRepo := repository.NewPostDBRepository(dbConn, userRepo)
	postUCase := usecases.NewPostUsecase(postRepo, categoryRepository)
	if err := postUCase.Create(&post, categories); err != nil {
		t.Error(err)
	}
}
