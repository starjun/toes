// Package middleware 提供 HTTP 中间件。
//
// 该包包含各种 HTTP 中间件，用于日志记录、认证、
// 限流、CORS 等横切关注点。
//
// 主要中间件:
//   - Logger: 请求日志
//   - CORS: 跨域支持
//   - RequestID: 请求追踪
//   - RateLimit: 限流
//
// 使用示例:
//
//	r.Use(middleware.Logger())
//	r.Use(middleware.Cors)
package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ErrLimitExceeded defines Limit exceeded error.
var ErrLimitExceeded = errors.New("Limit exceeded")

// Limit drops (HTTP status 429) the request if the limit is reached.
func Limit(maxEventsPerSec float64, maxBurstSize int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(maxEventsPerSec), maxBurstSize)

	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
			return
		}

		// Limit reached
		_ = c.Error(ErrLimitExceeded)
		c.AbortWithStatus(429)
	}
}
