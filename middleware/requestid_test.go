package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDGenerate(t *testing.T) {
	mw := RequestID()
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetRequestID(r.Context())
		if id == "" {
			t.Fatal("expected non-empty request id")
		}
		if len(id) < 32 {
			t.Fatalf("request id too short: %d", len(id))
		}
		w.Header().Set("X-Debug-ID", id)
		w.WriteHeader(200)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected X-Request-ID response header")
	}
}

func TestRequestIDReuse(t *testing.T) {
	mw := RequestID()
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetRequestID(r.Context())
		w.Write([]byte(id))
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Request-ID", "upstream-trace-42")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Body.String() != "upstream-trace-42" {
		t.Fatalf("expected reused id, got %q", w.Body.String())
	}
}

func TestRequestIDCustomHeader(t *testing.T) {
	mw := RequestIDWithConfig(RequestIDConfig{HeaderName: "X-Correlation-ID"})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Header().Get("X-Correlation-ID") == "" {
		t.Fatal("expected X-Correlation-ID header")
	}
	if w.Header().Get("X-Request-ID") != "" {
		t.Fatal("X-Request-ID should not be set with custom header name")
	}
}
