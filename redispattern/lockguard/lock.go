package lockguard

import (
	"time"
)

// Lock redis lock
type Lock struct {
	redis      rediser
	Key        string
	Value      string
	locked     bool
	retryTimes int
	expiration time.Duration
}
