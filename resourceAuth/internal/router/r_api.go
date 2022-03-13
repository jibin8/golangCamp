package router

import (
	"github.com/gin-gonic/gin"
	middleware2 "self/internal/middleware"
)

// RegisterAPI register api group router
func (a *Router) RegisterAPI(app *gin.Engine) {
	// 权限系统内部接口
	g := app.Group("/api")

	// 跳过token检查
	g.Use(middleware2.UserAuthMiddleware(a.Auth,
		middleware2.AllowPathPrefixSkipper("/api/pub/login"),
	))

	// 跳过鉴权
	g.Use(middleware2.CasbinMiddleware(a.CasbinEnforcer,
		middleware2.AllowPathPrefixSkipper("/api/pub"),
		middleware2.AllowPathPrefixSkipper("/api/external"),
	))

	pub := g.Group("/pub")
	{
		gLogin := pub.Group("login")
		{
			gLogin.POST("", a.LoginAPI.Login)
		}


		gUser := pub.Group("users")
		{
			gUser.POST("", a.AccountAPI.Create)
		}
	}

	external := g.Group("/external")
	{
		gLogin := external.Group("login")
		{
			gLogin.GET("authenticate", a.LoginAPI.Authenticate)
		}
	}
}
