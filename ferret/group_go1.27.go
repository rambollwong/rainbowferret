// This file provides Go 1.27+ method-style generic handler registration on
// *Group. Methods with their own type parameters require Go 1.27's generic
// methods support, hence the build constraint. The equivalent package-level
// functions (available on every Go version 1.18+) live in group.go.
//
// 本文件提供 Go 1.27+ 的 *Group 方法式泛型处理器注册。带独立类型参数的方法
// 需要 Go 1.27 新增的方法泛型支持，因此使用 build 约束。等价的全版本可用
// 包级函数（Go 1.18+）位于 group.go。
//
//go:build go1.27

package ferret

import (
	"github.com/rambollwong/rainbowferret/types"
	"github.com/rambollwong/rainbowferret/util"
)

// GetHandler registers a GET handler backed by a generic Handler interface.
// The request body is automatically decoded into T, and the response is
// automatically JSON-encoded.
//
// GetHandler 注册一个由泛型 Handler 接口支持的 GET 处理器。
// 请求体自动解码为 T，响应自动 JSON 编码。
func (g *Group) GetHandler[T, R any](pattern string, handler types.Handler[T, R], mws ...Middleware) {
	g.Get(pattern, util.HandleT(handler.Handle), mws...)
}

// PostHandler registers a POST handler backed by a generic Handler interface.
// The request body is automatically decoded into T, and the response is
// automatically JSON-encoded.
//
// PostHandler 注册一个由泛型 Handler 接口支持的 POST 处理器。
// 请求体自动解码为 T，响应自动 JSON 编码。
func (g *Group) PostHandler[T, R any](pattern string, handler types.Handler[T, R], mws ...Middleware) {
	g.Post(pattern, util.HandleT(handler.Handle), mws...)
}

// DeleteHandler registers a DELETE handler backed by a generic Handler interface.
//
// DeleteHandler 注册一个由泛型 Handler 接口支持的 DELETE 处理器。
func (g *Group) DeleteHandler[T, R any](pattern string, handler types.Handler[T, R], mws ...Middleware) {
	g.Delete(pattern, util.HandleT(handler.Handle), mws...)
}

// OptionsHandler registers an OPTIONS handler backed by a generic Handler interface.
//
// OptionsHandler 注册一个由泛型 Handler 接口支持的 OPTIONS 处理器。
func (g *Group) OptionsHandler[T, R any](pattern string, handler types.Handler[T, R], mws ...Middleware) {
	g.Options(pattern, util.HandleT(handler.Handle), mws...)
}

// PatchHandler registers a PATCH handler backed by a generic Handler interface.
//
// PatchHandler 注册一个由泛型 Handler 接口支持的 PATCH 处理器。
func (g *Group) PatchHandler[T, R any](pattern string, handler types.Handler[T, R], mws ...Middleware) {
	g.Patch(pattern, util.HandleT(handler.Handle), mws...)
}

// PutHandler registers a PUT handler backed by a generic Handler interface.
//
// PutHandler 注册一个由泛型 Handler 接口支持的 PUT 处理器。
func (g *Group) PutHandler[T, R any](pattern string, handler types.Handler[T, R], mws ...Middleware) {
	g.Put(pattern, util.HandleT(handler.Handle), mws...)
}

// MethodHandler registers a handler for the given HTTP method, backed by a
// generic Handler interface.
//
// MethodHandler 为指定 HTTP 方法注册一个由泛型 Handler 接口支持的处理器。
func (g *Group) MethodHandler[T, R any](method, pattern string, handler types.Handler[T, R], mws ...Middleware) {
	g.MethodFunc(method, pattern, util.HandleT(handler.Handle), mws...)
}
