package pprof

import (
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

func isInternalIP(ip string) bool {
	i := net.ParseIP(ip)
	if i.IsLoopback() {
		return true
	}
	if ip4 := i.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		default:
			return false
		}
	}
	return false
}

func auth(fn func(w http.ResponseWriter, r *http.Request)) func(c *gin.Context) {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !isInternalIP(ip) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		fn(c.Writer, c.Request)
	}
}

// InitRoutes inits pprof routes.
func InitRoutes(engine *gin.Engine) {
	engine.GET("/debug/pprof/", auth(pprof.Index))
	engine.GET("/debug/pprof/cmdline", auth(pprof.Cmdline))
	engine.GET("/debug/pprof/profile", auth(pprof.Profile))
	engine.GET("/debug/pprof/symbol", auth(pprof.Symbol))
	engine.GET("/debug/pprof/trace", auth(pprof.Trace))
	engine.GET("/debug/pprof/goroutine", auth(pprof.Index))
	engine.GET("/debug/pprof/allocs", auth(pprof.Index))
	engine.GET("/debug/pprof/block", auth(pprof.Index))
	engine.GET("/debug/pprof/heap", auth(pprof.Index))
	engine.GET("/debug/pprof/mutex", auth(pprof.Index))
	engine.GET("/debug/pprof/threadcreate", auth(pprof.Index))
}
