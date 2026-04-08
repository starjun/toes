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
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"toes/global"
)

// RequestID is a middleware that injects a 'X-Request-ID' into the context and request/response header of each request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		rid := c.GetHeader(global.Cfg.Header.Requestid)

		if rid == "" {
			rid = uuid.New().String()
			c.Set(global.Cfg.Header.Requestid, rid)
		}

		// Set known.XRequestIDKey header
		c.Writer.Header().Set(global.Cfg.Header.Requestid, rid)
		c.Next()
	}
}

func RealIp() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		rip := c.GetHeader(global.Cfg.Header.Realip)
		if rip == "" {
			// rip = c.Request.RemoteAddr
			rip = c.ClientIP()
		}
		c.Set(global.Cfg.Header.Realip, rip)
		c.Next()
	}
}
