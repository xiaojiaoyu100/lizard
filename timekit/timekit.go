package timekit

import "time"

// DurationToMillis converts duration to milliseconds.
func DurationToMillis(d time.Duration) int64 {
	return int64(d / time.Millisecond)
}

// NowInMillis returns timestamp in milliseconds.
func NowInMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// NowInSecs returns timestamp in seconds.
func NowInSecs() int64 {
	return time.Now().Unix()
}

// UTCNowTime returns current time in utc.
func UTCNowTime() time.Time {
	return time.Now().UTC()
}
