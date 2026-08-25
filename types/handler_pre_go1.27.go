package types

import (
	"context"
)

// Handler is the interface for generic request/response handlers.
// Handler 是泛型请求/响应处理器的接口。
type Handler[T any] interface {
	Handle(ctx context.Context, req T) (res any, err error)
}

// HandlerFunc is the function variant of Handler.
// HandlerFunc 是 Handler 的函数变体。
type HandlerFunc[T any] func(ctx context.Context, req T) (res any, err error)

// Handle implements the Handler interface by calling the underlying function.
// Handle 通过调用底层函数实现 Handler 接口。
func (f HandlerFunc[T]) Handle(ctx context.Context, req T) (res any, err error) {
	return f(ctx, req)
}
