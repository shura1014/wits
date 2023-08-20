package wits

import (
	"net/http"
	"os"
)

// Shutdown 发送终止信号
func (e *Engine) Shutdown() {
	quit <- os.Interrupt
}

func (e *Engine) preRun() {
	for _, handler := range e.PreRunFunc {
		handler()
	}
}
func (e *Engine) postNew() {
	for _, handler := range e.PostNewFunc {
		handler()
	}
}

// preRun 内置的一个preRun
// 如果需要使用: e.RegisterPreRunFunc(e.preRun)
func (e *Engine) defaultPreRun() {
	e.appName = GetAppName()
	// 检查session
	if enableSession {
		e.AddFilter(getSessionFilter())
	}
	// 统一结果返回
	if RestResultEnable() {
		e.SetHandlerBizError(commonHandler)
	}
	// 添加shutdown钩子的后置处理器
	e.RegistryServerPostProcessor(e.addShutdownHook)
}

func (e *Engine) defaultPostNew() {
	if enableGlobalTimeout {
		e.Use(Timeout())
	}

	if enableRecordCost {
		e.Use(costMiddleWareFunc)
	}
}

// server后置处理，注册一些shutdown钩子
func (e *Engine) addShutdownHook(server *http.Server) {
	for _, hook := range e.shutdownHook {
		server.RegisterOnShutdown(hook)
	}
}
