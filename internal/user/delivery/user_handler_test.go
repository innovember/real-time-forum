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
	"github.com/innovember/real-time-forum/internal/user/delivery"
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
func TestRegisterHandler(t *testing.T) {
	user := models.User{
		Nickname:  "newuser",
		Email:     "newuser@gmail.com",
		Password:  "qweasd123",
		FirstName: "Erha",
		LastName:  "D2",
		Age:       20,
		Gender:    "male",
	}
	body, err := json.Marshal(user)
	if err != nil {
		t.Error(err)
	}
	userRepo := repository.NewUserDBRepository(dbConn)
	userUCase := usecases.NewUserUsecase(userRepo)
	delivery := delivery.NewUserHandler(userUCase)
	handler := http.HandlerFunc(delivery.HandlerRegisterUser)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/user/signup", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if status := w.Code; status != http.StatusCreated {
		t.Error("didnt create user")
	}
}
