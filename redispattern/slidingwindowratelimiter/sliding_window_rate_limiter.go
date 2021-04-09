package slidingwindowratelimiter

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	script = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2]) 
local limit = tonumber(ARGV[3])
local pivot = now - window

redis.call('ZREMRANGEBYSCORE', key, 0, pivot)

local count = redis.call('ZCARD', key)
if count < limit then 
	redis.call('ZADD', key, now, now)
end

redis.call('EXPIRE', key, window / 1000000000)

return limit - count
`
)

func scriptDigest() (string, error) {
	s := sha1.New()
	_, err := io.WriteString(s, script)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(s.Sum(nil)), nil
}

// SlidingWindowRateLimiter represents a sliding window rate limiter.
type SlidingWindowRateLimiter struct {
	redis  rediser
	key    string
	window time.Duration
	limit  int64
}

// Option is necessary for creating a limiter.
type Option struct {
	Redis  rediser
	Key    string
	Window time.Duration
	Limit  int64
}

// New generates a rate limiter.
func New(o *Option) (*SlidingWindowRateLimiter, error) {
	if o.Window.Seconds() < 1 {
		return nil, errors.New("below one second is not supported")
	}
	sl := &SlidingWindowRateLimiter{
		redis:  o.Redis,
		key:    o.Key,
		window: o.Window,
		limit:  o.Limit,
	}
	return sl, nil
}

// Allow returns the sliding window rate limiter status.
func (sl *SlidingWindowRateLimiter) Allow() (bool, error) {
	digest, err := scriptDigest()
	if err != nil {
		return false, err
	}
	exist, err := sl.redis.ScriptExists(digest).Result()
	if err != nil {
		return false, err
	}
	if !exist[0] {
		_, err := sl.redis.ScriptLoad(script).Result()
		if err != nil {
			return false, err
		}
	}
	ret, err := sl.redis.EvalSha(digest, []string{sl.key}, time.Now().UnixNano(), sl.window.Nanoseconds(), sl.limit).Result()
	if err != nil {
		return false, err
	}
	switch v := ret.(type) {
	case int64:
		if v <= 0 {
			return false, errors.New("limit reached")
		}
		return true, nil
	default:
		return false, fmt.Errorf("sliding window rate limiter value: %#v, key = %s, window = %s, limit = %d", ret, sl.key, sl.window, sl.limit)
	}
}
