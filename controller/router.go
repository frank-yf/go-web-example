package controller

import (
	"github.com/frank-yf/go-web-example/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type routerRegister func(*gin.Engine)

// InitRouter 加载路由
func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Use(gin.CustomRecovery(Recovery)) // panic处理

	router.GET("/ping", Ping) // 心跳监测

	for _, register := range []routerRegister{registeLogic, registeHandle} {
		register(router)
	}

	utils.GetLogger().Debug("Initial router")
	return router
}

// registeLogic 业务逻辑相关接口
func registeLogic(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			renderData(c, "v1 response")
		})
	}
}

// registeHandle 服务管理相关接口
func registeHandle(r *gin.Engine) {
	handler := r.Group("/handler", Authorization)
	{
		handler.GET("/redis_stats", RedisPoolStats)
		handler.GET("/cache_stats", LocalCacheStats)

		redisSubRouter := handler.Group("/redis_sub")
		{
			redisSubRouter.GET("/", RedisSubscribes)
			redisSubRouter.GET("/cancel", CancelRedisSubscribe)
		}

		pprofRouter := handler.Group("/pprof")
		{
			pprofRouter.GET("/", PprofIndex)
			pprofRouter.GET("/cmdline", PprofCmdline)
			pprofRouter.GET("/profile", PprofProfile)
			pprofRouter.GET("/symbol", PprofSymbol)
			pprofRouter.POST("/symbol", PprofSymbol)
			pprofRouter.GET("/trace", PprofTrace)
			pprofRouter.GET("/:handler", PprofHandler)
		}
	}
	utils.GetLogger().Debug("Initial Handle router")

	fields := make([]zap.Field, 0, len(accounts))
	for k, v := range accounts {
		fields = append(fields, zap.String(k, v))
	}
	utils.GetLogger().Debug("Access handle interface with accounts", fields...)
}
