package common

import "syscall"

func GetTimeNowMillisecond() int64 {

	var tv syscall.Timeval
	syscall.Gettimeofday(&tv)

	return (int64(tv.Sec) * 1e3 + int64(tv.Usec) / 1e3)
}