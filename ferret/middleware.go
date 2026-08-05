package ferret

import "net/http"

// Middleware is the middleware type for http.Handler.
// Middleware 是用于 http.Handler 的中间件类型。
type Middleware func(next http.Handler) http.Handler

// MiddlewareFunc is the middleware type for http.HandlerFunc.
// MiddlewareFunc 是用于 http.HandlerFunc 的中间件类型。
type MiddlewareFunc func(next http.HandlerFunc) http.HandlerFunc

// Chain composes a list of middleware into a single Middleware. Middleware is
// executed in the order it is passed: the first element is the outermost
// (executed first on the way in, last on the way out).
//
// Chain 将一组中间件组合为单个 Middleware。中间件按传入顺序执行：
// 第一个元素是最外层（请求进入时最先执行，响应返回时最后执行）。
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
