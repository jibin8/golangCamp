package middleware

import (
	config2 "self/internal/config"
	ginplus2 "self/internal/ginplus"
	icontext2 "self/internal/icontext"
	"self/pkg/auth"
	"self/pkg/errors"
	"self/pkg/logger"
	"github.com/gin-gonic/gin"
)

func wrapAccountAuthContext(c *gin.Context, accountKey string) {
	ginplus2.SetAccountKey(c, accountKey)
	ctx := icontext2.NewAccountKey(c.Request.Context(), accountKey)
	ctx = logger.NewAccountKeyContext(ctx, accountKey)
	c.Request = c.Request.WithContext(ctx)
}

// UserAuthMiddleware 用户授权中间件
func UserAuthMiddleware(a auth.Auther, skippers ...SkipperFunc) gin.HandlerFunc {
	if !config2.C.JWTAuth.Enable {
		return func(c *gin.Context) {
			wrapAccountAuthContext(c, config2.C.Root.Username)
			c.Next()
		}
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		userID, info, err := a.ParseUserID(c.Request.Context(), ginplus2.GetToken(c))
		if err != nil {
			if err == auth.ErrInvalidToken {
				if config2.C.IsDebugMode() {
					wrapAccountAuthContext(c, config2.C.Root.Username)
					c.Next()
					return
				}
				ginplus2.ResError(c, errors.ErrInvalidToken)
				return
			}

			//ginplus.ResError(c, errors.WithStack(err))
			ginplus2.ResError(c, errors.ErrInvalidToken)
			return
		}

		//或者可利用json序列化再反序列化取出数据
		_, ok := info["AppKeys"]
		if !ok {
			ginplus2.ResError(c, errors.ErrInvalidToken)
			return
		}

		appKeys, ok := info["AppKeys"].([]interface{})
		if !ok {
			ginplus2.ResError(c, errors.ErrInvalidToken)
			return
		}

		appKey := c.GetHeader("RJ-AppKey")
		support := false
		for _, v := range appKeys {
			if rv, ok := v.(string); ok && rv == appKey {
				support = true
				break
			}
		}
		if !support {
			ginplus2.ResError(c, errors.ErrInvalidToken)
			return
		}

		wrapAccountAuthContext(c, userID)
		c.Next()
	}
}
