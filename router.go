package wits

import (
	"github.com/shura1014/common/container/tree"
	"strings"
)

type router struct {
	groups       []*routerGroup
	defaultGroup *routerGroup

	// 全局拦截器
	middlewares []MiddlewareFunc
}

func (r *router) Group(name string) *routerGroup {
	prefix := strings.HasPrefix(name, slash)
	if !prefix {
		name = slash + name
	}
	g := &routerGroup{
		groupName: name,
		handleMap: make(map[string]map[string]HandlerFunc),
		treeNode: &tree.Node{
			Name:     name,
			Children: make([]*tree.Node, 0),
		},
		middlewares:          make([]MiddlewareFunc, 0),
		middlewaresMethodMap: make(map[string]map[string][]MiddlewareFunc),
	}
	// 有过有全局拦截器，得添加进去
	if r.middlewares != nil {
		g.Use(r.middlewares...)
	}
	if name == slash {
		r.defaultGroup = g
	}
	r.groups = append(r.groups, g)
	return g
}

func (r *router) getRouterGroup(path string) *routerGroup {
	for _, routerGroup := range r.groups {
		if strings.HasPrefix(path, routerGroup.groupName+slash) {
			return routerGroup
		}
	}
	return r.defaultGroup
}

// Use
// 全局中间件
// middlewares 中间件
func (r *router) Use(middlewares ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *router) UseFront(middlewares ...MiddlewareFunc) {
	newMiddlewareFunc := make([]MiddlewareFunc, 0)
	newMiddlewareFunc = append(newMiddlewareFunc, middlewares...)
	newMiddlewareFunc = append(newMiddlewareFunc, r.middlewares...)
	r.middlewares = newMiddlewareFunc
}
