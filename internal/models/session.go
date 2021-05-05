package models

type Session struct {
	UserID    int    `json:"userID"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
}
