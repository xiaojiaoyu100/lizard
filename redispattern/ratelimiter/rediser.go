package ratelimiter

import "github.com/go-redis/redis"

type rediser interface {
	LLen(key string) *redis.IntCmd
	Exists(keys ...string) *redis.IntCmd
	TxPipeline() redis.Pipeliner
	RPushX(key string, values ...interface{}) *redis.IntCmd
}
