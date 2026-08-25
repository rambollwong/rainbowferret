// Example: REST API using the generic Handler pattern for automatic JSON
// binding and response encoding.
// 示例：使用泛型 Handler 模式实现自动 JSON 绑定和响应编码的 REST API。
//
//	go run _examples/rest-api/main.go
package main

import (
	"context"
	"net/http"

	"github.com/rambollwong/rainbowferret/ferret"
	"github.com/rambollwong/rainbowferret/middleware"
)

// --- request / response types -----------------------------------------------

type CreateUserReq struct {
	Name  string `json:"name"  param:"name"`
	Email string `json:"email" param:"email"`
}

type UserResp struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// --- handlers ---------------------------------------------------------------

func createUser(ctx context.Context, req CreateUserReq) (any, error) {
	return UserResp{ID: "42", Name: req.Name, Email: req.Email}, nil
}

func getUser(ctx context.Context, req struct{}) (any, error) {
	return UserResp{ID: "42", Name: "alice", Email: "alice@example.com"}, nil
}

func deleteUser(ctx context.Context, req struct{}) (any, error) {
	return nil, nil // 204 No Content
}

// --- main -------------------------------------------------------------------

func main() {
	sm := http.NewServeMux()
	root := ferret.NewRootGroup(sm,
		middleware.Logger(),
		middleware.RequestID(),
		middleware.ContentType("application/json"),
		middleware.ContentCharset("utf-8"),
	)

	api := root.Group("/api/v1")

	// Use generic Handler helpers — request body is auto-decoded into T,
	// response is auto-encoded as JSON.
	// 使用泛型 Handler 辅助函数 — 请求体自动解码为 T，响应自动 JSON 编码。
	ferret.GetHandler(api, "/users/{id}", getUser)
	ferret.PostHandler(api, "/users", createUser)
	ferret.DeleteHandler(api, "/users/{id}", deleteUser)

	http.ListenAndServe(":8080", root.Handler())
}
