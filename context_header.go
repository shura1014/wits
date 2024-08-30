package wits

import (
	"net"
	"strings"
)

// RequestHeader
/**************头部操作****************************/
func (c *Context) RequestHeader(key string) string {
	return c.R.Header.Get(key)
}

func (c *Context) ContentType() string {
	return c.R.Header.Get(HeaderContentType)
}

func (c *Context) GetMethod() string {
	return c.R.Method
}

func (c *Context) GetPath() string {
	return c.R.URL.Path
}

func (c *Context) Header(key, value string) {
	if value == "" {
		c.W.Header().Del(key)
		return
	}
	c.W.Header().Set(key, value)
}

func (c *Context) UserAgent() string {
	return c.RequestHeader(HeaderUserAgent)
}

func (c *Context) Referer() string {
	return c.RequestHeader(HeaderReferer)
}

// RemoteIP 获取客户端ip
func (c *Context) RemoteIP() string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(c.R.RemoteAddr))
	if err != nil {
		return ""
	}
	return ip
}

func (c *Context) ClientIP() string {
	var clientIP string

	// 检查该ip是不是本地ip 或者说它是一个代理ip，因此需要找出原ip
	proxyIps := c.RequestHeader(XForwardedFor)
	if proxyIps != "" {
		// XForwardedFor可能经过多重代理，取第一个就行
		i := strings.IndexAny(proxyIps, ",")
		if i > 0 {
			clientIP = strings.TrimSpace(proxyIps[:i])
		}
		clientIP = strings.TrimPrefix(clientIP, "[")
		clientIP = strings.TrimSuffix(clientIP, "]")
		return clientIP
	}

	clientIP = c.RequestHeader(XRealIP)
	if clientIP == "" {
		clientIP = c.RemoteIP()
	}

	clientIP = strings.TrimPrefix(clientIP, "[")
	clientIP = strings.TrimSuffix(clientIP, "]")
	return clientIP
}

/**************头部操作****************************/
