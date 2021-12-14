package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/frank-yf/go-web-example/utils"
	"github.com/gin-gonic/gin"
)

// Ping 健康监测接口
func Ping(c *gin.Context) {
	renderOK(c)
}

// RedisPoolStats redis连接池统计数据
func RedisPoolStats(c *gin.Context) {
	ctx := context.TODO()
	if ping, ok := utils.PingRedis(ctx); !ok {
		renderError(c, fmt.Sprintf("ping redis failed: %s", ping))
		return
	}
	stats := utils.GetRedisCli().PoolStats()
	renderData(c, stats)
}

// LocalCacheStats cache统计数据
func LocalCacheStats(c *gin.Context) {
	stats := utils.GetCacheCli().Stats()
	renderData(c, stats)
}

// RedisSubscribes redis订阅列表
func RedisSubscribes(c *gin.Context) {
	channels := utils.GetRedisSubPool().Lookup()
	renderData(c, channels)
}

// CancelRedisSubscribe 取消指定通道的redis订阅
func CancelRedisSubscribe(c *gin.Context) {
	channel := c.Query("channel")
	loaded, err := utils.GetRedisSubPool().Unsubscribe(c, channel)
	if err != nil {
		renderError(c, fmt.Sprintf("cancel subscribe '%s' error : %s", channel, err.Error()))
		return
	}
	if !loaded {
		renderError(c, fmt.Sprint("channel exist : ", channel))
		return
	}
	renderOK(c)
}

// Recovery 统一处理接口调用过程中的panic，避免影响web服务
func Recovery(c *gin.Context, recovered interface{}) {
	if err, ok := recovered.(string); ok {
		renderServerError(c, err)
	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
