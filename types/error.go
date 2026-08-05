package types

// HTTPError is an error that carries an HTTP status code.
// HTTPError 是一个携带 HTTP 状态码的错误。
type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

// NewHTTPError creates an HTTPError with the given code and message.
// NewHTTPError 使用指定的状态码和消息创建一个 HTTPError。
func NewHTTPError(code int, msg string) *HTTPError {
	return &HTTPError{Code: code, Message: msg}
}

// Error implements the error interface.
// Error 实现 error 接口。
func (e *HTTPError) Error() string { return e.Message }

// BadRequest returns a 400 HTTPError.
// BadRequest 返回一个 400 HTTPError。
func BadRequest(msg string) *HTTPError { return &HTTPError{400, msg} }

// NotFound returns a 404 HTTPError.
// NotFound 返回一个 404 HTTPError。
func NotFound(msg string) *HTTPError { return &HTTPError{404, msg} }

// Unauthorized returns a 401 HTTPError.
// Unauthorized 返回一个 401 HTTPError。
func Unauthorized(msg string) *HTTPError { return &HTTPError{401, msg} }

// UnsupportedMediaType returns a 415 HTTPError.
// UnsupportedMediaType 返回一个 415 HTTPError。
func UnsupportedMediaType(msg string) *HTTPError { return &HTTPError{415, msg} }

// Internal returns a 500 HTTPError.
// Internal 返回一个 500 HTTPError。
func Internal(msg string) *HTTPError { return &HTTPError{500, msg} }

// MethodNotAllowed returns a 405 HTTPError.
// MethodNotAllowed 返回一个 405 HTTPError。
func MethodNotAllowed(msg string) *HTTPError { return &HTTPError{405, msg} }
