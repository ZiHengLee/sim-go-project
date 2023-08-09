package time

import (
	"time"
)

func getTime(t ...time.Time) *time.Time {
	if len(t) == 0 {
		n := time.Now()
		return &n
	}
	return &t[0]
}

func Now() time.Time {
	return time.Now()
}

func Unix(t ...time.Time) int64 {
	return getTime(t...).Unix()
}

func UnixMilli(t ...time.Time) int64 {
	return getTime(t...).UnixNano() / 1000000
}

func UnixMicro(t ...time.Time) int64 {
	return getTime(t...).UnixNano() / 1000
}

func UnixNano(t ...time.Time) int64 {
	return getTime(t...).UnixNano()
}
