package util

import "reflect"

// IsNil reports whether v is nil, handling both untyped nil (any(nil)) and
// typed nil (e.g. a nil pointer, slice, map, channel, func, or interface).
// Non-nilable types (int, string, struct, etc.) always return false.
//
// IsNil 判断 v 是否为 nil，同时处理无类型 nil 和带类型 nil（如 nil 指针、
// slice、map、channel、func 或 interface）。不可为 nil 的类型（int、string、
// struct 等）始终返回 false。
func IsNil[T any](v T) bool {
	iface := any(v)
	if iface == nil {
		return true
	}
	rv := reflect.ValueOf(iface)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice,
		reflect.Map, reflect.Chan, reflect.Func:
		return rv.IsNil()
	}
	return false
}
