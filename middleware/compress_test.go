package middleware

import (
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompressGzip(t *testing.T) {
	mw := Compress()
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello compress"))
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	ce := w.Header().Get("Content-Encoding")
	if ce != "gzip" {
		t.Fatalf("expected gzip, got %q", ce)
	}

	zr, _ := gzip.NewReader(w.Body)
	raw, _ := io.ReadAll(zr)
	zr.Close()
	if string(raw) != "hello compress" {
		t.Fatalf("wrong body: %q", string(raw))
	}
}

func TestCompressDeflate(t *testing.T) {
	mw := Compress()
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello deflate"))
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Encoding", "deflate")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	ce := w.Header().Get("Content-Encoding")
	if ce != "deflate" {
		t.Fatalf("expected deflate, got %q", ce)
	}

	zr, _ := zlib.NewReader(w.Body)
	raw, _ := io.ReadAll(zr)
	zr.Close()
	if string(raw) != "hello deflate" {
		t.Fatalf("wrong body: %q", string(raw))
	}
}

func TestCompressNoEncoding(t *testing.T) {
	mw := Compress()
	h := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("plain"))
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)

	if w.Header().Get("Content-Encoding") != "" {
		t.Fatal("expected no Content-Encoding header")
	}
}
