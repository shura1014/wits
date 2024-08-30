package wits

import (
	"context"
	"net/http"
	"time"
)

var GlobalTimeout time.Duration
var enableGlobalTimeout bool

func init() {
	enableGlobalTimeout = GetBool(AppGlobalTimeoutEnable)
	GlobalTimeout = GetTime(AppGlobalTimeout, 5000) * time.Millisecond
}

// EnableGlobalTimeout 手动开启全局超时
func EnableGlobalTimeout() {
	enableGlobalTimeout = true
}

// SetGlobalTimeout 设置全局超时 毫秒
func SetGlobalTimeout(timeout time.Duration) {
	GlobalTimeout = timeout * time.Millisecond
}

// Timeout 超时中间件 timeout 毫秒
func Timeout(t ...time.Duration) MiddlewareFunc {
	timeout := GlobalTimeout
	if len(t) > 0 {
		timeout = t[0] * time.Millisecond
	}
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			c, cancelFunc := context.WithTimeout(context.Background(), timeout)
			defer cancelFunc()
			data := make(chan struct{}, 1)
			go func() {
				defer func() {
					close(data)
				}()
				next(ctx)
			}()
			select {
			case <-data:
			case <-c.Done():
				ctx.Fail(http.ErrHandlerTimeout.Error(), nil)
				return
			}
		}
	}
}
