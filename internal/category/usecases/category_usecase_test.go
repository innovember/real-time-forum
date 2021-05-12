package usecases_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/innovember/real-time-forum/internal/category/repository"
	"github.com/innovember/real-time-forum/internal/category/usecases"
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

func TestGetAllCategories(t *testing.T) {
	dbConn := setup()
	categoryRepo := repository.NewCategoryDBRepository(dbConn)
	categoryUcase := usecases.NewCategoryUsecase(categoryRepo)
	categories, err := categoryUcase.GetAllCategories()
	fmt.Println(categories)
	if err != nil {
		t.Error("all categories err ", err)
	}
}
