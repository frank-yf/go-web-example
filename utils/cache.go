package utils

import (
	"sync"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	cacheClient *cache.Cache
	cacheOnce   sync.Once
)

func GetCacheCli() *cache.Cache {
	cacheOnce.Do(initCacheCli)
	return cacheClient
}

// RegisteDeleteCache 从redis中订阅清空内存缓存
// 假如订阅事件已存在，则什么都不做
// 异步注册，避免获取订阅连接时影响业务逻辑性能
func RegisteDeleteCache(channel string) {
	Go(func() {
		GetRedisSubPool().Subscribe(nil, DeleteCacheFromRedisMessage, channel)
	})
}

// DeleteCacheFromRedisMessage 根据redis订阅消息清空内存缓存
func DeleteCacheFromRedisMessage(msg *redis.Message) {
	GetCacheCli().DeleteFromLocalCache(msg.Payload)
	GetLogger().Info("remove local cache",
		zap.String("channel", msg.Channel),
		zap.String("cacheKey", msg.Payload),
	)
}

func initCacheCli() {
	cacheClient = cache.New(&cache.Options{
		//Redis:        GetRedisCli(),
		LocalCache:   cache.NewTinyLFU(10, time.Minute),
		StatsEnabled: true,
	})
	GetLogger().Info("memory cache ready...")
}
