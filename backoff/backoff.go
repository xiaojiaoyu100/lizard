package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Strategy defines an interface.
type Strategy interface {
	Backoff(retry int) time.Duration
}

// LinearBackoffStrategy 提供了线性时间重试
type LinearBackoffStrategy struct {
	slope time.Duration
}

// Backoff 线性时间重试函数
func (stg LinearBackoffStrategy) Backoff(retry int) time.Duration {
	return time.Duration(retry) * stg.slope
}

// ConstantBackOffStrategy 常量时间重试
type ConstantBackOffStrategy struct {
	interval time.Duration
}

// Backoff 重试
func (stg ConstantBackOffStrategy) Backoff(retry int) time.Duration {
	return stg.interval
}

// ExponentialBackoff 指数回退
type ExponentialBackoff struct {
	base time.Duration
	cap  time.Duration
}

func (backoff ExponentialBackoff) expo(retry int) float64 {
	c := float64(backoff.cap)
	b := float64(backoff.base)
	r := float64(retry)
	return math.Min(c, math.Exp2(r)*b)
}

// ExponentialBackoffStrategy 指数重试
type ExponentialBackoffStrategy struct {
	ExponentialBackoff
}

// Backoff 重试
func (stg ExponentialBackoffStrategy) Backoff(retry int) time.Duration {
	return time.Duration(stg.expo(retry))
}

// ExponentialBackoffEqualJitterStrategy 指数jitter重试
type ExponentialBackoffEqualJitterStrategy struct {
	ExponentialBackoff
}

// Backoff 重试
func (stg ExponentialBackoffEqualJitterStrategy) Backoff(retry int) time.Duration {
	v := stg.expo(retry)
	u := uniform(0, v/2.0)
	return time.Duration(v/2.0 + u)
}

// ExponentialBackoffFullJitterStrategy 指数full jitter重试
type ExponentialBackoffFullJitterStrategy struct {
	ExponentialBackoff
}

// Backoff 重试
func (stg ExponentialBackoffFullJitterStrategy) Backoff(retry int) time.Duration {
	v := stg.expo(retry)
	u := uniform(0, v)
	return time.Duration(u)
}

// ExponentialBackoffDecorrelatedJitterStrategy 指数decorrelated jitter重试
type ExponentialBackoffDecorrelatedJitterStrategy struct {
	ExponentialBackoff
	sleep time.Duration
}

// uniform returns a number in [min, max)
func uniform(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// Backoff 重试
func (stg ExponentialBackoffDecorrelatedJitterStrategy) Backoff(retry int) time.Duration {
	c := float64(stg.cap)
	b := float64(stg.base)
	s := float64(stg.sleep)
	u := uniform(b, 3*s)
	s = math.Min(c, u)
	return time.Duration(s)
}
