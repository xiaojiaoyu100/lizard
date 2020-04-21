package ratelimiter

import (
	"time"
)

// RateLimiter 限流器
type RateLimiter interface {
	Limit() bool
}

type rateLimiter struct {
	redis    rediser
	key      string
	duration time.Duration
	limit    int64
}

type option struct {
	duration time.Duration
	limit    int64
}

// Setter 配置
type Setter func(*option)

// WithDuration 设置存活时间
func WithDuration(duration time.Duration) Setter {
	return func(o *option) {
		o.duration = duration
	}
}

// WithLimit 设置上限
func WithLimit(limit int64) Setter {
	return func(o *option) {
		o.limit = limit
	}
}

var defaultRateLimiterOption = option{
	duration: 1 * time.Second,
	limit:    1,
}

// New 限流器
func New(redis rediser, key string, setters ...Setter) RateLimiter {
	option := defaultRateLimiterOption
	for _, setter := range setters {
		setter(&option)
	}
	return &rateLimiter{
		redis:    redis,
		key:      key,
		duration: option.duration,
		limit:    option.limit,
	}
}

// Limit 触发限流
func (rl *rateLimiter) Limit() bool {
	current, _ := rl.redis.LLen(rl.key).Result()
	value := "0"
	if current >= rl.limit {
		return true
	}

	exist, err := rl.redis.Exists(rl.key).Result()
	if err != nil {
		return true
	}
	if exist == 0 {
		pipe := rl.redis.TxPipeline()
		pipe.RPush(rl.key, value)
		pipe.Expire(rl.key, rl.duration)
		_, err := pipe.Exec()
		if err != nil {
			return true
		}
	} else {
		rl.redis.RPushX(rl.key, value)
	}

	return false
}
