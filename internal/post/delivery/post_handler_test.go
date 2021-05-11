package delivery_test

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestCreate(t *testing.T) {
	dbConn := setup()
	post := models.Post{
		AuthorID: 1,
		Title:    "smth new",
		Content:  "some random text",
	}
	categories := []string{"new", "rand"}
	userRepo := userRepo.NewUserDBRepository(dbConn)
	postRepo := repository.NewPostDBRepository(dbConn, userRepo)
	postUCase := usecases.NewPostUsecase(postRepo, categoryRepo)
	body, err := json.Marshal(post)
	if err != nil {
		t.Error(err)
	}
	handler := http.HandlerFunc(delivery.HandlerCreatePost)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/post", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if status := w.Code; status != http.StatusCreated {
		t.Error("didnt create post")
	}
}
