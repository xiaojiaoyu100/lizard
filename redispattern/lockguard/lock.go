package lockguard

import (
	"time"
)

type Lock struct {
	redis      rediser
	Key        string
	Value      string
	Expiration time.Duration
}
