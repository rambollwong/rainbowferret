package util

import (
	"net/http/httptest"
	"testing"
)

func TestQueryParam(t *testing.T) {
	r := httptest.NewRequest("GET", "/?name=alice&age=30&active=true&pi=3.14&big=9223372036854775807", nil)

	if v := QueryParam(r, "name"); v != "alice" {
		t.Fatalf("expected alice, got %q", v)
	}
	if v := QueryParam(r, "missing"); v != "" {
		t.Fatalf("expected empty, got %q", v)
	}
}

func TestQueryParamDefault(t *testing.T) {
	r := httptest.NewRequest("GET", "/?name=alice", nil)

	if v := QueryParamDefault(r, "name", "bob"); v != "alice" {
		t.Fatalf("expected alice, got %q", v)
	}
	if v := QueryParamDefault(r, "missing", "default"); v != "default" {
		t.Fatalf("expected default, got %q", v)
	}
}

func TestQueryParamInt(t *testing.T) {
	r := httptest.NewRequest("GET", "/?age=30", nil)

	v, err := QueryParamInt(r, "age")
	if err != nil || v != 30 {
		t.Fatalf("expected 30, got %d err=%v", v, err)
	}
	_, err = QueryParamInt(r, "missing")
	if err == nil {
		t.Fatal("expected error for missing param")
	}
}

func TestQueryParamIntDefault(t *testing.T) {
	r := httptest.NewRequest("GET", "/?age=bad", nil)

	if v := QueryParamIntDefault(r, "age", 99); v != 99 {
		t.Fatalf("expected default 99, got %d", v)
	}
}

func TestQueryParamInt64(t *testing.T) {
	r := httptest.NewRequest("GET", "/?big=9223372036854775807", nil)

	v, err := QueryParamInt64(r, "big")
	if err != nil || v != 9223372036854775807 {
		t.Fatalf("expected max int64, got %d err=%v", v, err)
	}
}

func TestQueryParamFloat(t *testing.T) {
	r := httptest.NewRequest("GET", "/?pi=3.14", nil)

	v, err := QueryParamFloat(r, "pi")
	if err != nil || v != 3.14 {
		t.Fatalf("expected 3.14, got %f err=%v", v, err)
	}
}

func TestQueryParamBool(t *testing.T) {
	r := httptest.NewRequest("GET", "/?active=true&flag=0", nil)

	v, err := QueryParamBool(r, "active")
	if err != nil || !v {
		t.Fatalf("expected true, got %v err=%v", v, err)
	}
	v, err = QueryParamBool(r, "flag")
	if err != nil || v {
		t.Fatalf("expected false, got %v err=%v", v, err)
	}
}

func TestPathParam(t *testing.T) {
	// Go 1.22+ PathValue
	r := httptest.NewRequest("GET", "/users/42", nil)
	r.SetPathValue("id", "42")

	if v := PathParam(r, "id"); v != "42" {
		t.Fatalf("expected 42, got %q", v)
	}
	if v := PathParam(r, "missing"); v != "" {
		t.Fatalf("expected empty, got %q", v)
	}
}

func TestPathParamInt(t *testing.T) {
	r := httptest.NewRequest("GET", "/users/42", nil)
	r.SetPathValue("id", "42")

	v, err := PathParamInt(r, "id")
	if err != nil || v != 42 {
		t.Fatalf("expected 42, got %d err=%v", v, err)
	}
}

func TestPathParamIntDefault(t *testing.T) {
	r := httptest.NewRequest("GET", "/users/bad", nil)
	r.SetPathValue("id", "bad")

	if v := PathParamIntDefault(r, "id", 0); v != 0 {
		t.Fatalf("expected default 0, got %d", v)
	}
}

func TestPathParamFloat(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.SetPathValue("pi", "3.14")

	v, err := PathParamFloat(r, "pi")
	if err != nil || v != 3.14 {
		t.Fatalf("expected 3.14, got %f", v)
	}
}

func TestPathParamBool(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.SetPathValue("flag", "true")

	v, err := PathParamBool(r, "flag")
	if err != nil || !v {
		t.Fatalf("expected true, got %v", v)
	}
}
