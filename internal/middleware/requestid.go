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
			//rip = c.Request.RemoteAddr
			rip = c.ClientIP()
		}
		c.Set(global.Cfg.Header.Realip, rip)
		c.Next()
	}
}
