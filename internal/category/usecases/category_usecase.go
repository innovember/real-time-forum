package usecases

import (
	"github.com/innovember/real-time-forum/internal/category"
)

type CategoryUsecase struct {
	categoryRepo category.CategoryRepository
}

func NewCategoryUsecase(categoryRepo category.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: categoryRepo,
	}
}
