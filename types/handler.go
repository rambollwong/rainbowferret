package types

import (
	"context"
)

// Handler is the interface for generic request/response handlers with
// separate request (T) and response (R) type parameters.
//
// Handler 是泛型请求/响应处理器的接口，请求（T）与响应（R）类型参数分离。
type Handler[T, R any] interface {
	Handle(ctx context.Context, req T) (res R, err error)
}

// HandlerFunc is the function variant of Handler.
// HandlerFunc 是 Handler 的函数变体。
type HandlerFunc[T, R any] func(ctx context.Context, req T) (res R, err error)

// Handle implements the Handler interface by calling the underlying function.
// Handle 通过调用底层函数实现 Handler 接口。
func (f HandlerFunc[T, R]) Handle(ctx context.Context, req T) (res R, err error) {
	return f(ctx, req)
}

// SuccessStatuser is implemented by request types that want to control the
// success HTTP status code written for a non-nil response. When the bound
// request implements it, its return value replaces the default 200 OK.
//
// SuccessStatuser 由希望控制非 nil 响应成功状态码的请求类型实现。
// 当绑定后的请求实现该接口时，其返回值替代默认的 200 OK。
type SuccessStatuser interface {
	SuccessStatus() int
}
