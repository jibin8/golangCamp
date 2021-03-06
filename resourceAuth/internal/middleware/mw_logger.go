package middleware

import (
	ginplus2 "self/internal/ginplus"
	"mime"
	"net/http"
	"time"

	"self/pkg/logger"
	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware(skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		method := c.Request.Method
		span := logger.StartSpan(c.Request.Context(),
			logger.SetSpanTitle("Request"),
			logger.SetSpanFuncName(JoinRouter(method, p)))

		start := time.Now()

		fields := make(map[string]interface{})
		fields["ip"] = c.ClientIP()
		fields["method"] = method
		fields["url"] = c.Request.URL.String()
		fields["proto"] = c.Request.Proto
		fields["header"] = c.Request.Header
		fields["user_agent"] = c.GetHeader("User-Agent")
		fields["content_length"] = c.Request.ContentLength

		if method == http.MethodPost || method == http.MethodPut {
			mediaType, _, _ := mime.ParseMediaType(c.GetHeader("Content-Type"))
			if mediaType != "multipart/form-data" {
				if v, ok := c.Get(ginplus2.ReqBodyKey); ok {
					if b, ok := v.([]byte); ok {
						fields["body"] = string(b)
					}
				}
			}
		}
		c.Next()

		timeConsuming := time.Since(start).Nanoseconds() / 1e6
		fields["res_status"] = c.Writer.Status()
		fields["res_length"] = c.Writer.Size()

		if v, ok := c.Get(ginplus2.LoggerReqBodyKey); ok {
			if b, ok := v.([]byte); ok {
				fields["body"] = string(b)
			}
		}

		if v, ok := c.Get(ginplus2.ResBodyKey); ok {
			if b, ok := v.([]byte); ok {
				fields["res_body"] = string(b)
			}
		}

		fields[logger.AccountKeyKey] = ginplus2.GetAccountKey(c)
		span.WithFields(fields).Infof("[http] %s-%s-%s-%d(%dms)",
			p, c.Request.Method, c.ClientIP(), c.Writer.Status(), timeConsuming)
	}
}
