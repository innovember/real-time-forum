package helpers

import (
	"time"

	"github.com/innovember/real-time-forum/internal/consts"
)

func GetCurrentUnixTime() int64 {
	return time.Now().Unix()
}

func GetSessionExpireTime() int64 {
	return time.Now().Add(consts.SessionExpireDuration).Unix()
}
