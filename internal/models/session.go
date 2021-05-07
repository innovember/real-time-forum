package models

import (
	"time"

	"github.com/innovember/real-time-forum/internal/consts"
	uuid "github.com/satori/go.uuid"
)

type Session struct {
	Name      string    `json:"name"`
	UserID    int64     `json:"userID"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

func NewSession(userID int64) *Session {
	return &Session{
		Token:     uuid.NewV4().String(),
		Name:      consts.SessionName,
		UserID:    userID,
		ExpiresAt: time.Now().Add(consts.SessionExpireDuration),
	}
}
