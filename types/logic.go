package types

import (
	"context"
)

// Logic is the interface for generic request/response handlers.
// Logic 是泛型请求/响应处理器的接口。
type Logic[T any] interface {
	Execute(ctx context.Context, req T) (res any, err error)
}

// LogicFunc is the function variant of Logic.
// LogicFunc 是 Logic 的函数变体。
type LogicFunc[T any] func(ctx context.Context, req T) (res any, err error)
