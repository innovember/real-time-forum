package repository

import (
	"database/sql"

	"github.com/innovember/real-time-forum/internal/category"
)

type CategoryDBRepository struct {
	dbConn *sql.DB
}

func NewCategoryDBRepository(conn *sql.DB) category.CategoryRepository {
	return &CategoryDBRepository{dbConn: conn}
}
