package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/rambollwong/rainbowferret/types"
	"github.com/rambollwong/rainbowferret/util"
)

// Recoverer is middleware that recovers from panics, logs the stack trace,
// and responds with a 500 Internal Server Error. If the panic value is an
// *HTTPError it honours the embedded status code.
//
// Recoverer 是一个从 panic 中恢复的中间件，记录堆栈并响应 500 Internal Server
// Error。若 panic 值为 *HTTPError 则使用其嵌入的状态码。
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if httpErr, ok := rec.(*types.HTTPError); ok {
					util.WriteJSON(w, httpErr.Code, httpErr)
				} else {
					log.Printf("panic: %v\n%s", rec, debug.Stack())
					util.WriteJSON(w, 500, &types.HTTPError{Code: 500, Message: "internal server error"})
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
