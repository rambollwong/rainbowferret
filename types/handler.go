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

// SuccessStatuser is implemented by response types that want to control the
// success HTTP status code written by HandleT. When the response implements
// it, its return value replaces the default 200 OK (e.g. 201 Created).
//
// SuccessStatuser 由希望控制 HandleT 写出的成功状态码的响应类型实现。
// 当响应实现该接口时，其返回值替代默认的 200 OK（如 201 Created）。
type SuccessStatuser interface {
	SuccessStatus() int
}
