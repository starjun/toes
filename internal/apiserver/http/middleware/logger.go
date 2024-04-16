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
