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
	UpdateStatus(userID int64, status string) error
}
