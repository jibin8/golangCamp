package injector

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	config2 "self/internal/config"
	middleware2 "self/internal/middleware"
	router2 "self/internal/router"
)

// InitGinEngine 初始化gin引擎
func InitGinEngine(r router2.IRouter) *gin.Engine {
	gin.SetMode(config2.C.RunMode)

	app := gin.New()
	app.NoMethod(middleware2.NoMethodHandler())
	app.NoRoute(middleware2.NoRouteHandler())

	prefixes := r.Prefixes()

	// Trace ID
	app.Use(middleware2.TraceMiddleware(middleware2.AllowPathPrefixNoSkipper(prefixes...)))

	// Copy body
	app.Use(middleware2.CopyBodyMiddleware(middleware2.AllowPathPrefixNoSkipper(prefixes...)))

	// Access logger
	app.Use(middleware2.LoggerMiddleware(middleware2.AllowPathPrefixNoSkipper(prefixes...)))

	// Router register
	r.Register(app)

	// Swagger
	if config2.C.Swagger {
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return app
}
