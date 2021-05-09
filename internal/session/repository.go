package session

import (
	"github.com/innovember/real-time-forum/internal/models"
)

type SessionRepository interface {
	// Sessions
	SelectByToken(token string) (*models.Session, error)
	Insert(*models.Session) error
	Delete(token string) error
	DeleteTokens() error

	// Online Users
	InsertOnlineUser(userID int64) error
	DeleteOnlineUser(userID int64) error
	DeleteOnlineUsers() error
}
