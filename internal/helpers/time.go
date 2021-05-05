package helpers

import "time"

func GetTime() int64 {
	return time.Now().Unix()
}
