package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBaseAuthBearer(t *testing.T) {
	mw := BaseAuth(BaseAuthConfig{
		AuthFunc: func(scheme, key string) bool {
			return scheme == "Bearer" && key == "valid-token"
		},
	})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	tests := []struct {
		name     string
		header   string
		value    string
		wantCode int
	}{
		{"valid bearer", "Authorization", "Bearer valid-token", 200},
		{"invalid bearer", "Authorization", "Bearer wrong", 401},
		{"missing auth", "", "", 401},
		{"lowercase bearer", "Authorization", "bearer valid-token", 200},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			if tc.header != "" {
				r.Header.Set(tc.header, tc.value)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.wantCode {
				t.Fatalf("expected %d, got %d", tc.wantCode, w.Code)
			}
		})
	}
}

func TestBaseAuthWWWAuthenticate(t *testing.T) {
	mw := BaseAuth(BaseAuthConfig{
		AuthFunc: func(scheme, key string) bool { return false },
	})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if got := w.Header().Get("WWW-Authenticate"); got != "Bearer" {
		t.Fatalf("expected WWW-Authenticate Bearer, got %q", got)
	}
}

func TestBaseAuthMultipleSchemes(t *testing.T) {
	mw := BaseAuth(BaseAuthConfig{
		Schemes: []string{"Bearer", "Token"},
		AuthFunc: func(scheme, key string) bool {
			return (scheme == "Bearer" || scheme == "Token") && key == "secret"
		},
	})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	for _, scheme := range []string{"Bearer", "Token"} {
		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Authorization", scheme+" secret")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 200 {
			t.Fatalf("scheme %s: expected 200, got %d", scheme, w.Code)
		}
	}
}

func TestBaseAuthAPIKey(t *testing.T) {
	mw := BaseAuth(BaseAuthConfig{
		AuthFunc: func(scheme, key string) bool {
			return scheme == "X-API-Key" && key == "secret-key"
		},
	})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-API-Key", "secret-key")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// Bearer 优先于 X-API-Key
	r2 := httptest.NewRequest("GET", "/", nil)
	r2.Header.Set("Authorization", "Bearer wrong")
	r2.Header.Set("X-API-Key", "secret-key")
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	if w2.Code != 401 {
		t.Fatalf("expected 401 (Bearer precedence), got %d", w2.Code)
	}
}

func TestBaseAuthCustomAPIKeyHeader(t *testing.T) {
	mw := BaseAuth(BaseAuthConfig{
		APIKeyHeader: "X-Custom-Key",
		AuthFunc: func(scheme, key string) bool {
			return scheme == "X-Custom-Key" && key == "custom"
		},
	})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Custom-Key", "custom")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
