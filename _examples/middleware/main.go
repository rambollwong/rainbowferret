// Example: combining multiple middleware for a production-ready service.
// 示例：组合多个中间件构建生产级服务。
//
//	go run _examples/middleware/main.go
package main

import (
	"compress/gzip"
	"net/http"
	"time"

	"github.com/rambollwong/rainbowferret/ferret"
	"github.com/rambollwong/rainbowferret/middleware"
)

func main() {
	sm := http.NewServeMux()
	root := ferret.NewRootGroup(sm,
		middleware.Logger(),    // access logging
		middleware.Recoverer,   // panic recovery
		middleware.RequestID(), // unique request id
		middleware.ClientIPWithConfig(middleware.ClientIPConfig{ // resolve client IP
			TrustProxy: true,
		}),
		middleware.Compress(gzip.BestSpeed), // compress responses
		middleware.Timeout(5*time.Second),   // request timeout
	)

	// User endpoints — apply content-type check only to this group.
	// 用户端点 — 仅对该分组应用 Content-Type 检查。
	users := root.Group("/users", middleware.ContentType("application/json"))
	users.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"id":"` + r.PathValue("id") + `","name":"alice"}`))
	})
	users.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"ok":true}`))
	})

	// Slow endpoint to demonstrate timeout.
	// 慢端点用于演示超时。
	root.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		select {
		case <-ctx.Done():
			return
		case <-time.After(10 * time.Second):
			w.Write([]byte("done"))
		}
	})

	// Panic endpoint to demonstrate recovery.
	// 崩溃端点用于演示恢复。
	root.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	http.ListenAndServe(":8080", root.Handler())
}
