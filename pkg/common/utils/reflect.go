package utils

import (
	"reflect"
	"strings"
)

func GetEntityName(tp reflect.Type) string {
	if tp == nil {
		return ""
	}
	tp = GetRealType(tp)
	tbStr := tp.String()

	tbNmList := strings.Split(tbStr, ".")

	for i, v := range tbNmList {
		tbNmList[i] = HumpToUnderline(v)
	}
	return strings.Join(tbNmList, "_")
}

func GetRealType(tp reflect.Type) reflect.Type {
	if tp == nil {
		return nil
	}
	switch tp.Kind() {
	case reflect.Pointer:
		return GetRealType(tp.Elem())
	default:
	}
	return tp
}
func GetRealValue(value reflect.Value) reflect.Value {

	switch value.Kind() {
	case reflect.Pointer:
		return GetRealValue(value.Elem())
	default:
	}
	return value
}
func NewPointer(tp reflect.Type, slice bool) any {
	if slice {
		return reflect.New(reflect.SliceOf(tp)).Interface()
	}
	return reflect.New(tp).Interface()
}
func PrtToReal(obj any) any {
	if obj == nil {
		return nil
	}
	if v, ok := obj.(reflect.Value); ok {
		return GetRealValue(v)
	}
	val := GetRealValue(reflect.ValueOf(obj))

	switch val.Kind() {
	case reflect.Interface, reflect.Pointer:
		return GetRealValue(val).Interface()
	default:
	}
	return val.Interface()
}

// ConvertToValueArray
// 主要是把一些obj转换成数组的Value形式，这样后续可以统一处理
func ConvertToValueArray(obj any) []reflect.Value {
	var (
		val reflect.Value
	)
	if x, ok := obj.(reflect.Value); ok {
		val = x
	} else {
		val = reflect.ValueOf(obj)
	}

	var vs []reflect.Value
	if !val.IsValid() {
		return vs
	}
	switch val.Kind() {
	case reflect.Pointer:
		return ConvertToValueArray(val.Elem())
	case reflect.Slice:
		vs = make([]reflect.Value, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			vs = append(vs, ConvertToValueArray(v)...)
		}
		return vs
	case reflect.Struct, reflect.Map:
		vs = make([]reflect.Value, 1)
		vs[0] = val
		return vs
	default:
	}
	return vs
}
