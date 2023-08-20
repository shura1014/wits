package wits

import (
	"github.com/shura1014/logger"
	"net/http"
	"unicode"
)

type Map map[string]any

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func CreateTestContext(w http.ResponseWriter) (c *Context, r *Engine) {
	r = New()
	c = r.allocateContext()
	c.Reset(w, nil)
	return
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if address := GetAddress(); address != "" {
			Debug("The address is not set. Obtained through environment variables app.server.address: %s", address)
			Debug("The App will Running Listen Address is %s Get by Config or Environment Variables app.server.address.", address)
			return address
		} else if port := GetPort(); port != "" {
			address = ":" + port
			Debug("The address is not set. Obtained through environment variables app.server.port: %s", port)
			Debug("The App will Running Listen Address is %s Get by Config or Environment Variables app.server.port.", address)
			return address
		}
		Warn("The port is not set and no environment variables app.server.address or app.server.port are configured.")
		Debug("The App will Running Listen Address is :8888 Get by default.")
		return ":8888"
	case 1:
		Debug("The App will Running Listen Address is %s", logger.Blue(addr[0]))
		return addr[0]

	default:
		panic("too many address parameters")
	}
}

func assert(guard bool, text string) {
	if !guard {
		panic(text)
	}
}

func GetLogPrefix() string {
	prefix := GetString(AppLogPrefix)
	if prefix == "" {
		return DefaultLogPrefix
	}
	return prefix
}
