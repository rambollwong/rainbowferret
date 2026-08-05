package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContentEncoding(t *testing.T) {
	mw := ContentEncoding("gzip", "br")
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	tests := []struct {
		name string
		enc  string
		code int
	}{
		{"gzip", "gzip", 200},
		{"br", "br", 200},
		{"no encoding", "", 200},
		{"deflate rejected", "deflate", 415},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/", strings.NewReader("{}"))
			if tc.enc != "" {
				r.Header.Set("Content-Encoding", tc.enc)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.code {
				t.Fatalf("expected %d, got %d", tc.code, w.Code)
			}
		})
	}
}
