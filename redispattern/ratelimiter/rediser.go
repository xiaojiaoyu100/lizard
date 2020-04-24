package ratelimiter

import "github.com/go-redis/redis"

var (
	_ rediser = (*redis.Client)(nil)
	_ rediser = (*redis.Ring)(nil)
	_ rediser = (*redis.ClusterClient)(nil)
)

type rediser interface {
	LLen(key string) *redis.IntCmd
	Exists(keys ...string) *redis.IntCmd
	TxPipeline() redis.Pipeliner
	RPushX(key string, values interface{}) *redis.IntCmd
}
