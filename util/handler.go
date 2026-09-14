package util

import (
	"net/http"
	"reflect"

	"github.com/rambollwong/rainbowferret/types"
)

// HandleT wraps a generic HandlerFunc into a standard http.HandlerFunc.
// It automatically decodes the request body into T, invokes the handler
// function, and writes the result as JSON. When the handler function returns
// a nil value it responds with 204 No Content instead of a JSON null body.
//
// When the response implements types.SuccessStatuser, its SuccessStatus() code
// is used instead of the default 200 OK (e.g. 201 Created).
//
// HandleT 将泛型 HandlerFunc 包装为标准 http.HandlerFunc。
// 它自动将请求体解码为 T，调用处理函数，并将结果写为 JSON。
// 当处理函数返回 nil 值时，返回 204 No Content 而非 JSON null 响应体。
//
// 当响应实现了 types.SuccessStatuser 时，会使用其 SuccessStatus() 返回的
// 状态码替代默认的 200 OK（如 201 Created）。
func HandleT[T, R any](handlerFn types.HandlerFunc[T, R]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req T
		if err := bindRequest(r, &req); err != nil {
			e, ok := err.(*types.HTTPError)
			if !ok {
				e = types.BadRequest(err.Error())
			}
			WriteJSON(w, e.Code, e)
			return
		}
		// do something with req
		ctx := r.Context()
		res, err := handlerFn(ctx, req)
		if err != nil {
			if e, ok := err.(*types.HTTPError); ok {
				WriteJSON(w, e.Code, e)
			} else {
				e := types.Internal(err.Error())
				WriteJSON(w, e.Code, e)
			}
			return
		}
		// then write response to w
		if IsNil(res) {
			WriteNoContent(w)
		} else {
			code := http.StatusOK
			if ss, ok := any(res).(types.SuccessStatuser); ok {
				code = ss.SuccessStatus()
			}
			WriteJSON(w, code, res)
		}
	}
}

// bindRequest binds the request into target. When T is a pointer type it first
// allocates the pointed-to zero value so Bind receives a pointer to a struct
// (not a pointer to a pointer), which Bind's form/param filling requires.
//
// bindRequest 将请求绑定到 target。当 T 是指针类型时，先分配其指向的零值，
// 使 Bind 接收到指向结构体的指针（而非指向指针的指针），满足 Bind 的表单/
// 参数填充要求。
func bindRequest[T any](r *http.Request, target *T) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() == reflect.Pointer && !rv.IsNil() {
		elem := rv.Elem()
		if elem.Kind() == reflect.Pointer {
			if elem.IsNil() {
				elem.Set(reflect.New(elem.Type().Elem()))
			}
			return Bind(r, elem.Interface())
		}
	}
	return Bind(r, target)
}
