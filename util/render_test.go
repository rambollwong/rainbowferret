package util

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	WriteJSON(w, 201, map[string]string{"hello": "world"})

	if w.Code != 201 {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Fatalf("wrong content-type: %q", ct)
	}
	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["hello"] != "world" {
		t.Fatalf("wrong body: %v", body)
	}
}

func TestWriteNoContent(t *testing.T) {
	w := httptest.NewRecorder()
	WriteNoContent(w)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Fatalf("expected empty body")
	}
}

func TestWriteText(t *testing.T) {
	w := httptest.NewRecorder()
	WriteText(w, 200, "hello")

	if w.Body.String() != "hello" {
		t.Fatalf("expected hello, got %q", w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Fatalf("wrong content-type: %q", ct)
	}
}

func TestWriteXML(t *testing.T) {
	type Item struct {
		XMLName struct{} `xml:"item"`
		Value   string   `xml:"value"`
	}
	w := httptest.NewRecorder()
	WriteXML(w, 200, Item{Value: "hello"})

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	ct := w.Header().Get("Content-Type")
	if ct != "application/xml; charset=utf-8" {
		t.Fatalf("wrong content-type: %q", ct)
	}
	expected := xml.Header + "<item><value>hello</value></item>"
	if w.Body.String() != expected {
		t.Fatalf("wrong body: got %q, want %q", w.Body.String(), expected)
	}
}

func TestWriteStream(t *testing.T) {
	w := httptest.NewRecorder()
	WriteStream(w, 200, "text/csv", strings.NewReader("a,b,c"))

	if w.Body.String() != "a,b,c" {
		t.Fatalf("wrong body: %q", w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/csv" {
		t.Fatalf("wrong content-type: %q", ct)
	}
}
