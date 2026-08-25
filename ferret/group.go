package ferret

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/rambollwong/rainbowferret/types"
	"github.com/rambollwong/rainbowferret/util"
)

var _ Router = (*Group)(nil)

// Group is an HTTP routing group that wraps http.ServeMux with middleware
// chaining, sub-grouping, and automatic 404/405 distinction. It implements
// the Router interface.
//
// Group 是一个 HTTP 路由分组，基于 http.ServeMux 封装了中间件链、子分组、
// 以及自动的 404/405 区分功能。它实现了 Router 接口。
type Group struct {
	sm          *http.ServeMux
	tree        *http.ServeMux                 // path-only mux for 405 detection
	routes      map[string]map[string]struct{} // path -> set of methods
	middlewares []Middleware
	prefix      string

	notFoundFn         http.HandlerFunc
	methodNotAllowedFn http.HandlerFunc

	mu sync.Mutex
}

// sentinel is a no-op handler registered on tree for path matching.
// sentinel 是一个空操作处理器，注册在 tree 上用于路径匹配。
var sentinel http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

// NewRootGroup creates a new root Group backed by the given ServeMux.
// The tree and routes are initialized fresh and shared with all child groups
// created via Group().
//
// NewRootGroup 创建一个新的根 Group，使用给定的 ServeMux 作为后端。
// tree 和 routes 会被全新初始化，并通过 Group() 共享给所有子分组。
func NewRootGroup(sm *http.ServeMux, middleware ...Middleware) *Group {
	return &Group{
		sm:          sm,
		tree:        http.NewServeMux(),
		routes:      make(map[string]map[string]struct{}),
		middlewares: middleware,
		prefix:      "",
	}
}

// Handler returns the http.Handler for this group. Middleware is applied
// per-route at registration time (via handle/Handle/HandleFunc), so the
// returned handler serves as the top-level dispatcher with 404/405 distinction
// using the internal tree and routes.
//
// Handler 返回该分组的 http.Handler。中间件在路由注册时逐个应用
// （通过 handle/Handle/HandleFunc），返回的 handler 作为顶层分发器，
// 通过内部的 tree 和 routes 实现 404/405 区分。
func (g *Group) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, pattern := g.sm.Handler(r)
		if pattern != "" {
			g.sm.ServeHTTP(w, r)
			return
		}

		// No match on main mux — check if the path itself is registered.
		// 主 mux 未匹配——检查该路径本身是否已注册。
		if _, treePattern := g.tree.Handler(r); treePattern != "" {
			if g.methodNotAllowedFn != nil {
				g.methodNotAllowedFn.ServeHTTP(w, r)
				return
			}
			g.writeMethodNotAllowed(w, treePattern)
			return
		}

		if g.notFoundFn != nil {
			g.notFoundFn.ServeHTTP(w, r)
			return
		}
		http.NotFound(w, r)
	})
}

// Group creates a child group under the given path prefix. Child groups
// inherit the parent's middleware and share the same ServeMux, tree, and
// routes so that 404/405 detection works across the entire group tree.
// Panics if prefix does not start with '/'.
//
// Group 在给定路径前缀下创建一个子分组。子分组继承父分组的中间件，
// 并共享同一个 ServeMux、tree 和 routes，使 404/405 检测在整个分组树中生效。
// 如果 prefix 不以 '/' 开头则会 panic。
func (g *Group) Group(prefix string, middleware ...Middleware) *Group {
	if !strings.HasPrefix(prefix, "/") {
		panic(fmt.Sprintf("ferret: group prefix must begin with '/' in '%s'", prefix))
	}
	prefix = strings.TrimSuffix(prefix, "/")

	mws := make([]Middleware, 0, len(g.middlewares)+len(middleware))
	mws = append(mws, append(g.middlewares, middleware...)...)
	return &Group{
		sm:          g.sm,
		tree:        g.tree,
		routes:      g.routes,
		middlewares: mws,
		prefix:      g.prefix + prefix,
	}
}

// Use appends global middleware to the group. Middleware is executed in the
// order it is added (first in, first executed).
//
// Use 向分组追加全局中间件。中间件按添加顺序执行（先加入的先执行）。
func (g *Group) Use(middleware ...Middleware) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.middlewares = append(g.middlewares, middleware...)
}

// Handle registers an http.Handler for the given pattern. The pattern may
// include an optional HTTP method prefix (e.g. "GET /foo"). Route-specific
// middleware runs after the global chain.
//
// Handle 为指定模式注册一个 http.Handler。模式可包含可选的 HTTP 方法前缀
// （如 "GET /foo"）。特定路由的中间件在全局链之后执行。
func (g *Group) Handle(pattern string, h http.Handler, middleware ...Middleware) {
	g.mu.Lock()
	defer g.mu.Unlock()
	middlewareChain := Chain(append(g.middlewares, middleware...)...)
	g.sm.Handle(pattern, middlewareChain(h))
}

