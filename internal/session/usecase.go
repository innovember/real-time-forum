package session

import "github.com/innovember/real-time-forum/internal/models"

type SessionUsecase interface {
	CreateSession(session *models.Session) error
	GetByToken(token string) (*models.Session, error)
	DeleteSession(token string) error
	DeleteExpiredSessions()
}
