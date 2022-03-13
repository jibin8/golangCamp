package middleware

import (
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	config2 "self/internal/config"
	ginplus2 "self/internal/ginplus"
	"self/pkg/errors"
)

// CasbinMiddleware casbin中间件
func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) gin.HandlerFunc {
	cfg := config2.C.Casbin
	if !cfg.Enable {
		return EmptyMiddleware()
	}

	return func(c *gin.Context) {
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		p := c.Request.URL.Path
		m := c.Request.Method
		b := enforcer.Enforce(ginplus2.GetAccountKey(c), "Auth", "api", p, m)
		if !b {
			ginplus2.ResError(c, errors.ErrNoPerm)
			return
		}
		c.Next()
	}
}
