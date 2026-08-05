package ferret

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGroupRouting(t *testing.T) {
	sm := http.NewServeMux()
	g := NewRootGroup(sm)

	g.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("get"))
	})
	g.Post("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("post"))
	})

	h := g.Handler()

	tests := []struct {
		method string
		path   string
		code   int
		body   string
	}{
		{"GET", "/hello", 200, "get"},
		{"POST", "/hello", 200, "post"},
		{"PUT", "/hello", 405, ""},
		{"GET", "/unknown", 404, ""},
	}
	for _, tc := range tests {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != tc.code {
				t.Fatalf("expected %d, got %d", tc.code, w.Code)
			}
			if tc.body != "" && w.Body.String() != tc.body {
				t.Fatalf("expected body %q, got %q", tc.body, w.Body.String())
			}
		})
	}
}

func TestSubGroup(t *testing.T) {
	sm := http.NewServeMux()
	root := NewRootGroup(sm)
	api := root.Group("/api")

	api.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("users"))
	})
	root.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	h := root.Handler()

	r1 := httptest.NewRequest("GET", "/api/users", nil)
	w1 := httptest.NewRecorder()
	h.ServeHTTP(w1, r1)
	if w1.Code != 200 || w1.Body.String() != "users" {
		t.Fatalf("sub-group route failed: %d %s", w1.Code, w1.Body.String())
	}

	r2 := httptest.NewRequest("GET", "/health", nil)
	w2 := httptest.NewRecorder()
	h.ServeHTTP(w2, r2)
	if w2.Code != 200 || w2.Body.String() != "ok" {
		t.Fatalf("root route failed: %d %s", w2.Code, w2.Body.String())
	}
}

func TestMethodNotAllowed(t *testing.T) {
	sm := http.NewServeMux()
	g := NewRootGroup(sm)
	g.Get("/only-get", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	h := g.Handler()

	r := httptest.NewRequest("POST", "/only-get", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 405 {
		t.Fatalf("expected 405, got %d", w.Code)
	}
	allow := w.Header().Get("Allow")
	if allow != "GET" {
		t.Fatalf("expected Allow: GET, got %q", allow)
	}
}

func TestCustomNotFound(t *testing.T) {
	sm := http.NewServeMux()
	g := NewRootGroup(sm)
	g.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("custom not found"))
	})

	h := g.Handler()
	r := httptest.NewRequest("GET", "/nope", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Body.String() != "custom not found" {
		t.Fatalf("expected custom not found, got %q", w.Body.String())
	}
}

func TestCustomMethodNotAllowed(t *testing.T) {
	sm := http.NewServeMux()
	g := NewRootGroup(sm)
	g.Get("/only-get", func(w http.ResponseWriter, r *http.Request) {})
	g.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("nope"))
	})

	h := g.Handler()
	r := httptest.NewRequest("POST", "/only-get", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Body.String() != "nope" {
		t.Fatalf("expected custom 405 body, got %q", w.Body.String())
	}
}

func TestSubGroupMiddleware(t *testing.T) {
	sm := http.NewServeMux()
	root := NewRootGroup(sm)

	marker := ""
	markMW := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			marker += "A"
			next.ServeHTTP(w, r)
		})
	}

	api := root.Group("/api", markMW)
	api.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		marker += "B"
	})

	h := root.Handler()
	r := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if marker != "AB" {
		t.Fatalf("expected middleware applied, marker=%q", marker)
	}
}
