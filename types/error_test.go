package types

import (
	"net/http"
	"testing"
)

func TestHTTPError(t *testing.T) {
	e := NewHTTPError(400, "bad request")
	if e.Error() != "bad request" {
		t.Fatalf("expected 'bad request', got %q", e.Error())
	}
	if e.Code != 400 {
		t.Fatalf("expected 400, got %d", e.Code)
	}
}

func TestBadRequest(t *testing.T) {
	e := BadRequest("missing field")
	if e.Code != 400 {
		t.Fatalf("expected 400, got %d", e.Code)
	}
}

func TestNotFound(t *testing.T) {
	e := NotFound("resource not found")
	if e.Code != 404 {
		t.Fatalf("expected 404, got %d", e.Code)
	}
}

func TestUnauthorized(t *testing.T) {
	e := Unauthorized("invalid token")
	if e.Code != 401 {
		t.Fatalf("expected 401, got %d", e.Code)
	}
}

func TestUnsupportedMediaType(t *testing.T) {
	e := UnsupportedMediaType("text/plain")
	if e.Code != 415 {
		t.Fatalf("expected 415, got %d", e.Code)
	}
}

func TestInternal(t *testing.T) {
	e := Internal("something broke")
	if e.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", e.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	e := MethodNotAllowed("POST not allowed")
	if e.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", e.Code)
	}
}
