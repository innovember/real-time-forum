package helpers

import "time"

func GetCurrentUnixTime() int64 {
	return time.Now().Unix()
}
