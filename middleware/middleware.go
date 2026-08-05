package middleware

import (
	"os"
	"strings"
	"sync"
	"unicode"
)

// ctxKey is a unique type for context keys to avoid collisions between packages.
// ctxKey 是 context 键的唯一类型，用于避免包之间的键冲突。
type ctxKey string

var (
	hostname     string
	hostnameOnce sync.Once
)

// cachedHostname returns the sanitized hostname, cached after the first call.
// cachedHostname 返回清理后的主机名，首次调用后缓存。
func cachedHostname() string {
	hostnameOnce.Do(func() {
		h, err := os.Hostname()
		if err != nil {
			return
		}
		// Keep only the first label and strip non-printable characters.
		// 仅保留第一个标签并剔除非打印字符。
		h, _, _ = strings.Cut(h, ".")
		h = strings.TrimFunc(h, func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-'
		})
		if len(h) > 32 {
			h = h[:32]
		}
		hostname = h
	})
	return hostname
}
