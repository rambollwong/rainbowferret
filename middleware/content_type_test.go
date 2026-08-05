package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContentTypeAllowed(t *testing.T) {
	mw := ContentType("application/json", "application/xml")
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	tests := []struct {
		name string
		ct   string
		code int
	}{
		{"json", "application/json", 200},
		{"json+charset", "application/json; charset=utf-8", 200},
		{"xml", "application/xml", 200},
		{"plain rejected", "text/plain", 415},
		{"empty rejected", "", 415},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body := strings.NewReader("{}")
			r := httptest.NewRequest("POST", "/", body)
			if tc.ct != "" {
				r.Header.Set("Content-Type", tc.ct)
			}
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.code {
				t.Fatalf("expected %d, got %d", tc.code, w.Code)
			}
		})
	}
}

func TestContentTypeSkipNoBody(t *testing.T) {
	mw := ContentType("application/json")
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	for _, method := range []string{"GET", "HEAD", "OPTIONS", "TRACE"} {
		t.Run(method, func(t *testing.T) {
			r := httptest.NewRequest(method, "/", nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != 200 {
				t.Fatalf("%s: expected 200, got %d", method, w.Code)
			}
		})
	}
}
