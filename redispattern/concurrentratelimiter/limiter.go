package concurrentratelimiter

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"time"

	"github.com/xiaojiaoyu100/lizard/timekit"
)

const (
	enterScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local now = tonumber(ARGV[2])
local random = ARGV[3]
local ttl = tonumber(ARGV[4])

redis.call('zremrangebyscore', key, '-inf', now - ttl)

local count = redis.call("zcard", key)

if count < limit then
	redis.call("zadd", key, now, random)
	return 1
end

return 0
`
	leaveScript = `
local key = KEYS[1]
local random = ARGV[1]
local ret = redis.call("zrem", key, random)
return ret
`
)

func enterScriptDigest() (string, error) {
	hash := sha1.New()
	_, err := io.WriteString(hash, enterScript)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func leaveScriptDigest() (string, error) {
	hash := sha1.New()
	_, err := io.WriteString(hash, leaveScript)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Setting configures option.
type Setting func(o *Option) error

// Option 配置
type Option struct {
	ttl   int64 // time to live in millisecond
	limit int64 // maximum running limit
}

// WithTTL 存活期
func WithTTL(ttl time.Duration) Setting {
	return func(o *Option) error {
		o.ttl = timekit.DurationToMillis(ttl)
		return nil
	}
}

// WithLimit 上限
func WithLimit(limit int64) Setting {
	return func(o *Option) error {
		o.limit = limit
		return nil
	}
}

// ConcurrentRateLimiter 并发限流器
type ConcurrentRateLimiter struct {
	redis  rediser
	key    string
	option Option
}

// New 生成并发限流器
func New(redis rediser, key string, settings ...Setting) (*ConcurrentRateLimiter, error) {
	c := &ConcurrentRateLimiter{
		redis: redis,
		key:   key,
	}
	o := Option{
		ttl:   timekit.DurationToMillis(3 * time.Second),
		limit: 10,
	}
	for _, setting := range settings {
		if err := setting(&o); err != nil {
			return nil, err
		}
	}
	c.option = o
	return c, nil
}

// Enter 消耗
func (c *ConcurrentRateLimiter) Enter(random string) (bool, error) {
	d, err := enterScriptDigest()
	if err != nil {
		return false, err
	}
	exist, err := c.redis.ScriptExists(d).Result()
	if err != nil {
		return false, err
	}
	if !exist[0] {
		_, err := c.redis.ScriptLoad(enterScript).Result()
		if err != nil {
			return false, err
		}
	}
	ret, err := c.redis.EvalSha(d,
		[]string{c.key},
		c.option.limit,
		timekit.NowInMillis(),
		random,
		c.option.ttl,
	).Result()
	if err != nil {
		return false, err
	}

	r, ok := ret.(int64)
	if !ok {
		return false, errors.New("unexpected")
	}

	return r == 1, nil
}

// Leave 恢复
func (c *ConcurrentRateLimiter) Leave(random string) error {
	d, err := leaveScriptDigest()
	if err != nil {
		return err
	}
	exist, err := c.redis.ScriptExists(d).Result()
	if err != nil {
		return err
	}
	if !exist[0] {
		_, err := c.redis.ScriptLoad(leaveScript).Result()
		if err != nil {
			return err
		}
	}
	_, err = c.redis.EvalSha(d, []string{c.key}, random).Result()
	if err != nil {
		return err
	}
	return nil
}
