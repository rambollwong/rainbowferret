// This file registers routes using the Go 1.27+ method-style API on *Group,
// where each handler is a struct implementing types.Handler[T, R] with a
// concrete request (T) and response (R) type.
//
// 本文件使用 Go 1.27+ 的 *Group 方法式 API 注册路由，
// 每个处理器是实现 types.Handler[T, R] 的结构体，请求（T）与响应（R）
// 均为具体类型。
//
//go:build go1.27

package main

import (
	"context"

	"github.com/rambollwong/rainbowferret/ferret"
)

// GetUserHandler handles GET /users/{id}.
// GetUserHandler 处理 GET /users/{id}。
type GetUserHandler struct{}

func (GetUserHandler) Handle(ctx context.Context, _ struct{}) (UserResp, error) {
	return UserResp{ID: "42", Name: "alice", Email: "alice@example.com"}, nil
}

// CreateUserHandler handles POST /users.
// CreateUserHandler 处理 POST /users。
type CreateUserHandler struct{}

func (CreateUserHandler) Handle(ctx context.Context, req CreateUserReq) (UserResp, error) {
	return UserResp{ID: "42", Name: req.Name, Email: req.Email}, nil
}

// DeleteUserHandler handles DELETE /users/{id}.
// DeleteUserHandler 处理 DELETE /users/{id}。
type DeleteUserHandler struct{}

func (DeleteUserHandler) Handle(ctx context.Context, _ struct{}) (*UserResp, error) {
	// Returning a nil response produces 204 No Content.
	// 返回 nil 响应产生 204 No Content。
	return nil, nil
}

// registerRoutes wires the handlers onto the group using method-style calls.
// registerRoutes 通过方法式调用将处理器注册到分组上。
func registerRoutes(api *ferret.Group) {
	api.GetHandler("/users/{id}", GetUserHandler{})
	api.PostHandler("/users", CreateUserHandler{})
	api.DeleteHandler("/users/{id}", DeleteUserHandler{})
}
