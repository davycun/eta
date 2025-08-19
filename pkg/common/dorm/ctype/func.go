package ctype

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	gormFieldCache = map[reflect.Type]map[string]string{} //reflect.Type -> []string
)

func NewInt(i int64) Integer {
	if i == 0 {
		return Integer{}
	}
	return Integer{Valid: true, Data: i}
}
func NewIntPrt(i int64) *Integer {
	if i == 0 {
		return &Integer{}
	}
	return &Integer{Valid: true, Data: i}
}
func NewFloat(i float64) Float {
	if i == 0 {
		return Float{}
	}
	return Float{Valid: true, Data: i}
}
func NewFloatPrt(i float64) *Float {
	if i == 0 {
		return &Float{}
	}
	return &Float{Valid: true, Data: i}
}
func NewString(str string, valid bool) String {
	return String{Valid: valid, Data: str}
}
func NewStringPrt(str string) *String {
	if str == "" {
		return &String{}
	}
	return &String{Valid: true, Data: str}
}
func NewText(str string) Text {
	if str == "" {
		return Text{}
	}
	return Text{Valid: true, Data: str}
}
func NewTextPrt(str string) *Text {
	if str == "" {
		return &Text{}
	}
	return &Text{Valid: true, Data: str}
}
func NewBoolean(bl bool, valid bool) Boolean {
	return Boolean{Valid: valid, Data: bl}
}
func NewBooleanPrt(bl bool) *Boolean {
	return &Boolean{Valid: true, Data: bl}
}
func NewInt64Array(dt ...int64) Int64Array {
	if len(dt) < 1 {
		return Int64Array{}
	}
	return Int64Array{Data: dt, Valid: true}
}
func NewInt64ArrayPrt(dt ...int64) *Int64Array {
	if len(dt) < 1 {
		return &Int64Array{}
	}
	return &Int64Array{Data: dt, Valid: true}
}
func NewStringArray(dt ...string) StringArray {
	if len(dt) < 1 {
		return StringArray{}
	}
	return StringArray{Data: dt, Valid: true}
}
func NewStringArrayPrt(dt ...string) *StringArray {
	if len(dt) < 1 {
		return &StringArray{}
	}
	return &StringArray{Data: dt, Valid: true}
}
func NewJson(obj any) Json {
	if obj == nil {
		return Json{}
	}
	return Json{Data: obj, Valid: true}
}
func NewJsonPrt(obj any) *Json {
	j := NewJson(obj)
	return &j
}
func NewLocalTime(tm time.Time) LocalTime {
	if tm.IsZero() {
		return LocalTime{}
	}
	return LocalTime{Data: tm, Valid: true}
}
func NewLocalTimePrt(tm time.Time) *LocalTime {
	if tm.IsZero() {
		return &LocalTime{}
	}
	return &LocalTime{Data: tm, Valid: true}
}

func Len[T StringArray | Int64Array | *StringArray | *Int64Array](t T) int {

	switch x := any(t).(type) {
	case StringArray:
		return len(x.Data)
	case Int64Array:
		return len(x.Data)
	case *StringArray:
		if x == nil {
			return 0
		}
		return len(x.Data)
	case *Int64Array:
		if x == nil {
			return 0
		}
		return len(x.Data)
	}
	return 0
}
func ToSliceString[T StringArray | *StringArray](t ...T) []string {
	rs := make([]string, 0)
	for _, v := range t {
		if Len(v) > 0 {
			switch x := any(v).(type) {
			case StringArray:
				rs = append(rs, x.Data...)
			case *StringArray:
				if rs != nil {
					rs = append(rs, x.Data...)
				}
			}
		}
	}
	return rs
}
func ToSliceInt[T Int64Array | *Int64Array](t T) []int64 {
	rs := make([]int64, 0)
	if Len(t) > 0 {
		switch x := any(t).(type) {
		case Int64Array:
			return x.Data
		case *Int64Array:
			return x.Data
		}
	}
	return rs
}

func ToInt64[T any](t T) int64 {
	switch s := any(t).(type) {
	case Integer:
		return s.Data
	case *Integer:
		if s == nil {
			return 0
		}
		return s.Data
	case float64:
		return int64(s)
	case float32:
		return int64(s)
	case int:
		return int64(s)
	case int8:
		return int64(s)
	case int16:
		return int64(s)
	case int32:
		return int64(s)
	case int64:
		return s
	}
	return 0
}
func ToFloat[T any](t T) float64 {
	switch s := any(t).(type) {
	case Integer:
		return float64(s.Data)
	case *Integer:
		if s == nil {
			return 0
		}
		return float64(s.Data)
	case Float:
		return s.Data
	case *Float:
		if s == nil {
			return 0
		}
		return s.Data
	case float64:
		return s
	case float32:
		return float64(s)
	case int:
		return float64(s)
	case int8:
		return float64(s)
	case int16:
		return float64(s)
	case int32:
		return float64(s)
	case int64:
		return float64(s)
	case uint:
		return float64(s)
	case uint8:
		return float64(s)
	case uint16:
		return float64(s)
	case uint32:
		return float64(s)
	case uint64:
		return float64(s)
	}
	return 0
}
func Bool[T any](dst T) bool {
	switch x := any(dst).(type) {
	case Boolean:
		return x.Valid && x.Data
	case *Boolean:
		if x == nil {
			return false
		}
		return x.Valid && x.Data
	case string:
		return strings.TrimSpace(strings.ToLower(x)) == "true"
	case bool:
		return x
	}
	return !reflect.ValueOf(dst).IsZero()
}

