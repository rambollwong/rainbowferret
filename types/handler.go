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
