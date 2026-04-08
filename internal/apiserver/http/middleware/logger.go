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
	"time"

	"toes/global"

	"github.com/gin-gonic/gin"
)

// Logger 接收gin框架默认的日志.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		_ip := c.GetString(global.Cfg.Header.Realip)
		global.LogGin(c).Sugar().Infow(
			"GIN",
			"status", c.Writer.Status(),
			"method", c.Request.Method,
			"path", path,
			"query", query,
			// "ip", c.ClientIP(),
			"ip", _ip,
			"cost", cost,
			"errors", c.Errors.ByType(gin.ErrorTypePrivate).String(),
		)
	}
}
