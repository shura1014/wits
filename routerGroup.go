package wits

import (
	"github.com/shura1014/common/container/tree"
	"net/http"
	"strings"
)

type HandlerFunc func(ctx *Context)

type routerGroup struct {
	// 路由组名
	groupName string
	// 处理器 多级支持 /refresh -> POST -> 处理器
	// 处理器 多级支持 /refresh -> GET -> 处理器
	handleMap map[string]map[string]HandlerFunc

	treeNode *tree.Node

	// 路由组级别的通用中间价
	middlewares []MiddlewareFunc

	// 单个处理器的中间价
	middlewaresMethodMap map[string]map[string][]MiddlewareFunc
}

// @name /refresh
// @method POST
// @handler 处理器
// @middlewares 单个处理器的中间价
func (r *routerGroup) add(routerName string, method string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	if !strings.HasPrefix(routerName, "/") {
		routerName = "/" + routerName
	}
	_, ok := r.handleMap[routerName]
	if !ok {
		r.handleMap[routerName] = make(map[string]HandlerFunc)
		r.middlewaresMethodMap[routerName] = make(map[string][]MiddlewareFunc)
	}

	_, ok = r.handleMap[routerName][method]
	if ok {
		Debug("有重复的路由 name %s method %s", routerName, method)
	}

	r.handleMap[routerName][method] = handler
	r.middlewaresMethodMap[routerName][method] = append(r.middlewaresMethodMap[routerName][method], middlewares...)
	if IsOpenCors() {
		r.handleMap[routerName][http.MethodOptions] = handler
		r.middlewaresMethodMap[routerName][http.MethodOptions] = append(r.middlewaresMethodMap[routerName][method], middlewares...)
	}

	r.treeNode.Put(routerName)
}

// Use methodHandle
// 处理器执行逻辑，可以在前后做增强
// @handlerFunc 真正需要执行的业务逻辑
// @ctx 上下文
func (r *routerGroup) methodHandle(routerName string, method string, handlerFunc HandlerFunc, ctx *Context) {
	middlewares := r.middlewares
	if middlewares != nil {
		for _, middleware := range middlewares {
			// 包装函数做前置增强
			handlerFunc = middleware(handlerFunc)
		}
	}

	// 单个处理器的中间价优先处理
	methodMiddlewares := r.middlewaresMethodMap[routerName][method]
	if middlewares != nil {
		for _, middleware := range methodMiddlewares {
			// 包装函数做前置增强
			handlerFunc = middleware(handlerFunc)
		}
	}
	handlerFunc(ctx)
}

// Use
// 中间件使用
// middlewares 中间件
func (r *routerGroup) Use(middlewares ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}
