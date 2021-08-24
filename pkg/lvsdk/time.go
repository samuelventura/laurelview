package lvsdk

import "time"

func Millis(ms int) time.Duration {
	return time.Duration(ms) * time.Millisecond
}

func Future(ms int) time.Time {
	d := Millis(ms)
	return time.Now().Add(d)
}
