package concurrentratelimiter

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"strings"
	"time"

	"github.com/xiaojiaoyu100/lizard/redispattern"
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

var (
	enterScriptDigest string
	leaveScriptDigest string
)

func init() {
	e := sha1.New()
	io.WriteString(e, enterScript)
	enterScriptDigest = hex.EncodeToString(e.Sum(nil))

	l := sha1.New()
	io.WriteString(l, leaveScript)
	leaveScriptDigest = hex.EncodeToString(l.Sum(nil))
}

type Setting func(o *Option) error

type Option struct {
	ttl   int64 // time to live in millisecond
	limit int64 // maximum running limit
}

func WithTTL(ttl time.Duration) Setting {
	return func(o *Option) error {
		o.ttl = timekit.DurationToMillis(ttl)
		return nil
	}
}

func WithLimit(limit int64) Setting {
	return func(o *Option) error {
		o.limit = limit
		return nil
	}
}

type ConcurrentRateLimiter struct {
	runner redispattern.Runner
	key    string
	option Option
}

func New(runner redispattern.Runner, key string, settings ...Setting) (*ConcurrentRateLimiter, error) {
	c := &ConcurrentRateLimiter{
		runner: runner,
		key:    key,
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

func (c *ConcurrentRateLimiter) Enter(random string) (bool, error) {
	ok, err := c.runner.EvaSha1(enterScriptDigest,
		c.key,
		c.option.limit,
		timekit.NowInMillis(),
		random,
		c.option.ttl,
	)
	if err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT") {
		ok, err := c.runner.Eva(enterScript,
			c.key,
			c.option.limit,
			timekit.NowInMillis(),
			random,
			c.option.ttl,
		)
		if err != nil {
			return false, err
		}
		return ok == 1, nil
	}
	if err != nil {
		return false, err
	}
	return ok == 1, nil
}

func (c *ConcurrentRateLimiter) Leave(random string) error {
	_, err := c.runner.EvaSha1(leaveScriptDigest, c.key, random)
	if err != nil && strings.HasPrefix(err.Error(), "NOSCRIPT") {
		_, err := c.runner.Eva(leaveScript, c.key, random)
		if err != nil {
			return err
		}
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
