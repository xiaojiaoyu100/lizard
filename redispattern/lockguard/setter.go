package lockguard

// Setter 配置lock.
type Setter func(l *Lock) error

// WithRetryTimes configures lock retry times.
func WithRetryTimes(t int) Setter {
	return func(l *Lock) error {
		l.retryTimes = t
		return nil
	}
}
