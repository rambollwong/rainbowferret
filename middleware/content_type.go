package middleware

import (
	"mime"
	"net/http"
	"strings"

	"github.com/rambollwong/rainbowferret/types"
	"github.com/rambollwong/rainbowferret/util"
)

// noBodyMethods is the set of HTTP methods that carry no request body and are
// therefore exempt from Content-Type checks.
//
// noBodyMethods 是不带请求体的 HTTP 方法集合，因此豁免 Content-Type 检查。
var noBodyMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodHead:    true,
	http.MethodOptions: true,
	http.MethodTrace:   true,
}

// ContentType returns middleware that rejects requests whose Content-Type is
// not in the allowed list with a 415 Unsupported Media Type response. Matching
// is case-insensitive and ignores parameters such as "; charset=utf-8".
// Methods with no body (GET, HEAD, OPTIONS, TRACE) are always skipped.
//
// ContentType 返回一个中间件，拒绝 Content-Type 不在允许列表中的请求，
// 返回 415 Unsupported Media Type。匹配忽略大小写和参数。
// 无请求体的方法始终跳过。
func ContentType(allowed ...string) func(next http.Handler) http.Handler {
	allow := make(map[string]bool, len(allowed))
	for _, t := range allowed {
		allow[strings.ToLower(strings.TrimSpace(t))] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if noBodyMethods[r.Method] {
				next.ServeHTTP(w, r)
				return
			}

			ct := r.Header.Get("Content-Type")
			mediaType, _, err := mime.ParseMediaType(ct)
			if err != nil {
				mediaType = strings.TrimSpace(strings.SplitN(ct, ";", 2)[0])
			}

			if !allow[strings.ToLower(mediaType)] {
				util.WriteJSON(w, http.StatusUnsupportedMediaType, &types.HTTPError{
					Code:    http.StatusUnsupportedMediaType,
					Message: "unsupported media type: " + mediaType,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
