package ratelimiter

// Limiter limit interface.
type Limiter interface {
	Allow() (bool, func())
}
