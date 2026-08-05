package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTimeoutNormal(t *testing.T) {
	mw := Timeout(1 * time.Second)
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestTimeoutFires(t *testing.T) {
	mw := Timeout(10 * time.Millisecond)
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate work that respects context cancellation.
		select {
		case <-time.After(100 * time.Millisecond):
			w.Write([]byte("too late"))
		case <-r.Context().Done():
			return
		}
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 504 {
		t.Fatalf("expected 504, got %d body=%q", w.Code, w.Body.String())
	}
}

func TestTimeoutHandlerWritesBeforeTimeout(t *testing.T) {
	// Handler writes response before context expires → 200
	mw := Timeout(1 * time.Second)
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		// context expires after, but response already started
		select {
		case <-time.After(100 * time.Millisecond):
		case <-r.Context().Done():
		}
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("expected 200 (already wrote header), got %d", w.Code)
	}
}
