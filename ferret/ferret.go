package ferret

import (
	"io/fs"
	"net/http"
)

// Router is the interface that HTTP routing groups must implement.
// Router 是 HTTP 路由分组必须实现的接口。
type Router interface {
	Handler() http.Handler

	Use(middleware ...Middleware)

	Handle(pattern string, h http.Handler, middleware ...Middleware)
	HandleFunc(pattern string, h http.HandlerFunc, middleware ...Middleware)

	Method(method, pattern string, h http.Handler, middleware ...Middleware)
	MethodFunc(method, pattern string, h http.HandlerFunc, middleware ...Middleware)

	Connect(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Delete(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Get(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Head(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Options(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Patch(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Post(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Put(pattern string, h http.HandlerFunc, middleware ...Middleware)
	Trace(pattern string, h http.HandlerFunc, middleware ...Middleware)

	NotFound(h http.HandlerFunc)

	MethodNotAllowed(h http.HandlerFunc)

	// StaticFileHandler serves static files from a local directory.
	// The pattern is automatically suffixed with "/" for subtree matching.
	//
	// StaticFileHandler 从本地目录提供静态文件服务。
	// pattern 会自动添加 "/" 后缀以匹配子树。
	Static(pattern, dir string)

	// StaticFS serves static files from an fs.FS (e.g. embed.FS, os.DirFS).
	// StaticFS 从 fs.FS（如 embed.FS、os.DirFS）提供静态文件服务。
	StaticFS(pattern string, fsys fs.FS)
}
