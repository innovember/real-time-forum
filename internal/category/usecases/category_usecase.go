package usecases

import (
	"github.com/innovember/real-time-forum/internal/category"
	"github.com/innovember/real-time-forum/internal/models"
)

type CategoryUsecase struct {
	categoryRepo category.CategoryRepository
}

func NewCategoryUsecase(categoryRepo category.CategoryRepository) *CategoryUsecase {
	return &CategoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (cu *CategoryUsecase) GetAllCategories() (categories []models.Category, err error) {
	if categories, err = cu.categoryRepo.SelectAllCategories(); err != nil {
		return nil, err
	}
	return categories, nil
}
