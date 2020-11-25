package slidingwindowratelimiter

import (
	"crypto/sha1"
	"encoding/hex"
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

// New generates a rate limiter.
func New(redis rediser, key string, window time.Duration, limit int64) (*SlidingWindowRateLimiter, error) {
	sl := &SlidingWindowRateLimiter{
		redis:  redis,
		key:    key,
		window: window,
		limit:  limit,
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
		return v > 0, nil
	default:
		return false, fmt.Errorf("sliding window rate limiter err: %#v, key = %s, window = %s, limit = %d", ret, sl.key, sl.window, sl.limit)
	}
}
