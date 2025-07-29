package utils

import (
	"github.com/davycun/eta/pkg/common/logger"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"regexp"
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
func ContainAnyInsensitive(str []string, target ...string) bool {
	if len(str) < 1 || len(target) < 1 {
		return false
	}
	mp := make(map[string]string)
	for i, _ := range str {
		v := strings.ToLower(str[i])
		mp[v] = v
	}

	for _, v := range target {
		_, ok := mp[strings.ToLower(v)]
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

// IsMatchedUri
// pattern示例: 以method@url格式，其中method表示请求方法,如果有多个用逗号隔开（匹配所有可以用*来代替），url表示的是url请求的path（支持正则表达式）
// 示例：
// 1. POST,GET@/api/a/b，表示针对/api/a/b的GET请求进行缓存
// 2. GET@/api/a/.*，表示以/api/a/开头的所有GET请求都缓存
// 3. GET@.*，表示所有GET请求都缓存
// 4. *@.*, 表示所有请求都缓存
// 5. 以URI匹配为优先级，比如配置里面有["GET@/api/a/*","*@.*"]，
// 6. 针对/api/a/b的POST请求，优先匹配到了"GET@/api/a/*"的URI，但是这个规则是只能针对GET请求，所以/api/a/b的POST请求不匹配，尽管patters中有"*@*"
func IsMatchedUri(method string, uri string, patterns ...string) bool {

	for _, pattern := range patterns {
		ci := strings.Split(pattern, "@")
		//如果配置不是method@pattern的格式，直接忽略
		if len(ci) != 2 {
			logger.Warnf("[%s] is not uri pattern ", pattern)
			continue
		}
		matched, err := regexp.MatchString(ci[1], uri)
		if err != nil {
			logger.Errorf("pattern [%s] match uri [%s] err %s", pattern, uri, err)
			return false
		}
		// 如果uri匹配成功即matched为true，就以当前的规则来判断是否匹配成功，
		// 如果没有匹配成功即matched为false，就继续往后匹配
		//ci[0] 可能是逗号隔开的多个
		if matched {
			return ci[0] == "*" || method == "" || strings.Contains(strings.ToLower(ci[0]), strings.ToLower(method))
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
