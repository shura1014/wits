package wits

import (
	"github.com/shura1014/common/container/tree"
	"github.com/shura1014/logger"
	"net/http"
)

type FilterChain []*Filter

// Filter 过滤器
// 在接收到请求执行beforeFunc
// 请求完成之后执行afterFunc
type Filter struct {
	Order       int    // 排序，order较小的优先执行
	Include     string // 路径匹配，匹配到的才会执行 BeforeFunc and AfterFunc
	Exclude     string // 排除某个路径，即使include包含也不会执行
	includeTrie *tree.Node
	excludeTrie *tree.Node
	BeforeFunc  func(ctx *Context) bool
	AfterFunc   func(ctx *Context)
}

// Chain 实现beforeFunc先执行的afterFunc应该后执行
func (r Filter) Chain(handler func(ctx *Context)) func(ctx *Context) {
	return func(ctx *Context) {
		if r.BeforeFunc != nil {
			if !r.BeforeFunc(ctx) {
				return
			}
		}
		handler(ctx)
		if r.AfterFunc != nil {
			r.AfterFunc(ctx)
		}
	}
}

func (r FilterChain) Len() int {
	return len(r)
}

func (r FilterChain) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Less 比较，将order大的排到前面，实际上是最后执行
func (r FilterChain) Less(i, j int) bool {
	return r[i].Order > r[j].Order
}

// handler 路由的真正处理逻辑，不同的逻辑可以实现不同的功能，例如代理转发
func (r *FilterChain) doFilter(ctx *Context, handler RouterHandler) {
	defer func() {
		ctx.W.Flush()
	}()
	var (
		path = ctx.R.URL.Path
	)
	// handler之前 接收到连接的处理逻辑，校验请求的合法性，如ip白名单处理
	if r != nil {
		for _, filter := range *r {
			if filter.Exclude != "" {
				excludeNode := filter.excludeTrie.Get(path, nil)
				if excludeNode != nil && excludeNode.IsEnd {
					continue
				}
			}

			if filter.Include != "" {
				includeNode := filter.includeTrie.Get(path, nil)

				if includeNode != nil && includeNode.IsEnd {
					handler = filter.Chain(handler)
				}
			}
		}
	}
	handler(ctx)
	return
}

func IpCheck(IsBlock ...bool) *Filter {
	block := false
	if len(IsBlock) > 0 {
		block = true
	}
	return &Filter{
		Order:   0,
		Include: "/**",
		BeforeFunc: func(ctx *Context) bool {
			// 白名单处理
			ip := ctx.RemoteIP()

			check := ctx.e.acl.AclCheck(ip)
			// 检查到了，但是列表全是黑名单处理，那么 return false
			if ctx.e.acl.IsEnable() && check && block {
				_ = ctx.String(http.StatusUnauthorized, "Deny illegal access from %s", ip)
				return false
			}

			// 没有检查到，但是列表全是白名单处理，那么 return false
			if ctx.e.acl.IsEnable() && !check && !block {
				_ = ctx.String(http.StatusUnauthorized, "Deny illegal access from %s", ip)
				return false
			}

			return true
		},
	}
}

var RequestLog = &Filter{
	Order:   1,
	Include: "/**",
	BeforeFunc: func(ctx *Context) bool {
		ctx.DEBUG("received new connect clientIp: %s,method: %s, path: %s", logger.Blue(ctx.RemoteIP()), logger.Blue(ctx.GetMethod()), logger.Blue(ctx.R.URL.Path))
		return true
	},
}
