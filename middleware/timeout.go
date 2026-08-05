package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/rambollwong/rainbowferret/types"
	"github.com/rambollwong/rainbowferret/util"
)

// Timeout returns middleware that cancels the request context after the given
// duration. When the context expires before the handler writes a response a
// 504 Gateway Timeout JSON response is sent.
//
// The handler runs in the same goroutine as the caller; panics are therefore
// still visible to Recoverer. Handlers that perform blocking I/O should
// respect r.Context().Done() so the timeout takes effect.
//
// Timeout 返回在指定时长后取消请求 context 的中间件。
// 若 context 到期时 handler 尚未写入响应，则返回 504 Gateway Timeout JSON。
//
// handler 与调用者在同一 goroutine 中运行，因此 Recoverer 仍可捕获 panic。
// 执行阻塞 I/O 的 handler 应监听 r.Context().Done() 以使超时生效。
func Timeout(d time.Duration) func(next http.Handler) http.Handler {
	return TimeoutWithConfig(TimeoutConfig{Timeout: d})
}

// TimeoutConfig holds the configuration for the Timeout middleware.
// TimeoutConfig 存储 Timeout 中间件的配置。
type TimeoutConfig struct {
	// Timeout is the maximum duration the handler is allowed to run.
	// Defaults to 30 seconds when <= 0.
	// Timeout 是 handler 允许运行的最大时长，<= 0 时默认为 30 秒。
	Timeout time.Duration

	// OnTimeout is called when the timeout fires. When nil a default 504
	// JSON response is written.
	// OnTimeout 在超时触发时调用。为 nil 时写入默认的 504 JSON 响应。
	OnTimeout http.HandlerFunc
}

// TimeoutWithConfig returns a Timeout middleware configured with cfg.
// TimeoutWithConfig 返回使用 cfg 配置的 Timeout 中间件。
func TimeoutWithConfig(cfg TimeoutConfig) func(next http.Handler) http.Handler {
	d := cfg.Timeout
	if d <= 0 {
		d = 30 * time.Second
	}
	onTimeout := cfg.OnTimeout
	if onTimeout == nil {
		onTimeout = defaultTimeoutHandler
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()

			rw := &headerWriter{ResponseWriter: w}
			next.ServeHTTP(rw, r.WithContext(ctx))

			// If the context expired and the handler hasn't started
			// a response yet, send the timeout response.
			// 若 context 到期且 handler 尚未开始写入响应，
			// 则发送超时响应。
			if ctx.Err() != nil && !rw.wrote {
				onTimeout.ServeHTTP(w, r)
			}
		})
	}
}

// headerWriter records whether WriteHeader has been called.
// headerWriter 记录 WriteHeader 是否已被调用。
type headerWriter struct {
	http.ResponseWriter
	wrote bool
}

func (hw *headerWriter) WriteHeader(code int) {
	hw.wrote = true
	hw.ResponseWriter.WriteHeader(code)
}

func (hw *headerWriter) Write(b []byte) (int, error) {
	if !hw.wrote {
		hw.wrote = true
		hw.ResponseWriter.WriteHeader(http.StatusOK)
	}
	return hw.ResponseWriter.Write(b)
}

// Unwrap returns the underlying ResponseWriter.
// Unwrap 返回底层 ResponseWriter。
func (hw *headerWriter) Unwrap() http.ResponseWriter {
	return hw.ResponseWriter
}

var defaultTimeoutHandler http.HandlerFunc = func(w http.ResponseWriter, _ *http.Request) {
	util.WriteJSON(w, http.StatusGatewayTimeout, &types.HTTPError{
		Code:    http.StatusGatewayTimeout,
		Message: http.StatusText(http.StatusGatewayTimeout),
	})
}
