// This file registers routes using the package-function API, which works on
// every Go version (1.18+). Handlers use the concrete response type R instead
// of any for compile-time type safety. Go 1.27+ projects may prefer the
// method-style API in routes_go1.27.go.
//
// 本文件使用包级函数 API 注册路由，该 API 在所有 Go 版本（1.18+）均可用。
// 处理器使用具体响应类型 R 而非 any，以获得编译期类型安全。
// Go 1.27+ 项目亦可优先使用 routes_go1.27.go 中的方法式 API。
//
//go:build !go1.27

package main

import (
	"context"

	"github.com/rambollwong/rainbowferret/ferret"
)

func getUser(ctx context.Context, _ struct{}) (UserResp, error) {
	return UserResp{ID: "42", Name: "alice", Email: "alice@example.com"}, nil
}

func createUser(ctx context.Context, req CreateUserReq) (UserResp, error) {
	return UserResp{ID: "42", Name: req.Name, Email: req.Email}, nil
}

func deleteUser(ctx context.Context, _ struct{}) (*UserResp, error) {
	return nil, nil // 204 No Content
}

// registerRoutes wires the handler functions onto the group using
// package-level helper functions.
//
// registerRoutes 通过包级辅助函数将处理器函数注册到分组上。
func registerRoutes(api *ferret.Group) {
	ferret.GetHandler(api, "/users/{id}", getUser)
	ferret.PostHandler(api, "/users", createUser)
	ferret.DeleteHandler(api, "/users/{id}", deleteUser)
}
