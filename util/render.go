package util

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
)

// WriteJSON writes a JSON-encoded response with the given status code.
// WriteJSON 以给定的状态码写入 JSON 编码的响应。
func WriteJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json error: %v", err)
	}
}

// WriteNoContent writes a 204 No Content response with an empty body.
// WriteNoContent 写入一个无内容的 204 响应。
func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// WriteText writes a plain text response with the given status code.
// WriteText 以给定的状态码写入纯文本响应。
func WriteText(w http.ResponseWriter, code int, text string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	if _, err := io.WriteString(w, text); err != nil {
		log.Printf("write text error: %v", err)
	}
}

// WriteXML writes an XML-encoded response with the given status code.
// The XML declaration header is automatically prepended to the output.
//
// WriteXML 以给定的状态码写入 XML 编码的响应，
// 输出中会自动添加 XML 声明头。
func WriteXML(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(code)
	if _, err := io.WriteString(w, xml.Header); err != nil {
		log.Printf("write xml header error: %v", err)
		return
	}
	if err := xml.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write xml error: %v", err)
	}
}

// WriteStream reads from r and writes it to the response body with the given
// status code and Content-Type. It is suitable for streaming large payloads,
// file downloads, or Server-Sent Events (with "text/event-stream").
//
// WriteStream 从 r 中读取数据并写入响应体，使用指定的状态码和 Content-Type，
// 适用于流式传输大文件、文件下载或 Server-Sent Events（配合 "text/event-stream"）。
func WriteStream(w http.ResponseWriter, code int, contentType string, r io.Reader) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	if _, err := io.Copy(w, r); err != nil {
		log.Printf("write stream error: %v", err)
	}
}
