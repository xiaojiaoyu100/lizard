package lockguard

import (
	"context"
	"crypto/rand"
	"crypto/rc4"
	"errors"
	"fmt"
	"time"

	"github.com/xiaojiaoyu100/lizard/backoff"
	"github.com/xiaojiaoyu100/lizard/errorkit"
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

// New 生成一个锁
func New(redis rediser, key string, setters ...Setter) (*LockGuard, error) {
	if key == "" {
		return nil, errors.New("key length is zero")
	}

	guard := new(LockGuard)
	l := Lock{
		redis:      redis,
		Key:        key,
		Value:      "",
		retryTimes: 3,
		expiration: 1 * time.Minute,
	}
	for _, setter := range setters {
		if err := setter(&l); err != nil {
			return nil, err
		}
	}
	guard.lock = l
	return guard, nil
}

// Run 锁住
func (guard *LockGuard) Run(ctx context.Context, handler Handler) error {
	for i := 0; i < guard.lock.retryTimes; i++ {
		guard.obtain()
		if !guard.lock.locked {
			if guard.lock.retryTimes > 1 {
				time.Sleep(backoff.ExponentialBackoffFullJitterStrategy{
					ExponentialBackoff: backoff.ExponentialBackoff{
						Base: 20 * time.Millisecond,
						Cap:  100 * time.Millisecond,
					}}.Backoff(i))
			}
			continue
		}

		t := time.NewTicker(guard.tickInterval())
		stop := make(chan struct{})
		errChan := make(chan error)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					if _, ok := r.(error); ok {
						errChan <- errorkit.WithStack(r.(error))
					} else {
						errChan <- errorkit.WithStack(fmt.Errorf("%+v", r))
					}
				}
				close(stop)
				errChan <- nil
			}()
			handler(ctx)
		}()

		go func() {
			for {
				select {
				case <-t.C:
					guard.renewTTL()
				case <-stop:
					return
				}
			}
		}()

		var err error
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case err = <-errChan:
		}
		guard.UnLock()
		t.Stop()
		return err
	}
	return fmt.Errorf("key: %s, err: %w", guard.lock.Key, errLockNotObtained)
}

func (guard *LockGuard) renewTTL() {
	guard.lock.redis.Expire(guard.lock.Key, guard.lock.expiration)
}

func (guard *LockGuard) tickInterval() time.Duration {
	return time.Second * 10
}

func (guard *LockGuard) obtain() {
	src := make([]byte, len(redisLockKey))
	_, err := rand.Read(src)
	if err != nil {
		return
	}
	redisLockKeyByte := make([]byte, len(redisLockKey))
	copy(redisLockKeyByte[:], redisLockKey)
	cipher, err := rc4.NewCipher(redisLockKeyByte)
	if err != nil {
		return
	}
	cipher.XORKeyStream(src, src)
	guard.lock.Value = string(src)
	cmd := guard.lock.redis.SetNX(guard.lock.Key, guard.lock.Value, guard.lock.expiration)
	flag, err := cmd.Result()
	if err != nil {
		return
	}
	guard.lock.locked = flag
}

// UnLock 解锁
func (guard *LockGuard) UnLock() {
	if !guard.lock.locked {
		return
	}
	keys := []string{guard.lock.Key}
	guard.lock.redis.Eval(delLuaScript, keys, guard.lock.Value)
}
