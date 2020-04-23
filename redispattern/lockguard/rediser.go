package lockguard

import (
	"time"

	"github.com/go-redis/redis"
)

type rediser interface {
	SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Eval(script string, keys []string, args ...interface{}) *redis.Cmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
}
