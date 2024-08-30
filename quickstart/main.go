package main

import (
	"github.com/shura1014/wits"
	"os"
	"time"
)

func main() {
	// 添加环境变量
	_ = os.Setenv(wits.AppRestResultEnable, "true")

	engine := wits.Default()
	group := engine.Group("api")
	group.GET("probe", func(ctx *wits.Context) {
		ctx.Success("ok", time.Now().UnixNano())
	})
	engine.RunServer()
}
