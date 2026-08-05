package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
)

// clientIPKey is the context key under which the resolved client IP is stored.
// Uses the same ctxKey type declared in middleware.go for consistency.
//
// clientIPKey 是存储解析后的客户端 IP 的 context 键。
// 使用 middleware.go 中声明的 ctxKey 类型以保持一致。
const clientIPKey ctxKey = "ClientIP"

// GetClientIP returns the client IP stored in ctx, or "" if none is found.
// GetClientIP 返回存储在 ctx 中的客户端 IP，未找到则返回 ""。
func GetClientIP(ctx context.Context) string {
	ip, _ := ctx.Value(clientIPKey).(string)
	return ip
}

// WithClientIP returns a child context with the given client IP attached.
// WithClientIP 返回一个附加了指定客户端 IP 的子 context。
func WithClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, clientIPKey, ip)
}

// ClientIP returns middleware that resolves the client IP for every request.
// By default it uses r.RemoteAddr (stripping the port). When TrustProxy is
// enabled it also inspects the X-Forwarded-For and X-Real-IP headers.
//
// ClientIP 返回一个为每个请求解析客户端 IP 的中间件。
// 默认使用 r.RemoteAddr（去掉端口）。当 TrustProxy 启用时，
// 还会检查 X-Forwarded-For 和 X-Real-IP 头。
func ClientIP() func(next http.Handler) http.Handler {
	return ClientIPWithConfig(ClientIPConfig{})
}

// ClientIPConfig holds the configuration for the ClientIP middleware.
// ClientIPConfig 存储 ClientIP 中间件的配置。
type ClientIPConfig struct {
	// TrustProxy enables parsing of proxy headers (X-Forwarded-For and
	// X-Real-IP) to obtain the real client IP. Enable it only when your
	// service runs behind a trusted reverse proxy; otherwise a client
	// could spoof its IP.
	//
	// TrustProxy 启用解析代理头（X-Forwarded-For 和 X-Real-IP）获取真实
	// 客户端 IP。仅在服务运行于可信反向代理之后时启用；否则客户端可伪造 IP。
	TrustProxy bool

	// RealIPHeader overrides the header name used when TrustProxy is
	// enabled. Defaults to "X-Real-IP".
	// RealIPHeader 覆盖 TrustProxy 启用时使用的真实 IP 头名，
	// 默认为 "X-Real-IP"。
	RealIPHeader string
}

// ClientIPWithConfig returns a ClientIP middleware configured with cfg.
// ClientIPWithConfig 返回使用 cfg 配置的 ClientIP 中间件。
func ClientIPWithConfig(cfg ClientIPConfig) func(next http.Handler) http.Handler {
	trustProxy := cfg.TrustProxy
	realIPHeader := cfg.RealIPHeader
	if realIPHeader == "" {
		realIPHeader = "X-Real-IP"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := resolveIP(r, trustProxy, realIPHeader)
			ctx := context.WithValue(r.Context(), clientIPKey, ip)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// resolveIP extracts the best-guess client IP from the request.
//
// Without TrustProxy the result is always derived from r.RemoteAddr.
// With TrustProxy the precedence is:
//   1. X-Forwarded-For (first entry)
//   2. RealIPHeader (default X-Real-IP)
//   3. r.RemoteAddr
//
// resolveIP 从请求中提取最佳猜测的客户端 IP。
//
// 未启用 TrustProxy 时始终从 r.RemoteAddr 获取。
// 启用 TrustProxy 时优先级为：
//   1. X-Forwarded-For（第一个条目）
//   2. RealIPHeader（默认 X-Real-IP）
//   3. r.RemoteAddr
func resolveIP(r *http.Request, trustProxy bool, realIPHeader string) string {
	if trustProxy {
		if ip := firstForwardedFor(r); ip != "" {
			return ip
		}
		if ip := r.Header.Get(realIPHeader); ip != "" {
			return ip
		}
	}
	return stripPort(r.RemoteAddr)
}

// firstForwardedFor returns the first IP in the X-Forwarded-For chain,
// trimming whitespace. An empty string means the header is absent or empty.
//
// firstForwardedFor 返回 X-Forwarded-For 链中的第一个 IP（去除空白）。
// 空字符串表示该头不存在或为空。
func firstForwardedFor(r *http.Request) string {
	ff := r.Header.Get("X-Forwarded-For")
	if ff == "" {
		return ""
	}
	if i := strings.IndexByte(ff, ','); i >= 0 {
		return strings.TrimSpace(ff[:i])
	}
	return strings.TrimSpace(ff)
}

// stripPort removes the port portion from an address in "host:port" format.
// IPv6 addresses are handled correctly.
//
// stripPort 去除 "host:port" 格式地址中的端口部分，正确处理 IPv6 地址。
func stripPort(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}
