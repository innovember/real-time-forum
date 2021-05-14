package usecases_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	categoryRepo "github.com/innovember/real-time-forum/internal/category/repository"
	commentRepo "github.com/innovember/real-time-forum/internal/comment/repository"
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
	commentRepo := commentRepo.NewCommentDBRepository(dbConn, userRepo)
	postRepo := repository.NewPostDBRepository(dbConn, userRepo, commentRepo)
	postUCase := usecases.NewPostUsecase(postRepo, categoryRepository)
	if err := postUCase.Create(&post, categories); err != nil {
		t.Error(err)
	}
}

func TestGetPost(t *testing.T) {
	dbConn := setup()
	userRepo := userRepo.NewUserDBRepository(dbConn)
	commentRepo := commentRepo.NewCommentDBRepository(dbConn, userRepo)
	categoryRepository := categoryRepo.NewCategoryDBRepository(dbConn)
	postRepo := repository.NewPostDBRepository(dbConn, userRepo, commentRepo)
	postUCase := usecases.NewPostUsecase(postRepo, categoryRepository)
	post, err := postUCase.GetPostByID(1)
	fmt.Printf("%+v\n", post)
	if err != nil {
		t.Error(err)
	}
}

func TestGetPosts(t *testing.T) {
	dbConn := setup()
	userRepo := userRepo.NewUserDBRepository(dbConn)
	commentRepo := commentRepo.NewCommentDBRepository(dbConn, userRepo)
	categoryRepository := categoryRepo.NewCategoryDBRepository(dbConn)
	postRepo := repository.NewPostDBRepository(dbConn, userRepo, commentRepo)
	postUCase := usecases.NewPostUsecase(postRepo, categoryRepository)
	// for i := 0; i < 30; i++ {
	// 	post := models.Post{
	// 		AuthorID: 1,
	// 		Title:    fmt.Sprintf("smth new title #%d", i),
	// 		Content:  fmt.Sprintf("some random text #%d", i),
	// 	}
	// 	categories := []string{"new", "random", "smth"}
	// 	time.Sleep(1 * time.Second)
	// 	if err := postUCase.Create(&post, categories); err != nil {
	// 		t.Error(err)
	// 	}
	// }
	posts, err := postUCase.GetAllPosts(&models.InputGetPosts{
		Option:     "all",
		LastPostID: 24,
		Limit:      5,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, post := range posts {
		fmt.Printf("%+v\n", post)
		fmt.Println("-------------------------------------------------------")
	}
	if err != nil {
		t.Error(err)
	}
}
