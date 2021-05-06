package user

import "github.com/innovember/real-time-forum/internal/models"

type UserUsecase interface {
	Create(user *models.User) error
	// GetByEmail(email string) (*models.User, error)
	// GetByUsername(email string) (*models.User, error)
	// GetByID(userID int64) (*models.User, error)
	// CheckPassword(user *models.User, password string) error
}
