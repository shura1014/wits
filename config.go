package wits

import (
	"github.com/shura1014/cfg"
	"github.com/shura1014/common/env"
	"time"
)

var conf *cfg.Config

var defaultCfgDir = "config"

func init() {
	// dir: app.config.dir
	value, ok := env.GetEnv("app.config.dir")
	if !ok {
		value = defaultCfgDir
	}
	c, err := cfg.LoadConfig(value, "app.yaml", "app.yml", "app.properties")
	if err != nil || c == nil {
		return
	}
	conf = c
}

func Get(pattern string) (value any, err error) {
	return conf.Get(pattern)
}

func GetString(pattern string, def ...string) (value string) {
	return conf.GetString(pattern, def...)
}

func GetInt64(pattern string, def ...int64) (value int64) {
	return conf.GetInt64(pattern, def...)
}

func GetInt(pattern string, def ...int) (value int) {
	return conf.GetInt(pattern, def...)
}

func GetBool(pattern string, def ...bool) (value bool) {
	return conf.GetBool(pattern, def...)
}

func GetTime(pattern string, def ...time.Duration) (value time.Duration) {
	return conf.GetTime(pattern, def...)
}

func GetArray(pattern string) (value []any) {
	return conf.GetArray(pattern)
}

func GetStringMap(pattern string) (value map[string]any) {
	return conf.GetStringMap(pattern)
}

func GetAll() (data map[string]any, err error) {
	return conf.GetAll()
}

func GetPort() string {
	return GetString(AppPort)
}

func GetAddress() string {
	return GetString(AppAddress)
}

func GetAppName() string {
	return GetString(AppName)
}

// IsOpenCors 是否打开跨域
// 开发阶段打开
func IsOpenCors() bool {
	return GetBool(AppServerCorsEnable)
}
