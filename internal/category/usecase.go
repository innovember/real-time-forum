package category

import "github.com/innovember/real-time-forum/internal/models"

type CategoryUsecase interface {
	GetAllCategories() (categories []models.Category, err error)
}
