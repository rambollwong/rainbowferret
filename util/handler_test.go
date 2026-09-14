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

// createdResp implements types.SuccessStatuser, returning 201 Created.
// createdResp 实现 types.SuccessStatuser，返回 201 Created。
type createdResp struct {
	ID string `json:"id"`
}

func (createdResp) SuccessStatus() int { return http.StatusCreated }

func TestHandleTSuccessStatus(t *testing.T) {
	type createReq struct {
		Name string `json:"name"`
	}
	h := HandleT[createReq, createdResp](
		func(ctx context.Context, req createReq) (createdResp, error) {
			return createdResp{ID: "1"}, nil
		},
	)

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"alice"}`))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%q", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"id":"1"`) {
		t.Fatalf("expected body to contain id, got %q", w.Body.String())
	}
}
