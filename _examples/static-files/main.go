// Example: serving static files from a local directory and embed.FS.
// 示例：从本地目录和 embed.FS 提供静态文件服务。
//
//	go run _examples/static-files/main.go
package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/rambollwong/rainbowferret/ferret"
	"github.com/rambollwong/rainbowferret/middleware"
)

//go:embed public
var publicFS embed.FS

func main() {
	sm := http.NewServeMux()
	root := ferret.NewRootGroup(sm,
		middleware.Logger(),
		middleware.RequestID(),
		middleware.Compress(),
	)

	// Serve from a local directory.
	// 从本地目录提供服务。
	root.Static("/assets", "./public")
	// → GET /assets/index.html

	// Serve from embedded files (Go 1.16+ embed.FS).
	// 从嵌入文件提供服务。
	sub, _ := fs.Sub(publicFS, "public")
	root.StaticFS("/embed", sub)
	// → GET /embed/index.html

	// Fallback route.
	root.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`Try <a href="/assets/">/assets/</a> or <a href="/embed/">/embed/</a>`))
	})

	http.ListenAndServe(":8080", root.Handler())
}
