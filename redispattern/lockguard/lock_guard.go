package lockguard

import (
	"crypto/rand"
	"crypto/rc4"

	"time"

	"github.com/go-redis/redis"
	lock "github.com/xiaojiaoyu100/lizard/redispattern/lockguard/internal"
)

const (
	redisLockKey = "HHsYC5oVzLjFuWE4KMz923QT"

	delLuaScript = `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`
)

type LockGuard struct {
	lock lock.Lock
}

type Setter func(l *lock.Lock)

func WithExpiration(expiration time.Duration) Setter {
	return func(l *lock.Lock) {
		l.Expiration = expiration
	}
}

func New(client *redis.Client, key string, setters ...Setter) *LockGuard {
	guard := new(LockGuard)
	l := lock.Lock{
		Client:     client,
		Key:        key,
		Value:      "",
		Expiration: 30 * time.Second,
	}
	for _, setter := range setters {
		setter(&l)
	}
	guard.lock = l
	return guard
}

func (guard *LockGuard) Lock() bool {
	src := make([]byte, len(redisLockKey))
	_, err := rand.Read(src)
	if err != nil {
		return false
	}
	redisLockKeyByte := make([]byte, len(redisLockKey))
	copy(redisLockKeyByte[:], redisLockKey)
	cipher, err := rc4.NewCipher(redisLockKeyByte)
	if err != nil {
		return false
	}
	cipher.XORKeyStream(src, src)
	guard.lock.Value = string(src)
	cmd := guard.lock.Client.SetNX(guard.lock.Key, guard.lock.Value, guard.lock.Expiration)
	flag, err := cmd.Result()
	if err != nil {
		return false
	}
	return flag
}

func (guard *LockGuard) UnLock() {
	if len(guard.lock.Key) > 0 && len(guard.lock.Value) > 0 {
		keys := []string{guard.lock.Key}
		guard.lock.Client.Eval(delLuaScript, keys, guard.lock.Value)
	}
}
