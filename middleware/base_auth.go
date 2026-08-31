package middleware

import (
	"net/http"
	"strings"

	"github.com/rambollwong/rainbowferret/types"
	"github.com/rambollwong/rainbowferret/util"
)

// BaseAuthConfig holds the configuration for the BaseAuth middleware.
// BaseAuthConfig 存储 BaseAuth 中间件的配置。
type BaseAuthConfig struct {
	// AuthFunc decides whether a credential is valid. The scheme argument
	// identifies the credential source: the matched Authorization scheme
	// (e.g. "Bearer") or the API key header name (e.g. "X-API-Key").
	//
	// AuthFunc 决定凭证是否有效。scheme 参数标识凭证来源：匹配到的
	// Authorization 方案（如 "Bearer"）或 API key 头名（如 "X-API-Key"）。
	AuthFunc func(scheme, key string) bool

	// Schemes is the list of Authorization header prefixes accepted, matched
	// case-insensitively. Defaults to ["Bearer"].
	//
	// Schemes 是 Authorization 头可接受的认证方案（前缀）列表，大小写不敏感
	// 匹配。默认为 ["Bearer"]。
	Schemes []string

	// APIKeyHeader is the header name used as the API key fallback. When empty
	// the API key fallback is disabled and only Authorization is inspected.
	//
	// APIKeyHeader 是作为 API key 回退的请求头名。为空时禁用 API key 回退，
	// 仅检查 Authorization 头。
	APIKeyHeader string
}

// BaseAuth returns middleware that authenticates each request. Credentials are
// read from the Authorization header (matching one of cfg.Schemes) or, as a
// fallback, from cfg.APIKeyHeader. Missing or invalid credentials receive
// 401 Unauthorized with a WWW-Authenticate challenge header.
//
// BaseAuth 返回一个为每个请求鉴权的中间件。凭证从 Authorization 头（匹配
// cfg.Schemes 之一）读取，或回退到 cfg.APIKeyHeader。缺失或无效的凭证返回
// 401 Unauthorized，并附带 WWW-Authenticate challenge 头。
func BaseAuth(cfg BaseAuthConfig) func(next http.Handler) http.Handler {
	schemes := cfg.Schemes
	if len(schemes) == 0 {
		schemes = []string{"Bearer"}
	}
	apiKeyHeader := cfg.APIKeyHeader
	if apiKeyHeader == "" {
		apiKeyHeader = "X-API-Key"
	}
	challenge := schemes[0]

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scheme, key, ok := extractCredential(r, schemes, apiKeyHeader)
			if !ok || cfg.AuthFunc == nil || !cfg.AuthFunc(scheme, key) {
				w.Header().Set("WWW-Authenticate", challenge)
				util.WriteJSON(w, http.StatusUnauthorized, types.Unauthorized("unauthorized"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// extractCredential reads the credential from the request. It first tries each
// Authorization scheme (case-insensitive) and then falls back to the API key
// header. ok is false when no credential is present.
//
// extractCredential 从请求中读取凭证。先依次尝试每个 Authorization 方案
// （大小写不敏感），然后回退到 API key 头。无凭证时 ok 为 false。
func extractCredential(r *http.Request, schemes []string, apiKeyHeader string) (scheme, key string, ok bool) {
	if auth := r.Header.Get("Authorization"); auth != "" {
		fields := strings.Fields(auth)
		if len(fields) == 2 {
			for _, s := range schemes {
				if strings.EqualFold(fields[0], s) {
					return s, fields[1], true
				}
			}
		}
	}
	if apiKeyHeader != "" {
		if key := strings.TrimSpace(r.Header.Get(apiKeyHeader)); key != "" {
			return apiKeyHeader, key, true
		}
	}
	return "", "", false
}
