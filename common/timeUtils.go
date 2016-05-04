package common

import (
	"syscall"
	"time"
)

func GetTimeNowMillisecond() int64 {

	var tv syscall.Timeval
	syscall.Gettimeofday(&tv)

	return (int64(tv.Sec) * 1e3 + int64(tv.Usec) / 1e3)
}

func MillisecondToTime(ms int64) time.Time {

	return time.Unix(0, ms * int64(time.Millisecond))
}