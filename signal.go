package wits

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var quit = make(chan os.Signal)

func (e *Engine) handlerSignal() {
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	signal.Notify(quit,
		syscall.SIGINT,  // 键盘的中断
		syscall.SIGQUIT, // 键盘的退出
		syscall.SIGKILL, // 杀死
		syscall.SIGTERM, // kill 15 软件终止信号 优雅停机处理
		syscall.SIGUSR1, // 用户自定义信号一
		syscall.SIGUSR2, // 用户自定义的信号2
	)
	var sig os.Signal
	for {
		sig = <-quit
		Info("Receive a signal %d", sig)
		switch sig {
		case syscall.SIGTERM:
			e.gracefullyShutdown()
			return
		case syscall.SIGKILL, syscall.SIGINT, syscall.SIGQUIT:
			Info("Server exiting the pid is %d", os.Getpid())
			return
		case syscall.SIGUSR1:
			e.reload()
		default:
		}
	}

}

// 优雅停机 Connection refused
func (e *Engine) gracefullyShutdown() {
	gracefulTime := GetTime(appServerGracefulTime, 10)
	Info("The application will terminate in %d seconds", gracefulTime)
	ctx, cancel := context.WithTimeout(context.Background(), gracefulTime*time.Second)
	defer cancel()
	pid := os.Getpid()
	// 如果一个连接都没有，那么会立刻Shutdown
	if err := e.srv.Shutdown(ctx); err != nil {
		Error("Server Shutdown: %v the pid is %d", err, pid)
	}
	Info("Server exiting the pid is %d", pid)
}
