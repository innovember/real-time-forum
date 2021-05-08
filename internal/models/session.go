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
	Expires   time.Time `json:"expires"`
	ExpiresAt int64     `json:"expiresAt"`
}

func NewSession(userID int64) *Session {
	return &Session{
		Token:     uuid.NewV4().String(),
		Name:      consts.SessionName,
		UserID:    userID,
		Expires:   time.Now().Add(consts.SessionExpireDuration),
		ExpiresAt: time.Now().Add(consts.SessionExpireDuration).Unix(),
	}
}
