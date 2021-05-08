package limitermiddleware

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xiaojiaoyu100/lizard/ratelimiter"
	"github.com/xiaojiaoyu100/lizard/ratelimiter/bbr"
)

type LimiterMiddleware struct {
	abortHttpErr  interface{}
	latestLogTime int64
	limiterGroup  *bbr.Group
	logger        ratelimiter.Logger
}

func NewLimiterMiddleware(conf *bbr.Config, logger ratelimiter.Logger, abortHttpErr interface{}) *LimiterMiddleware {
	return &LimiterMiddleware{
		limiterGroup: bbr.NewGroup(conf),
		abortHttpErr: abortHttpErr,
		logger:       logger,
	}
}

func (lm *LimiterMiddleware) RateLimit() gin.HandlerFunc {
	return func(context *gin.Context) {
		limiter := lm.limiterGroup.Get(fmt.Sprintf("%s_%s", context.Request.Method, context.Request.URL.Path))
		isAllow, done := limiter.Allow()
		if !isAllow {
			context.AbortWithStatusJSON(http.StatusTooManyRequests, lm.abortHttpErr)
			return
		}
		context.Next()
		if done != nil {
			done()
		}
		lm.print(context.Request.Method, context.Request.URL.Path, limiter)
	}
}

func (lm *LimiterMiddleware) print(method, path string, limiter ratelimiter.Limiter) {
	now := time.Now().UnixNano()
	if now-atomic.LoadInt64(&lm.latestLogTime) <= int64(time.Second*3) {
		return
	}
	atomic.StoreInt64(&lm.latestLogTime, now)
	ratelimiter.SafeLog(lm.logger, func(logger ratelimiter.Logger) {
		lm.logger.Printf("bbr method: %s path: %s stat:%+v", method, path, limiter.(*bbr.BBR).Stat())
	})
}
