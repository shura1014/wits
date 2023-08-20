package wits

import (
	"time"
)

type MiddlewareFunc func(handlerFunc HandlerFunc) HandlerFunc

var enableRecordCost bool

func init() {
	enableRecordCost = GetBool(AppRecordCostEnable)
}

// 请求耗时，全局拦截器
func costMiddleWareFunc(next HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		start := time.Now()
		next(ctx)
		ctx.DEBUG("path %s cost %s", ctx.R.RequestURI, time.Now().Sub(start))
	}

}
