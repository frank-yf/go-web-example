package utils

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	redisSubPool     *redisSubscriptionPool
	redisSubPoolOnce sync.Once
)

// GetRedisSubPool 获取redis订阅连接池
func GetRedisSubPool() *redisSubscriptionPool {
	redisSubPoolOnce.Do(func() {
		redisSubPool = &redisSubscriptionPool{
			m:      new(sync.Map),
			length: new(int32),
			ctx:    GetRedisCli().Context(),
		}
		GetLogger().Info("redis subscription pool ready...")
	})
	return redisSubPool
}

// redisSubscriptionPool redis订阅连接池
type redisSubscriptionPool struct {
	m      *sync.Map
	length *int32

	// ctx 连接池中全局设定的context，当订阅事件没有设定例外的context，则以此为准
	ctx context.Context
}

func (p redisSubscriptionPool) withContext(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	return p.ctx
}

// Subscribe 注册订阅事件
func (p redisSubscriptionPool) Subscribe(ctx context.Context, consume func(*redis.Message), channel string) (loaded bool) {
	pubSub := GetRedisCli().Subscribe(p.withContext(ctx), channel)
	_, loaded = p.m.LoadOrStore(channel, pubSub)
	if !loaded {
		// 订阅事件不存在，则订阅到redis客户端
		Go(func() {
			ch := pubSub.Channel()
			for msg := range ch {
				consume(msg)
				GetLogger().Debug("consume subscribe message", zap.String("message", msg.String()))
			}
			GetLogger().Debug("subscribe channel is stopped", zap.String("channel", channel))
		})
		atomic.AddInt32(p.length, 1)
		GetLogger().Info("registe redis subscribe", zap.String("channel", channel))
	} else {
		// 订阅事件已存在，将新建的订阅连接关闭掉，避免连接逃逸
		err := pubSub.Close()
		if err != nil {
			GetLogger().Error("close escape redis subscription connection error", zap.Error(err))
		}
	}
	return
}

// Unsubscribe 取消订阅事件
// 如果订阅通道存在，则取消订阅事件；假如不存在则loaded会返回false
func (p redisSubscriptionPool) Unsubscribe(ctx context.Context, channel string) (loaded bool, err error) {
	v, loaded := p.m.LoadAndDelete(channel)
	if loaded {
		pubSub := v.(*redis.PubSub)
		err = p.unsubscribe(ctx, channel, pubSub)
	}
	return
}

// Close 关闭订阅连接池，取消已订阅的所有事件
// 当遇到无法正确关闭的订阅连接时，之后的所有订阅连接将不会被关闭，操作失败
func (p *redisSubscriptionPool) Close() (err error) {
	p.m.Range(func(k, v interface{}) bool {
		err = p.unsubscribe(nil, k.(string), v.(*redis.PubSub))
		suc := err == nil
		if suc {
			p.m.Delete(k)
		}
		return suc
	})
	if err == nil {
		GetLogger().Debug("redis subscription pool closed")
	}
	return
}

// unsubscribe 取消订阅事件
// 假如抛出异常即表示订阅事件未能成功取消，redis连接没有被释放
func (p redisSubscriptionPool) unsubscribe(ctx context.Context, channel string, pubSub *redis.PubSub) (err error) {
	// 从redis客户端取消订阅
	err = pubSub.Unsubscribe(p.withContext(ctx), channel)
	if err != nil {
		return
	}
	GetLogger().Debug("subscribe is canceled", zap.String("channel", channel))

	// 关闭订阅连接
	err = pubSub.Close()
	if err != nil {
		return
	}
	GetLogger().Debug("subscribe connection is closed", zap.String("channel", channel))

	// 池中订阅数递减
	atomic.AddInt32(p.length, -1)
	return
}

// Lookup 查看已订阅的通道名称
func (p redisSubscriptionPool) Lookup() (channels []string) {
	channels = make([]string, 0, p.Len())
	p.m.Range(func(k, _ interface{}) bool {
		channels = append(channels, k.(string))
		return true
	})
	return
}

// Len 订阅池中已有的订阅连接
func (p redisSubscriptionPool) Len() int {
	i := atomic.LoadInt32(p.length)
	return int(i)
}
