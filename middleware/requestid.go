package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// requestIDKey is the context key under which the request ID is stored.
// requestIDKey 是存储请求 ID 的 context 键。
const requestIDKey ctxKey = "RequestID"

// RequestIDHeader is the default HTTP header name for the request ID.
// RequestIDHeader 是请求 ID 的默认 HTTP 头名。
const RequestIDHeader = "X-Request-ID"

// GetRequestID returns the request ID stored in ctx, or "" if none is found.
// GetRequestID 返回存储在 ctx 中的请求 ID，未找到则返回 ""。
func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

// WithRequestID returns a child context with the given request ID attached.
// WithRequestID 返回一个附加了指定请求 ID 的子 context。
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestID returns middleware that ensures every request has a unique request
// ID. It reads the ID from the RequestIDHeader when present (to support
// distributed trace propagation); otherwise it generates a new one prefixed
// with the machine hostname to avoid collisions in clustered deployments.
//
// Generated IDs look like: node3-ebaa9bf0e27166f24aa0cb0052dfca92
// (hostname prefix + "-" + 32 hex chars from crypto/rand).
//
// RequestID 返回一个确保每个请求都有唯一 ID 的中间件。
// 当请求头中存在 ID 时直接复用（支持分布式追踪透传）；
// 否则生成带本机 hostname 前缀的新 ID，避免集群部署时的碰撞。
//
// 生成的 ID 格式：node3-ebaa9bf0e27166f24aa0cb0052dfca92
// （hostname 前缀 + "-" + 32 位 hex 随机串）。
func RequestID() func(next http.Handler) http.Handler {
	return RequestIDWithConfig(RequestIDConfig{})
}

// RequestIDConfig holds the configuration for the RequestID middleware.
// RequestIDConfig 存储 RequestID 中间件的配置。
type RequestIDConfig struct {
	// HeaderName is the HTTP header used to read / write the request ID.
	// Defaults to RequestIDHeader ("X-Request-ID").
	// HeaderName 是用于读取/写入请求 ID 的 HTTP 头，默认为 "X-Request-ID"。
	HeaderName string

	// Prefix is prepended to every generated ID to help identify the source
	// node in a cluster. When left empty, os.Hostname is used automatically.
	// Set it to an explicit value to override the automatic hostname, or to
	// a single space " " to disable the prefix entirely.
	//
	// Prefix 会添加到每个生成的 ID 前面，用于在集群中识别来源节点。
	// 为空时自动使用 os.Hostname。可设置为显式值覆盖自动检测的主机名，
	// 或设为单个空格 " " 以完全禁用前缀。
	Prefix string

	// Generator creates a new request ID. When nil, a crypto/rand-based
	// hex generator (with prefix) is used.
	// Generator 创建一个新的请求 ID。为 nil 时使用基于 crypto/rand 的 hex 生成器（带前缀）。
	Generator func() string
}

// RequestIDWithConfig returns a RequestID middleware configured with cfg.
// RequestIDWithConfig 返回使用 cfg 配置的 RequestID 中间件。
func RequestIDWithConfig(cfg RequestIDConfig) func(next http.Handler) http.Handler {
	headerName := cfg.HeaderName
	if headerName == "" {
		headerName = RequestIDHeader
	}
	gen := cfg.Generator
	if gen == nil {
		prefix := cfg.Prefix
		switch prefix {
		case "":
			prefix = cachedHostname()
		case " ":
			prefix = ""
		}
		gen = func() string { return newID(prefix) }
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(headerName)
			if id == "" {
				id = gen()
			}

			ctx := context.WithValue(r.Context(), requestIDKey, id)
			w.Header().Set(headerName, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// --- internal ----------------------------------------------------------------

// newID returns a request ID with the given prefix. The random part is 16
// bytes from crypto/rand encoded as 32 hex characters.
//
// newID 返回带指定前缀的请求 ID。随机部分为 16 字节 crypto/rand，hex 编码为 32 字符。
func newID(prefix string) string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		if prefix != "" {
			return prefix + "-fallback"
		}
		return "fallback-request-id"
	}
	buf := make([]byte, 32)
	hex.Encode(buf, b)
	if prefix == "" {
		return string(buf)
	}
	return prefix + "-" + string(buf)
}