// HandleFunc is the http.HandlerFunc variant of Handle.
//
// HandleFunc 是 Handle 的 http.HandlerFunc 变体。
func (g *Group) HandleFunc(pattern string, h http.HandlerFunc, middleware ...Middleware) {
	g.mu.Lock()
	defer g.mu.Unlock()
	middlewareChain := Chain(append(g.middlewares, middleware...)...)
	g.sm.Handle(pattern, middlewareChain(h))
}

// Method registers an http.Handler for the given HTTP method and path pattern.
//
// Method 为指定的 HTTP 方法和路径模式注册一个 http.Handler。
func (g *Group) Method(method, pattern string, handler http.Handler, mws ...Middleware) {
	g.handle(method, pattern, handler.ServeHTTP, mws...)
}

// MethodFunc is the http.HandlerFunc variant of Method.
//
// MethodFunc 是 Method 的 http.HandlerFunc 变体。
func (g *Group) MethodFunc(method, pattern string, handler http.HandlerFunc, mws ...Middleware) {
	g.handle(method, pattern, handler, mws...)
}

// Connect registers a CONNECT handler for the given path.
//
// Connect 为指定路径注册一个 CONNECT 处理器。
func (g *Group) Connect(pattern string, handler http.HandlerFunc, mws ...Middleware) {
	g.handle(http.MethodConnect, pattern, handler, mws...)
}

// Delete registers a DELETE handler for the given path.
//
// Delete 为指定路径注册一个 DELETE 处理器。
func (g *Group) Delete(pattern string, handler http.HandlerFunc, mws ...Middleware) {
	g.handle(http.MethodDelete, pattern, handler, mws...)
}

// Get registers a GET handler for the given path.
//
// Get 为指定路径注册一个 GET 处理器。
func (g *Group) Get(pattern string, handler http.HandlerFunc, mws ...Middleware) {
	g.handle(http.MethodGet, pattern, handler, mws...)
}

// Head registers a HEAD handler for the given path.
//
// Head 为指定路径注册一个 HEAD 处理器。
func (g *Group) Head(pattern string, handler http.HandlerFunc, mws ...Middleware) {
	g.handle(http.MethodHead, pattern, handler, mws...)
}

// Options registers an OPTIONS handler for the given path.
//
// Options 为指定路径注册一个 OPTIONS 处理器。
func (g *Group) Options(pattern string, handler http.HandlerFunc, mws ...Middleware) {
	g.handle(http.MethodOptions, pattern, handler, mws...)
}

// Patch registers a PATCH handler for the given path.
//
// Patch 为指定路径注册一个 PATCH 处理器。
func (g *Group) Patch(pattern string, handler http.HandlerFunc, mws ...Middleware) {
	g.handle(http.MethodPatch, pattern, handler, mws...)
}

// Post registers a POST handler for the given path.
//
// Post 为指定路径注册一个 POST 处理器。
func (g *Group) Post(pattern string, handler http.HandlerFunc, mws ...Middleware) {
	g.handle(http.MethodPost, pattern, handler, mws...)
}

// Put registers a PUT handler for the given path.
//
// Put 为指定路径注册一个 PUT 处理器。
func (g *Group) Put(pattern string, handler http.HandlerFunc, mws ...Middleware) {
	g.handle(http.MethodPut, pattern, handler, mws...)
}

// Trace registers a TRACE handler for the given path.
//
// Trace 为指定路径注册一个 TRACE 处理器。
func (g *Group) Trace(pattern string, handler http.HandlerFunc, mws ...Middleware) {
	g.handle(http.MethodTrace, pattern, handler, mws...)
}

// NotFound sets the handler called when no route matches the request path.
// If not set, http.NotFound is used as the fallback.
//
// NotFound 设置当请求路径未匹配到任何路由时调用的处理器。
// 若未设置，则使用 http.NotFound 作为兜底。
func (g *Group) NotFound(handler http.HandlerFunc) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.notFoundFn = handler
}

// MethodNotAllowed sets the handler called when the path matches a registered
// route but the HTTP method is not allowed. If set, it takes precedence over
// the default 405 response (which includes an Allow header).
//
// MethodNotAllowed 设置当路径匹配到已注册路由但 HTTP 方法不允许时调用的处理器。
// 若设置，它将覆盖默认的 405 响应（默认响应会包含 Allow 头）。
func (g *Group) MethodNotAllowed(handler http.HandlerFunc) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.methodNotAllowedFn = handler
}

// handle registers a method+path handler on the group. It records the route
// for 405 detection and Allow-header generation.
//
// handle 在分组上注册一个 method+path 处理器，
// 同时记录路由信息用于 405 检测和 Allow 头生成。
func (g *Group) handle(method, pattern string, handler http.HandlerFunc, mws ...Middleware) {
	if !strings.HasPrefix(pattern, "/") {
		pattern = "/" + pattern
	}

	fullPath := g.prefix + pattern
	fullPattern := method + " " + fullPath

	g.mu.Lock()
	defer g.mu.Unlock()

	middlewareChain := Chain(append(g.middlewares, mws...)...)
	g.sm.Handle(fullPattern, middlewareChain(handler))

	// Register on the path-only tree for 405 detection (once per path).
	// 在纯路径 tree 上注册（每个路径仅一次），用于 405 检测。
	if _, exists := g.routes[fullPath]; !exists {
		g.routes[fullPath] = make(map[string]struct{})
		g.tree.Handle(fullPath, sentinel)
	}
	g.routes[fullPath][method] = struct{}{}
}