func ToString[T any](str T) string {
	switch s := any(str).(type) {
	case String:
		return s.Data
	case *String:
		if s == nil {
			return ""
		}
		return s.Data
	case Integer:
		if IsValid(s) {
			return fmt.Sprintf("%d", s.Data)
		}
		return ""
	case *Integer:
		if IsValid(s) {
			return fmt.Sprintf("%d", s.Data)
		}
		return ""
	case Float:
		if IsValid(s) {
			return strconv.FormatFloat(s.Data, 'f', -1, 64)
		}
		return ""
	case *Float:
		if IsValid(s) {
			return strconv.FormatFloat(s.Data, 'f', -1, 64)
		}
		return ""
	case Text:
		return s.Data
	case *Text:
		return s.Data
	case LocalTime:
		if IsValid(s) {
			return s.Data.Format(time.RFC3339Nano)
		}
		return ""
	case *LocalTime:
		if IsValid(s) {
			return s.Data.Format(time.RFC3339Nano)
		}
		return ""
	case Boolean:
		return fmt.Sprintf("%t", s.Data)
	case *Boolean:
		return fmt.Sprintf("%t", s.Data)
	case time.Time:
		return s.Format(time.RFC3339Nano)
	case *time.Time:
		if s == nil {
			return ""
		}
		return s.Format(time.RFC3339Nano)
	case StringArray:
		if IsValid(s) {
			return strings.Join(s.Data, ",")
		}
		return ""
	case *StringArray:
		if IsValid(s) {
			return strings.Join(s.Data, ",")
		}
		return ""
	case Int64Array:
		if IsValid(s) {
			s1 := make([]string, 0)
			for _, v := range s.Data {
				s1 = append(s1, fmt.Sprintf("%d", v))
			}
			return strings.Join(s1, ",")
		}
		return ""
	case *Int64Array:
		if IsValid(s) {
			s1 := make([]string, 0)
			for _, v := range s.Data {
				s1 = append(s1, fmt.Sprintf("%d", v))
			}
			return strings.Join(s1, ",")
		}
		return ""
	case Json:
		if IsValid(s) {
			dt, _ := json.Marshal(s.Data)
			return string(dt)
		}
		return ""
	case *Json:
		if IsValid(s) {
			dt, _ := json.Marshal(s.Data)
			return string(dt)
		}
		return ""
	case *string:
		if s == nil {
			return ""
		}
		return *s
	case string:
		return s
	case []byte:
		return string(s)
	case int8, int16, int32, int64, int, uint8, uint16, uint, uint32, uint64:
		return fmt.Sprintf("%d", s)
	default:
		return fmt.Sprintf("%v", str)
	}
}
func Concat[T String | string | *String](sep string, dt ...T) *String {
	if len(dt) < 1 {
		return NewStringPrt("")
	}
	bd := strings.Builder{}
	hasPre := false
	for _, v := range dt {
		var (
			tmp string
		)

		switch x := any(v).(type) {
		case string:
			tmp = x
		case String:
			tmp = x.Data
		case *String:
			if x != nil {
				tmp = x.Data
			}
		default:
			tmp = ""
		}
		if tmp == "" {
			continue
		}
		if hasPre && sep != "" {
			bd.WriteString(sep)
		}
		bd.WriteString(tmp)
		hasPre = true
	}
	return NewStringPrt(bd.String())
}

func IsValid(t any) bool {
	switch x := t.(type) {
	case String:
		return x.Valid
	case *String:
		if x != nil {
			return x.Valid
		}
		return false
	case Integer:
		return x.Valid
	case *Integer:
		if x != nil {
			return x.Valid
		}
		return false
	case Float:
		return x.Valid
	case *Float:
		if x != nil {
			return x.Valid
		}
		return false
	case Boolean:
		return x.Valid
	case *Boolean:
		if x != nil {
			return x.Valid
		}
		return false
	case Int64Array:
		return x.Valid
	case *Int64Array:
		if x != nil {
			return x.Valid
		}
		return false
	case StringArray:
		return x.Valid
	case *StringArray:
		if x != nil {
			return x.Valid
		}
		return false
	case Geometry:
		return x.Valid
	case *Geometry:
		if x != nil {
			return x.Valid
		}
		return false
	case Json:
		return x.Valid
	case *Json:
		if x != nil {
			return x.Valid
		}
		return false
	case LocalTime:
		return x.Valid
	case *LocalTime:
		if x != nil {
			return x.Valid
		}
		return false

	}
	return false
}
func EqualsString(src *String, dest *String) bool {
	if src == nil || !src.Valid || dest == nil || !dest.Valid {
		return false
	}

	return src.Data == dest.Data
}

func ToSliceMap(dt any) []map[string]interface{} {
	if x, ok := dt.([]map[string]interface{}); ok {
		return x
	}
	return []map[string]interface{}{}
}
func ToSlice(dt any) []interface{} {
	if x, ok := dt.([]interface{}); ok {
		return x
	}
	return []interface{}{}
}
func ToMap(dt any) map[string]interface{} {
	if x, ok := dt.(map[string]interface{}); ok {
		return x
	}
	return map[string]interface{}{}
}

func GetMapValue(dt map[string]interface{}, key string) interface{} {
	if key == "" || dt == nil {
		return nil
	}
	fields := strings.Split(key, ".")
	if len(fields) == 1 {
		return dt[fields[0]]
	}
	child := dt[fields[0]]

	if x, ok := child.(map[string]interface{}); ok {
		return GetMapValue(x, strings.Join(fields[1:], "."))
	}
	return child
}
