package lvsdk

import "time"

func Millis(ms int64) time.Duration {
	return time.Duration(ms) * time.Millisecond
}

func Future(ms int64) time.Time {
	d := Millis(ms)
	return time.Now().Add(d)
}
