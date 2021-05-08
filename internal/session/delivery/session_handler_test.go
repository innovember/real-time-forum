package delivery_test

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/innovember/real-time-forum/internal/models"
	sessionDelivery "github.com/innovember/real-time-forum/internal/session/delivery"
	sessionRepo "github.com/innovember/real-time-forum/internal/session/repository"
	sessionUsecase "github.com/innovember/real-time-forum/internal/session/usecases"
	userRepo "github.com/innovember/real-time-forum/internal/user/repository"
	userUsecase "github.com/innovember/real-time-forum/internal/user/usecases"
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
func TestHandlerLogin(t *testing.T) {
	var (
		input models.InputUserSignIn
	)
	input = models.InputUserSignIn{
		Nickname: "newuser",
		Password: "qweasd123",
	}
	body, _ := json.Marshal(input)

	userRepository := userRepo.NewUserDBRepository(dbConn)
	sessionRepository := sessionRepo.NewSessionDBRepository(dbConn)

	userUsecase := userUsecase.NewUserUsecase(userRepository)
	sessionUsecase := sessionUsecase.NewSessionUsecase(sessionRepository)
	sessionHander := sessionDelivery.NewSessionHandler(sessionUsecase, userUsecase)

	handler := http.HandlerFunc(sessionHander.HandlerLogin)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/session/login",
		strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(rec, req)
	if status := rec.Code; status != http.StatusOK {
		bytes, _ := ioutil.ReadAll(rec.Body)
		t.Error(string(bytes))
		t.Error("login fail")
	}
}
