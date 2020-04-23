package lockguard

import (
	"time"
)

// Setter 配置lock.
type Setter func(l *Lock) error

// WithExpiration 设置过期时间
func WithExpiration(expiration time.Duration) Setter {
	return func(l *Lock) error {
		l.expiration = expiration
		return nil
	}
}

// WithRetryTimes configures lock retry times.
func WithRetryTimes(t int) Setter {
	return func(l *Lock) error {
		l.retryTimes = t
		return nil
	}
}
