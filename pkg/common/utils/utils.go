package utils

import (
	"github.com/davycun/eta/pkg/common/logger"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

const (
	TimeLayout          = "2006-01-02 15:04:05"
	TimeLayoutZone      = "2006-01-02 15:04:05Z07:00"
	TimeLayoutZoneMilli = "2006-01-02 15:04:05.999Z07:00"
	TimeLayoutZoneMicro = "2006-01-02 15:04:05.999999Z07:00"
	TimeLayoutZoneNano  = "2006-01-02 15:04:05.999999999Z07:00"
)

func CurrentTimeStr() string {

	return time.Now().Format(TimeLayout)
}

func IsZero(obj interface{}) bool {
	if obj == nil {
		return true
	}
	of := reflect.ValueOf(obj)
	return of.IsZero()
}
func IsNotZero(obj interface{}) bool {
	if obj == nil {
		return false
	}
	of := reflect.ValueOf(obj)
	return !of.IsZero()
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func BytesToString(b []byte) string {
	if len(b) < 1 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

// Interface2bytes 将接口转换为字节切片
func Interface2bytes(i interface{}) (b []byte) {
	switch v := i.(type) {
	case string:
		b = StringToBytes(v)
	case []byte:
		b = v
	}
	return
}

func Join(sep string, target ...int64) string {
	bd := strings.Builder{}
	for i, v := range target {
		if i > 0 {
			bd.WriteString(sep)
		}
		bd.WriteString(strconv.FormatInt(v, 10))
	}
	return bd.String()
}

// ContainAny str 里是否包含 target 中的任一一个
func ContainAny[T comparable](str []T, target ...T) bool {
	if len(str) < 1 || len(target) < 1 {
		return false
	}
	mp := make(map[T]T)
	for i, _ := range str {
		mp[str[i]] = str[i]
	}

	for _, v := range target {
		_, ok := mp[v]
		if ok {
			return true
		}
	}
	return false
}

// ContainAll elem 是否 allSlice 的子集
func ContainAll[T comparable](allSlice []T, elem ...T) bool {
	if len(elem) > len(allSlice) {
		return false
	}
	for _, e := range elem {
		if !ContainAny(allSlice, e) {
			return false
		}
	}
	return true
}

func Merge[T comparable](src []T, str ...T) []T {

	mp := make(map[T]T)
	for _, v := range src {
		mp[v] = v
	}
	for _, v := range str {
		switch x := any(v).(type) {
		case int8, int16, int32, int64, int, uint8, uint16, uint, uint32, uint64:
			if x == 0 {
				continue
			}
		case float32, float64:
			if x == 0 {
				continue
			}
		case string:
			if x == "" {
				continue
			}
		}
		_, ok := mp[v]
		if !ok {
			src = append(src, v)
			mp[v] = v
		}
	}
	return src
}

func IsMatchedUri(uri string, uris ...string) bool {
	if uris == nil || len(uris) < 1 {
		return false
	}
	for _, v := range uris {
		if strings.Contains(v, "*") {
			v = v[:strings.Index(v, "*")]
			if strings.HasPrefix(uri, v) {
				return true
			}
		}
		if v == uri {
			return true
		}
	}
	return false
}

// DifferenceOfStringSlices 使用map模拟集合求差集
func DifferenceOfStringSlices[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]bool)
	var diff []T

	// 将slice2中的元素添加到集合中
	for _, item := range slice2 {
		set[item] = true
	}

	// 遍历slice1，检查集合中是否存在该元素
	// 如果不存在，则是差集的一部分
	for _, item := range slice1 {
		if _, found := set[item]; !found {
			diff = append(diff, item)
		}
	}

	return diff
}

// IntersectionOfStringSlices 使用map模拟集合求交集
func IntersectionOfStringSlices[T comparable](slice1, slice2 []T) []T {
	// 创建一个 map 用于存储 slice1 中的元素
	m := make(map[T]bool)
	for _, v := range slice1 {
		m[v] = true
	}

	// 遍历 slice2，如果元素在 map 中存在，则添加到结果切片中
	var intersection []T
	for _, v := range slice2 {
		if m[v] {
			intersection = append(intersection, v)
		}
	}

	return intersection
}

func StructToMap(s interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	marshal, err := jsoniter.Marshal(s)
	if err != nil {
		logger.Warnf("convert struct to map failed, err: %v", err)
		return out
	}
	err = jsoniter.Unmarshal(marshal, &out)
	if err != nil {
		logger.Warnf("convert struct to map failed, err: %v", err)
	}
	return out
}
