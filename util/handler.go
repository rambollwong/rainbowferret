package util

import (
	"net/http"

	"github.com/rambollwong/rainbowferret/types"
)

// HandleT wraps a generic HandlerFunc into a standard http.HandlerFunc.
// It automatically decodes the request body into T, invokes the handler
// function, and writes the result as JSON. When the handler function returns
// a nil value it responds with 204 No Content instead of a JSON null body.
//
// HandleT 将泛型 HandlerFunc 包装为标准 http.HandlerFunc。
// 它自动将请求体解码为 T，调用处理函数，并将结果写为 JSON。
// 当处理函数返回 nil 值时，返回 204 No Content 而非 JSON null 响应体。
func HandleT[T any](handlerFn types.HandlerFunc[T]) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req T
		if err := Bind(r, &req); err != nil {
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
		if res == nil {
			WriteNoContent(w)
		} else {
			WriteJSON(w, http.StatusOK, res)
		}
	}
}
