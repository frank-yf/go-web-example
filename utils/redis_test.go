package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisPing(t *testing.T) {
	ctx := context.TODO()
	res, bn := PingRedis(ctx)
	assert.Truef(t, bn, "无法获取redis连接信息，ping返回信息：%s", res)
}
