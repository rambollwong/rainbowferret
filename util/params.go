package util

import (
	"net/http"
	"strconv"
)

// --- Query parameters --------------------------------------------------------

// QueryParam returns the first value for the named query parameter, or "".
// QueryParam 返回指定查询参数的第一个值，不存在则返回 ""。
func QueryParam(r *http.Request, name string) string { return r.URL.Query().Get(name) }

// QueryParamDefault returns the query parameter value or defaultVal when the
// parameter is missing.
// QueryParamDefault 返回查询参数的值，参数不存在时返回 defaultVal。
func QueryParamDefault(r *http.Request, name, defaultVal string) string {
	if q := r.URL.Query(); q.Has(name) {
		return q.Get(name)
	}
	return defaultVal
}

// QueryParamInt parses the named query parameter as int.
// QueryParamInt 将指定查询参数解析为 int。
func QueryParamInt(r *http.Request, name string) (int, error) {
	return strconv.Atoi(r.URL.Query().Get(name))
}

// QueryParamIntDefault returns the query parameter as int, or defaultVal when
// the parameter is missing or cannot be parsed.
// QueryParamIntDefault 返回查询参数的 int 值，参数不存在或解析失败时返回 defaultVal。
func QueryParamIntDefault(r *http.Request, name string, defaultVal int) int {
	v, err := strconv.Atoi(r.URL.Query().Get(name))
	if err != nil {
		return defaultVal
	}
	return v
}

// QueryParamInt64 parses the named query parameter as int64.
// QueryParamInt64 将指定查询参数解析为 int64。
func QueryParamInt64(r *http.Request, name string) (int64, error) {
	return strconv.ParseInt(r.URL.Query().Get(name), 10, 64)
}

// QueryParamInt64Default returns the query parameter as int64, or defaultVal
// when the parameter is missing or cannot be parsed.
// QueryParamInt64Default 返回查询参数的 int64 值，参数不存在或解析失败时返回 defaultVal。
func QueryParamInt64Default(r *http.Request, name string, defaultVal int64) int64 {
	v, err := strconv.ParseInt(r.URL.Query().Get(name), 10, 64)
	if err != nil {
		return defaultVal
	}
	return v
}

// QueryParamFloat parses the named query parameter as float64.
// QueryParamFloat 将指定查询参数解析为 float64。
func QueryParamFloat(r *http.Request, name string) (float64, error) {
	return strconv.ParseFloat(r.URL.Query().Get(name), 64)
}

// QueryParamFloatDefault returns the query parameter as float64, or defaultVal
// when the parameter is missing or cannot be parsed.
// QueryParamFloatDefault 返回查询参数的 float64 值，参数不存在或解析失败时返回 defaultVal。
func QueryParamFloatDefault(r *http.Request, name string, defaultVal float64) float64 {
	v, err := strconv.ParseFloat(r.URL.Query().Get(name), 64)
	if err != nil {
		return defaultVal
	}
	return v
}

// QueryParamBool parses the named query parameter as bool. Accepted values
// are those recognised by strconv.ParseBool: "1", "t", "T", "true", "TRUE",
// "True", "0", "f", "F", "false", "FALSE", "False".
// QueryParamBool 将指定查询参数解析为 bool。可接受的值由 strconv.ParseBool 定义。
func QueryParamBool(r *http.Request, name string) (bool, error) {
	return strconv.ParseBool(r.URL.Query().Get(name))
}

// QueryParamBoolDefault returns the query parameter as bool, or defaultVal when
// the parameter is missing or cannot be parsed.
// QueryParamBoolDefault 返回查询参数的 bool 值，参数不存在或解析失败时返回 defaultVal。
func QueryParamBoolDefault(r *http.Request, name string, defaultVal bool) bool {
	v, err := strconv.ParseBool(r.URL.Query().Get(name))
	if err != nil {
		return defaultVal
	}
	return v
}

// --- Path parameters (Go 1.22+ r.PathValue) ----------------------------------

// PathParam returns the value for the named path wildcard.
// PathParam 返回指定路径通配符的值。
func PathParam(r *http.Request, name string) string { return r.PathValue(name) }

// PathParamInt parses the named path parameter as int.
// PathParamInt 将指定路径参数解析为 int。
func PathParamInt(r *http.Request, name string) (int, error) {
	return strconv.Atoi(r.PathValue(name))
}

// PathParamIntDefault returns the path parameter as int, or defaultVal when
// the parameter is missing or cannot be parsed.
// PathParamIntDefault 返回路径参数的 int 值，参数不存在或解析失败时返回 defaultVal。
func PathParamIntDefault(r *http.Request, name string, defaultVal int) int {
	v, err := strconv.Atoi(r.PathValue(name))
	if err != nil {
		return defaultVal
	}
	return v
}

// PathParamInt64 parses the named path parameter as int64.
// PathParamInt64 将指定路径参数解析为 int64。
func PathParamInt64(r *http.Request, name string) (int64, error) {
	return strconv.ParseInt(r.PathValue(name), 10, 64)
}

// PathParamInt64Default returns the path parameter as int64, or defaultVal when
// the parameter is missing or cannot be parsed.
// PathParamInt64Default 返回路径参数的 int64 值，参数不存在或解析失败时返回 defaultVal。
func PathParamInt64Default(r *http.Request, name string, defaultVal int64) int64 {
	v, err := strconv.ParseInt(r.PathValue(name), 10, 64)
	if err != nil {
		return defaultVal
	}
	return v
}

// PathParamFloat parses the named path parameter as float64.
// PathParamFloat 将指定路径参数解析为 float64。
func PathParamFloat(r *http.Request, name string) (float64, error) {
	return strconv.ParseFloat(r.PathValue(name), 64)
}

// PathParamFloatDefault returns the path parameter as float64, or defaultVal
// when the parameter is missing or cannot be parsed.
// PathParamFloatDefault 返回路径参数的 float64 值，参数不存在或解析失败时返回 defaultVal。
func PathParamFloatDefault(r *http.Request, name string, defaultVal float64) float64 {
	v, err := strconv.ParseFloat(r.PathValue(name), 64)
	if err != nil {
		return defaultVal
	}
	return v
}

// PathParamBool parses the named path parameter as bool. Accepted values are
// those recognised by strconv.ParseBool.
// PathParamBool 将指定路径参数解析为 bool。
func PathParamBool(r *http.Request, name string) (bool, error) {
	return strconv.ParseBool(r.PathValue(name))
}

// PathParamBoolDefault returns the path parameter as bool, or defaultVal when
// the parameter is missing or cannot be parsed.
// PathParamBoolDefault 返回路径参数的 bool 值，参数不存在或解析失败时返回 defaultVal。
func PathParamBoolDefault(r *http.Request, name string, defaultVal bool) bool {
	v, err := strconv.ParseBool(r.PathValue(name))
	if err != nil {
		return defaultVal
	}
	return v
}
