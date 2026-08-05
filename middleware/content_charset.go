package middleware

import (
	"mime"
	"net/http"

	"github.com/rambollwong/rainbowferret/types"
	"github.com/rambollwong/rainbowferret/util"
)

// ContentCharset returns middleware that rejects requests whose Content-Type
// charset parameter is not in the allowed list with a 415 Unsupported Media
// Type response. Requests without an explicit charset are accepted (most
// clients do not send it). Matching is case-insensitive. Methods with no body
// (GET, HEAD, OPTIONS, TRACE) are always skipped.
//
// ContentCharset 返回一个中间件，拒绝 Content-Type 中 charset 参数不在
// 允许列表中的请求，返回 415 Unsupported Media Type。未指定 charset 的
// 请求会被放行。匹配忽略大小写。无请求体的方法始终跳过。
func ContentCharset(allowed ...string) func(next http.Handler) http.Handler {
	allow := make(map[string]bool, len(allowed))
	for _, c := range allowed {
		allow[http.CanonicalHeaderKey(c)] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if noBodyMethods[r.Method] {
				next.ServeHTTP(w, r)
				return
			}

			ct := r.Header.Get("Content-Type")
			_, params, err := mime.ParseMediaType(ct)
			if err != nil {
				// Malformed Content-Type; let ContentType middleware
				// (if any) handle it. Accept it here.
				next.ServeHTTP(w, r)
				return
			}

			charset := params["charset"]
			if charset == "" {
				// No charset parameter — accept.
				next.ServeHTTP(w, r)
				return
			}

			if !allow[http.CanonicalHeaderKey(charset)] {
				util.WriteJSON(w, http.StatusUnsupportedMediaType, &types.HTTPError{
					Code:    http.StatusUnsupportedMediaType,
					Message: "unsupported charset: " + charset,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
