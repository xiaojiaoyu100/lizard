package ratecounter

import (
	"sync/atomic"
	"time"
)

// Counter 计数
type Counter int64

// Reset 重置
func (c *Counter) Reset() {
	atomic.StoreInt64((*int64)(c), 0)
}

// Incr 增加
func (c *Counter) Incr(value int64) {
	atomic.AddInt64((*int64)(c), value)
}

// Value 当前值
func (c *Counter) Value() int64 {
	return atomic.LoadInt64((*int64)(c))
}

// RateCounter 速率计数器
type RateCounter struct {
	counter    Counter
	interval   time.Duration
	resolution int
	partials   []Counter
	current    int32
	running    int32
}

type option struct {
	Resolution int
	Interval   time.Duration
}

// Setter configures option.
type Setter func(option *option)

// WithInterval configures interval.
func WithInterval(interval time.Duration) Setter {
	return func(option *option) {
		option.Interval = interval
	}
}

// WithResolution configures resolution.
func WithResolution(resolution int) Setter {
	if resolution < 1 {
		resolution = 1
	}
	return func(option *option) {
		option.Resolution = resolution
	}
}

// New generates a RateCounter.
func New(setters ...Setter) *RateCounter {
	rateCounter := &RateCounter{}
	option := option{
		Interval:   time.Second,
		Resolution: 20,
	}
	for _, setter := range setters {
		setter(&option)
	}
	rateCounter.interval = option.Interval
	rateCounter.resolution = option.Resolution
	rateCounter.partials = make([]Counter, option.Resolution)
	rateCounter.current = 0
	rateCounter.running = 0
	return rateCounter
}

func (r *RateCounter) run() {
	if ok := atomic.CompareAndSwapInt32(&r.running, 0, 1); !ok {
		return
	}

	go func() {
		ticker := time.NewTicker(time.Duration(float64(r.interval) / float64(r.resolution)))

		for range ticker.C {
			current := atomic.LoadInt32(&r.current)
			next := (int(current) + 1) % r.resolution
			r.counter.Incr(-1 * r.partials[next].Value())
			r.partials[next].Reset()
			atomic.CompareAndSwapInt32(&r.current, current, int32(next))
			if r.counter.Value() == 0 {
				atomic.StoreInt32(&r.running, 0)
				ticker.Stop()

				return
			}
		}
	}()
}

// Incr 增加计数
func (r *RateCounter) Incr(val int64) {
	r.counter.Incr(val)
	r.partials[atomic.LoadInt32(&r.current)].Incr(val)
	r.run()
}

// Rate 当前速率
func (r *RateCounter) Rate() int64 {
	return r.counter.Value()
}
