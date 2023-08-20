package main

import (
	"github.com/shura1014/common/utils/timeutil"
	"github.com/shura1014/wits"
)

func main() {
	engine := wits.Default()
	group := engine.Group("api")
	group.GET("probe", func(ctx *wits.Context) {
		_ = ctx.Success(timeutil.Now())
	})
	engine.RunServer()
}
