// Example: sub-group routing with shared and route-specific middleware.
// 示例：使用共享中间件和路由专属中间件的子分组路由。
//
//	go run _examples/sub-group/main.go
package main

import (
	"net/http"

	"github.com/rambollwong/rainbowferret/ferret"
	"github.com/rambollwong/rainbowferret/middleware"
)

func main() {
	sm := http.NewServeMux()

	// Root group with a global logger.
	// 根分组带有全局日志中间件。
	root := ferret.NewRootGroup(sm, middleware.Logger())

	// Public routes — no extra middleware.
	// 公开路由 — 无额外中间件。
	root.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// API sub-group with a request timeout.
	// API 子分组带有请求超时中间件。
	api := root.Group("/api", middleware.Timeout(5*1e9)) // 5 seconds

	// V1 sub-group inherits both Logger and Timeout.
	// V1 子分组继承 Logger 和 Timeout。
	v1 := api.Group("/v1")
	v1.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("[user1, user2]"))
	})
	v1.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":"` + r.PathValue("id") + `","name":"alice"}`))
	})

	// V2 sub-group with its own request-id middleware on top.
	// V2 子分组额外叠加 request-id 中间件。
	v2 := api.Group("/v2", middleware.RequestID())
	v2.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("[userA, userB]"))
	})

	http.ListenAndServe(":8080", root.Handler())
}
