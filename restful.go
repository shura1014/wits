package wits

import "net/http"

const ANY = "ANY"

func (r *routerGroup) ANY(name string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	r.add(name, ANY, handlerFunc, middlewares...)
}

func (r *routerGroup) POST(name string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	r.add(name, http.MethodPost, handlerFunc, middlewares...)
}

func (r *routerGroup) GET(name string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	r.add(name, http.MethodGet, handlerFunc, middlewares...)
}

func (r *routerGroup) DELETE(name string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	r.add(name, http.MethodDelete, handlerFunc, middlewares...)
}
func (r *routerGroup) PUT(name string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	r.add(name, http.MethodPut, handlerFunc, middlewares...)
}
func (r *routerGroup) PATCH(name string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	r.add(name, http.MethodPatch, handlerFunc, middlewares...)
}
func (r *routerGroup) OPTIONS(name string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	r.add(name, http.MethodOptions, handlerFunc, middlewares...)
}
func (r *routerGroup) HEAD(name string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	r.add(name, http.MethodHead, handlerFunc, middlewares...)
}
