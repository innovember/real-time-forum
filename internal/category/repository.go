package category

import (
	"database/sql"

	"github.com/innovember/real-time-forum/internal/models"
)

type CategoryRepository interface {
	Insert(postID int64, categories []string) (err error)
	SelectAllCategories() (categories []models.Category, err error)
	SelectByName(name string, tx *sql.Tx) (id int64, err error)
	IsCategoryExist(category string, tx *sql.Tx) (bool, error)
}
