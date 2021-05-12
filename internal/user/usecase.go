package user

import "github.com/innovember/real-time-forum/internal/models"

type UserUsecase interface {
	Create(user *models.User) error
	GetByNickname(nickname string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByID(userID int64) (*models.User, error)
	UpdateActivity(userID int64) error
	CheckPassword(*models.InputUserSignIn) error
	GetByEmailOrNickname(login string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
}
