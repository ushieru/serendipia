package utils

import "time"

func GetTimeStamp() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}
