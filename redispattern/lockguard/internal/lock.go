package lock

import (
	"time"

	"github.com/go-redis/redis"
)

type Lock struct {
	Client     *redis.Client
	Key        string
	Value      string
	Expiration time.Duration
}
