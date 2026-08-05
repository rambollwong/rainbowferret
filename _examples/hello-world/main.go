// Example: minimal Hello World server.
// 示例：最小化的 Hello World 服务。
//
//	go run _examples/hello-world/main.go
package main

import (
	"net/http"

	"github.com/rambollwong/rainbowferret/ferret"
)

func main() {
	sm := http.NewServeMux()
	root := ferret.NewRootGroup(sm)

	root.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	root.Get("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		w.Write([]byte("Hello, " + name + "!"))
	})

	http.ListenAndServe(":8080", root.Handler())
}
