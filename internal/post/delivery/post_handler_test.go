package delivery_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	categoryRepo "github.com/innovember/real-time-forum/internal/category/repository"
	"github.com/innovember/real-time-forum/internal/consts"
	"github.com/innovember/real-time-forum/internal/models"
	"github.com/innovember/real-time-forum/internal/post/delivery"
	"github.com/innovember/real-time-forum/internal/post/repository"
	"github.com/innovember/real-time-forum/internal/post/usecases"
	userRepo "github.com/innovember/real-time-forum/internal/user/repository"
	userUcase "github.com/innovember/real-time-forum/internal/user/usecases"
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
	post := models.InputPost{
		Title:      "smth new",
		Content:    "some random text",
		Categories: []string{"new", "random"},
	}
	userRepo := userRepo.NewUserDBRepository(dbConn)
	userUcase := userUcase.NewUserUsecase(userRepo)
	categoryRepository := categoryRepo.NewCategoryDBRepository(dbConn)
	postRepo := repository.NewPostDBRepository(dbConn, userRepo)
	postUCase := usecases.NewPostUsecase(postRepo, categoryRepository)
	delivery := delivery.NewPostHandler(postUCase, userUcase)
	body, err := json.Marshal(post)
	if err != nil {
		t.Error(err)
	}
	handler := http.HandlerFunc(delivery.HandlerCreatePost)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/post", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), consts.ConstAuthedUserParam, 1)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req.WithContext(ctx))
	if status := w.Code; status != http.StatusCreated {
		t.Error("didnt create post")
	}
}
