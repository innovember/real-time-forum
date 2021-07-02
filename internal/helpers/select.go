package helpers

import "github.com/innovember/real-time-forum/internal/models"

func SelectSecondUser(users []models.User, userID int64) int64 {
	for _, u := range users {
		if u.ID != userID {
			return u.ID
		}
	}
	return -1
}
