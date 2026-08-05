package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture the status code and the
// number of bytes written to the response body.
//
// responseWriter 包装 http.ResponseWriter，用于捕获状态码和写入响应体的字节数。
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// WriteHeader captures the status code and delegates to the underlying writer.
// WriteHeader 捕获状态码并委托给底层 writer。
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the number of bytes written. If WriteHeader has not been
// called explicitly it defaults to 200 OK.
//
// Write 捕获写入的字节数。若 WriteHeader 未被显式调用则默认为 200 OK。
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// Logger returns middleware that logs every request using the standard log
// package. The log line includes method, path, status code, duration, and
// response body size.
//
// Format: METHOD /path 200 1.234ms 42B
//
// Logger 返回一个使用标准 log 包记录每个请求的中间件。
// 日志行包含方法、路径、状态码、耗时和响应体大小。
//
// 格式: METHOD /path 200 1.234ms 42B
func Logger() func(next http.Handler) http.Handler {
	return LoggerWithConfig(LoggerConfig{})
}

// LoggerConfig holds the configuration for the Logger middleware.
// LoggerConfig 存储 Logger 中间件的配置。
type LoggerConfig struct {
	// Logf is the function used to emit log lines. When nil, log.Printf is used.
	// Logf 是用于输出日志行的函数。为 nil 时使用 log.Printf。
	Logf func(format string, args ...any)
}

// LoggerWithConfig returns a Logger middleware configured with cfg.
// LoggerWithConfig 返回使用 cfg 配置的 Logger 中间件。
func LoggerWithConfig(cfg LoggerConfig) func(next http.Handler) http.Handler {
	logf := cfg.Logf
	if logf == nil {
		logf = log.Printf
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w}
			next.ServeHTTP(rw, r)
			logf("%s %s %d %v %dB",
				r.Method,
				r.URL.RequestURI(),
				rw.status,
				time.Since(start).Truncate(time.Microsecond),
				rw.size,
			)
		})
	}
}
