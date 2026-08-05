package middleware

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"strings"
	"sync"
)

// Compress returns middleware that compresses the response body using one of
// the supported encodings. The encoding is negotiated from the request's
// Accept-Encoding header; when the client does not accept any supported
// encoding the response is passed through uncompressed.
//
// An optional compression level may be passed (0-9, or the gzip package
// constants NoCompression / BestSpeed / DefaultCompression / BestCompression).
// When omitted, DefaultCompression (-1) is used.
//
// Supported encodings (all from the standard library):
//   - gzip   (compress/gzip)
//   - deflate (compress/zlib per RFC 2616)
//
// Compress 返回一个压缩响应体的中间件。编码通过请求的 Accept-Encoding 头协商；
// 当客户端不接受任何支持的编码时，响应不压缩。
//
// 可指定可选的压缩级别（0-9，或 gzip 包的常量和
// NoCompression / BestSpeed / DefaultCompression / BestCompression）。
// 省略时使用 DefaultCompression (-1)。
//
// 支持的编码（均来自标准库）：
//   - gzip    (compress/gzip)
//   - deflate (compress/zlib，遵循 RFC 2616)
func Compress(level ...int) func(next http.Handler) http.Handler {
	lvl := flate.DefaultCompression
	if len(level) > 0 {
		lvl = level[0]
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			enc := negotiateEncoding(r.Header.Get("Accept-Encoding"))
			if enc == "" {
				next.ServeHTTP(w, r)
				return
			}

			cw := &compressResponseWriter{
				ResponseWriter: w,
				encoding:       enc,
				level:          lvl,
			}
			defer cw.Close()

			next.ServeHTTP(cw, r)
		})
	}
}

// --- encoding negotiation ---------------------------------------------------

// negotiateEncoding picks the first supported encoding that appears in the
// Accept-Encoding header. Returns "" when nothing matches.
//
// negotiateEncoding 选择 Accept-Encoding 中第一个出现的受支持编码。
// 无匹配时返回 ""。
func negotiateEncoding(ae string) string {
	for _, token := range strings.FieldsFunc(ae, func(r rune) bool { return r == ',' || r == ';' }) {
		t := strings.TrimSpace(token)
		switch strings.ToLower(t) {
		case "gzip":
			return "gzip"
		case "deflate":
			return "deflate"
		case "identity":
			return ""
		}
	}
	return ""
}

// --- response writer wrapper ------------------------------------------------

type compressResponseWriter struct {
	http.ResponseWriter
	encoding string
	level    int
	writer   io.WriteCloser
	buf      *bufio.Writer
	once     sync.Once
}

func (cw *compressResponseWriter) WriteHeader(code int) {
	cw.prepareHeaders()
	cw.ResponseWriter.WriteHeader(code)
}

func (cw *compressResponseWriter) prepareHeaders() {
	cw.once.Do(func() {
		cw.ResponseWriter.Header().Del("Content-Length")
		cw.ResponseWriter.Header().Set("Content-Encoding", cw.encoding)
	})
}

func (cw *compressResponseWriter) initWriter() {
	switch cw.encoding {
	case "gzip":
		cw.writer, _ = gzip.NewWriterLevel(cw.ResponseWriter, cw.level)
	case "deflate":
		cw.writer, _ = zlib.NewWriterLevel(cw.ResponseWriter, cw.level)
	}
	cw.buf = bufio.NewWriterSize(cw.writer, 4096)
}

func (cw *compressResponseWriter) Write(b []byte) (int, error) {
	cw.prepareHeaders()
	if cw.buf == nil {
		cw.initWriter()
	}
	return cw.buf.Write(b)
}

func (cw *compressResponseWriter) Flush() {
	if cw.buf == nil {
		return
	}
	cw.buf.Flush()
	if gw, ok := cw.writer.(*gzip.Writer); ok {
		gw.Flush()
	}
	if f, ok := cw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (cw *compressResponseWriter) Close() error {
	if cw.buf == nil {
		return nil
	}
	if err := cw.buf.Flush(); err != nil {
		return err
	}
	return cw.writer.Close()
}
