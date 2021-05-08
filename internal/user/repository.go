package user

import "github.com/innovember/real-time-forum/internal/models"

type UserRepository interface {
	Insert(user *models.User) error
	SelectByEmail(email string) (*models.User, error)
	SelectByNickname(username string) (*models.User, error)
	SelectByID(userID int64) (*models.User, error)
	UpdateActivity(userID int64) (err error)
}
