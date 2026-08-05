package middleware

import (
	"net/http"
	"strings"

	"github.com/rambollwong/rainbowferret/types"
	"github.com/rambollwong/rainbowferret/util"
)

// ContentEncoding returns middleware that rejects requests whose
// Content-Encoding header is not in the allowed list with a 415 Unsupported
// Media Type response. An empty or missing Content-Encoding header is always
// accepted. Matching is case-insensitive. Methods with no body (GET, HEAD,
// OPTIONS, TRACE) are always skipped.
//
// ContentEncoding 返回一个中间件，拒绝 Content-Encoding 不在允许列表中的
// 请求，返回 415 Unsupported Media Type。空或缺失的 Content-Encoding 始终放行。
// 匹配忽略大小写。无请求体的方法始终跳过。
func ContentEncoding(allowed ...string) func(next http.Handler) http.Handler {
	allow := make(map[string]bool, len(allowed))
	for _, e := range allowed {
		allow[strings.ToLower(strings.TrimSpace(e))] = true
	}
	// "identity" means no encoding; always allowed when the header is empty.
	allow[""] = true

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if noBodyMethods[r.Method] {
				next.ServeHTTP(w, r)
				return
			}

			enc := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Encoding")))
			if !allow[enc] {
				util.WriteJSON(w, http.StatusUnsupportedMediaType, &types.HTTPError{
					Code:    http.StatusUnsupportedMediaType,
					Message: "unsupported content encoding: " + enc,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
