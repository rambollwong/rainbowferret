package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	var buf strings.Builder
	mw := LoggerWithConfig(LoggerConfig{
		Logf: func(format string, args ...any) {
			buf.WriteString(fmt.Sprintf(format, args...))
			buf.WriteString("\n")
		},
	})
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		w.Write([]byte("ok"))
	}))

	r := httptest.NewRequest("POST", "/api/v1/test", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 201 {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	if !strings.Contains(buf.String(), "POST /api/v1/test 201") {
		t.Fatalf("unexpected log line: %s", buf.String())
	}
}
