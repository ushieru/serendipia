package utils

import "time"

func GetTimeStamp() int64 {
	return time.Now().UTC().UnixMilli() / 1e3
}
