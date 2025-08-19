package utils

import (
	"errors"
	"github.com/duke-git/lancet/v2/slice"
	"reflect"
)

func ConvertToSlice[T any](data any, dst *[]T) {

	var (
		dt []T
	)
	switch x := data.(type) {
	case *[]T:
		dt = *(x)
	case []T:
		dt = x
	case *T:
		dt = append(dt, *(x))
	case T:
		dt = append(dt, x)
	}
	*dst = dt
}

func AppendNoEmpty(src []string, dst ...string) []string {
	if len(dst) < 1 {
		return src
	}
	for _, v := range dst {
		if v != "" {
			src = append(src, v)
		}
	}
	return src
}

func Chunk(data any, size int) ([]any, error) {
	// data 切片
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Slice {
		return nil, errors.New("data must be a slice type")
	}
	if val.Len() == 0 {
		return nil, nil
	}

	var dataChunks []any
	start, end := 0, size
	for start < val.Len() {
		if end > val.Len() {
			end = val.Len()
		}
		dataChunks = append(dataChunks, val.Slice(start, end).Interface())
		start = start + size
		end = end + size
	}
	return dataChunks, nil
}

func SliceLen(data any) (int, error) {
	// data 切片
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Slice {
		return 0, errors.New("data must be a slice type")
	}
	return val.Len(), nil
}

func SliceRemoveElemByIndexes(data any, idx []int) (any, error) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Slice {
		return nil, errors.New("data must be a slice type")
	}

	newData := reflect.MakeSlice(val.Type(), 0, val.Len())
	if val.Len() == 0 {
		return newData.Interface(), nil
	}

	for i := range val.Len() {
		if slice.Contain(idx, i) {
			continue
		}
		newData = reflect.Append(newData, val.Index(i))
	}
	return newData.Interface(), nil
}

func IsEmptySlice(data any) bool {
	if data == nil {
		return true
	}
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if !slice.Contain([]reflect.Kind{reflect.Slice, reflect.Array}, val.Kind()) {
		return true
	}
	return val.Len() == 0
}