// writeMethodNotAllowed writes a 405 Method Not Allowed response with an Allow
// header listing the methods registered for the matched path pattern.
//
// writeMethodNotAllowed 写入 405 Method Not Allowed 响应，
// 并在 Allow 头中列出该路径模式已注册的所有方法。
func (g *Group) writeMethodNotAllowed(w http.ResponseWriter, pattern string) {
	methods := g.routes[pattern]
	allowed := make([]string, 0, len(methods))
	for m := range methods {
		allowed = append(allowed, m)
	}
	sort.Strings(allowed)
	w.Header().Set("Allow", strings.Join(allowed, ", "))
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// Static serves static files from a local directory. The pattern is
// automatically suffixed with "/" for subtree matching. The handler is
// registered as a GET route so it participates in middleware chaining
// and 405 detection.
//
// Static 从本地目录提供静态文件服务。pattern 自动添加 "/" 后缀以匹配子树。
// handler 注册为 GET 路由，参与中间件链和 405 检测。
func (g *Group) Static(pattern, dir string) {
	g.StaticFS(pattern, os.DirFS(dir))
}

// StaticFS serves static files from an fs.FS (e.g. embed.FS, os.DirFS).
// The handler is registered as a GET route so it participates in middleware
// chaining and 405 detection.
//
// StaticFS 从 fs.FS（如 embed.FS、os.DirFS）提供静态文件服务。
// handler 注册为 GET 路由，参与中间件链和 405 检测。
func (g *Group) StaticFS(pattern string, fsys fs.FS) {
	if !strings.HasSuffix(pattern, "/") {
		pattern += "/"
	}
	fullPath := g.prefix + pattern
	handler := http.StripPrefix(fullPath, http.FileServerFS(fsys))
	g.handle(http.MethodGet, pattern, handler.ServeHTTP)
}

// GetHandler registers a GET handler backed by a generic HandlerFunc.
// The request body is automatically decoded into T, and the response is
// automatically JSON-encoded.
//
// GetHandler 注册一个由泛型 HandlerFunc 支持的 GET 处理器。
// 请求体自动解码为 T，响应自动 JSON 编码。
func GetHandler[T any](g *Group, pattern string, handlerFn types.HandlerFunc[T], mws ...Middleware) {
	g.Get(pattern, util.HandleT(handlerFn), mws...)
}

// PostHandler registers a POST handler backed by a generic HandlerFunc.
// The request body is automatically decoded into T, and the response is
// automatically JSON-encoded.
//
// PostHandler 注册一个由泛型 HandlerFunc 支持的 POST 处理器。
// 请求体自动解码为 T，响应自动 JSON 编码。
func PostHandler[T any](g *Group, pattern string, handlerFn types.HandlerFunc[T], mws ...Middleware) {
	g.Post(pattern, util.HandleT(handlerFn), mws...)
}

// DeleteHandler registers a DELETE handler backed by a generic HandlerFunc.
//
// DeleteHandler 注册一个由泛型 HandlerFunc 支持的 DELETE 处理器。
func DeleteHandler[T any](g *Group, pattern string, handlerFn types.HandlerFunc[T], mws ...Middleware) {
	g.Delete(pattern, util.HandleT(handlerFn), mws...)
}

// OptionsHandler registers an OPTIONS handler backed by a generic HandlerFunc.
//
// OptionsHandler 注册一个由泛型 HandlerFunc 支持的 OPTIONS 处理器。
func OptionsHandler[T any](g *Group, pattern string, handlerFn types.HandlerFunc[T], mws ...Middleware) {
	g.Options(pattern, util.HandleT(handlerFn), mws...)
}

// PatchHandler registers a PATCH handler backed by a generic HandlerFunc.
//
// PatchHandler 注册一个由泛型 HandlerFunc 支持的 PATCH 处理器。
func PatchHandler[T any](g *Group, pattern string, handlerFn types.HandlerFunc[T], mws ...Middleware) {
	g.Patch(pattern, util.HandleT(handlerFn), mws...)
}

// PutHandler registers a PUT handler backed by a generic HandlerFunc.
//
// PutHandler 注册一个由泛型 HandlerFunc 支持的 PUT 处理器。
func PutHandler[T any](g *Group, pattern string, handlerFn types.HandlerFunc[T], mws ...Middleware) {
	g.Put(pattern, util.HandleT(handlerFn), mws...)
}

// MethodHandler registers a handler for the given HTTP method, backed by a
// generic HandlerFunc.
//
// MethodHandler 为指定 HTTP 方法注册一个由泛型 HandlerFunc 支持的处理器。
func MethodHandler[T any](g *Group, method, pattern string, handlerFn types.HandlerFunc[T], mws ...Middleware) {
	g.MethodFunc(method, pattern, util.HandleT(handlerFn), mws...)
}
