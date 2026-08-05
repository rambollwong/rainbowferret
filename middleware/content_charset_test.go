package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContentCharset(t *testing.T) {
	mw := ContentCharset("utf-8")
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	tests := []struct {
		name string
		ct   string
		code int
	}{
		{"utf-8", "application/json; charset=utf-8", 200},
		{"UTF-8 case", "application/json; charset=UTF-8", 200},
		{"no charset", "application/json", 200},
		{"empty ct", "", 200},
		{"iso rejected", "text/plain; charset=iso-8859-1", 415},
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
