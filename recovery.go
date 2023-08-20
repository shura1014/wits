package wits

import (
	"github.com/shura1014/common/goerr"
)

func Recovery(next HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				bizError := goerr.Wrap(err)
				ctx.HandlerError(bizError)
			}
		}()
		next(ctx)
	}
}
