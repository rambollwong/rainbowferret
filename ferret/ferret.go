package ferret

import (
	"net/http"
)

// Router is the interface that HTTP routing groups must implement.
// Router 是 HTTP 路由分组必须实现的接口。
type Router interface {
	Handler() http.Handler

	Use(middleware ...Middleware)

	Handle(pattern string, h http.Handler, middleware ...Middleware)
	HandleFunc(pattern string, h http.HandlerFunc, middleware ...Middleware)

	Method(method, pattern string, h http.Handler, middleware ...Middleware)
	MethodFunc(method, pattern string, h http.HandlerFunc, middleware ...Middleware)

	Connect(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Delete(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Get(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Head(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Options(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Patch(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Post(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Put(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Trace(pattern string, h http.HandlerFunc, middleware ...Middleware)

	NotFound(h http.HandlerFunc)

	MethodNotAllowed(h http.HandlerFunc)
}
