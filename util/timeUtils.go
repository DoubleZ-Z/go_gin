package util

import (
	"fmt"
	"time"
)

func UnixToTimeString(timestamp int64) string {
	fmt.Println(timestamp)
	unix := time.Unix(int64(timestamp), 0)
	return unix.Format("2006-01-02 15:04:05")
}

func UnixToTimestampString(timestamp int64) string {
	fmt.Println(timestamp)
	unix := time.Unix(int64(timestamp), 0)
	return unix.Format("2006-01-02 15:04:05.111")
}
