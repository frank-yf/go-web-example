// 在程序中应该避免使用野生的 goroutine，设定该工具类提供统一的 goroutine 创建方式

package utils

import (
	"context"

	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
)

// Go 统一的 goroutine 创建，避免因为 panic 导致主进程退出
func Go(f func()) {
	go func() {
		defer func() {
			res := recover()
			if res != nil {
				GetLogger().S.Error("panic: %+v", res)
			}
		}()
		f()
	}()
}

// GoWithGroup 使用 errgroup 创建 goroutine
func GoWithGroup(f func() error) (func() error, context.Context) {
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() (err error) {
		defer func() {
			res := recover()
			if res != nil {
				GetLogger().S.Error("panic: %+v", res)
			}
			if e, ok := res.(error); ok {
				err = multierr.Append(err, e)
			}
		}()
		err = f()
		return
	})
	return g.Wait, ctx
}
