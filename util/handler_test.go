package util

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rambollwong/rainbowferret/types"
)

func TestHandleTPointerRequest(t *testing.T) {
	type userReq struct {
		Name string `json:"name"`
	}

	h := HandleT(
		types.HandlerFunc[*userReq, map[string]string](
			func(ctx context.Context, req *userReq) (map[string]string, error) {
				if req == nil {
					t.Error("req should not be nil")
					return nil, nil
				}
				return map[string]string{"name": req.Name}, nil
			},
		),
	)

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"alice"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%q", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "alice") {
		t.Fatalf("expected body to contain alice, got %q", w.Body.String())
	}
}

func TestHandleTValueRequest(t *testing.T) {
	type userReq struct {
		Name string `json:"name"`
	}

	h := HandleT(
		func(ctx context.Context, req userReq) (map[string]string, error) {
			return map[string]string{"name": req.Name}, nil
		},
	)

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"bob"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%q", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "bob") {
		t.Fatalf("expected body to contain bob, got %q", w.Body.String())
	}
}
