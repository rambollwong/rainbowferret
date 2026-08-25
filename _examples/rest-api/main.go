// Example: REST API using the generic Handler pattern for automatic JSON
// binding and response encoding.
//
// Route registration is split by Go version:
//   - routes_go1.27.go    — Go 1.27+ method-style registration
//   - routes_pre_go1.27.go — pre-Go 1.27 package-function registration
//
// 示例：使用泛型 Handler 模式实现自动 JSON 绑定和响应编码的 REST API。
//
// 路由注册按 Go 版本拆分：
//   - routes_go1.27.go     — Go 1.27+ 方法式注册
//   - routes_pre_go1.27.go — Go 1.27 之前的包级函数式注册
//
//	go run _examples/rest-api/main.go
package main

import (
	"net/http"

	"github.com/rambollwong/rainbowferret/ferret"
	"github.com/rambollwong/rainbowferret/middleware"
)

// CreateUserReq is the request type for creating a user.
// CreateUserReq 是创建用户的请求类型。
type CreateUserReq struct {
	Name  string `json:"name"  param:"name"`
	Email string `json:"email" param:"email"`
}

// UserResp is the response type representing a user.
// UserResp 是表示用户的响应类型。
type UserResp struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	sm := http.NewServeMux()
	root := ferret.NewRootGroup(sm,
		middleware.Logger(),
		middleware.RequestID(),
		middleware.ContentType("application/json"),
		middleware.ContentCharset("utf-8"),
	)

	api := root.Group("/api/v1")
	registerRoutes(api)

	http.ListenAndServe(":8080", root.Handler())
}
