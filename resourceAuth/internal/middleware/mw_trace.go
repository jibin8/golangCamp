package middleware

import (
	"github.com/gin-gonic/gin"
	icontext2 "self/internal/icontext"
	"self/pkg/logger"
	"self/pkg/trace"
)

// TraceMiddleware 跟踪ID中间件
func TraceMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		// 优先从请求头中获取请求ID
		traceID := c.GetHeader("X-Request-Id")
		if traceID == "" {
			traceID = trace.NewID()
		}

		ctx := icontext2.NewTraceID(c.Request.Context(), traceID)
		ctx = logger.CreateTraceIDContext(ctx, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-Trace-Id", traceID)

		c.Next()
	}
}
