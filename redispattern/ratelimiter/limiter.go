package ratelimiter

import (
	"time"
)

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

type Setter func(*option)

func WithDuration(duration time.Duration) Setter {
	return func(o *option) {
		o.duration = duration
	}
}

func WithLimit(limit int64) Setter {
	return func(o *option) {
		o.limit = limit
	}
}

var defaultRateLimiterOption = option{
	duration: 1 * time.Second,
	limit:    1,
}

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

func (rl *rateLimiter) Limit() bool {
	current, _ := rl.redis.LLen(rl.key).Result()
	value := "0"
	if current >= rl.limit {
		return true
	} else {
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
	}
	return false
}
