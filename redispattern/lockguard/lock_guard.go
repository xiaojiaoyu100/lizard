package lockguard

import (
	"crypto/rand"
	"crypto/rc4"
	"time"
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

// LockGuard provides distributed lock.
type LockGuard struct {
	lock Lock
}

// Setter 配置lock.
type Setter func(l *Lock)

// WithExpiration 设置过期时间
func WithExpiration(expiration time.Duration) Setter {
	return func(l *Lock) {
		l.Expiration = expiration
	}
}

// New 生成一个锁
func New(redis rediser, key string, setters ...Setter) *LockGuard {
	guard := new(LockGuard)
	l := Lock{
		redis:      redis,
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

// Lock 锁住
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
	cmd := guard.lock.redis.SetNX(guard.lock.Key, guard.lock.Value, guard.lock.Expiration)
	flag, err := cmd.Result()
	guard.lock.locked = flag
	if err != nil {
		return false
	}
	return flag
}

// UnLock 解锁
func (guard *LockGuard) UnLock() {
	if !guard.lock.locked {
		return
	}
	keys := []string{guard.lock.Key}
	guard.lock.redis.Eval(delLuaScript, keys, guard.lock.Value)
}
