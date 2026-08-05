package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientIPNoProxy(t *testing.T) {
	mw := ClientIP()
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(GetClientIP(r.Context())))
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "192.168.1.100:54321"
	r.Header.Set("X-Forwarded-For", "1.2.3.4")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Body.String() != "192.168.1.100" {
		t.Fatalf("expected 192.168.1.100, got %q", w.Body.String())
	}
}

func TestClientIPTrustProxy(t *testing.T) {
	mw := ClientIPWithConfig(ClientIPConfig{TrustProxy: true})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(GetClientIP(r.Context())))
	}))

	tests := []struct {
		name     string
		xff      string
		xri      string
		remote   string
		expected string
	}{
		{"xff first", "1.2.3.4, 5.6.7.8", "", "10.0.0.1:12345", "1.2.3.4"},
		{"x-real-ip", "", "8.8.8.8", "10.0.0.1:12345", "8.8.8.8"},
		{"fallback remote", "", "", "10.0.0.1:12345", "10.0.0.1"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.RemoteAddr = tc.remote
			if tc.xff != "" {
				r.Header.Set("X-Forwarded-For", tc.xff)
			}
			if tc.xri != "" {
				r.Header.Set("X-Real-IP", tc.xri)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Body.String() != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, w.Body.String())
			}
		})
	}
}
