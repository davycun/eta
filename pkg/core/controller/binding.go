package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/validity"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"reflect"
	"strings"
	"time"
)

func BindBody(c *gin.Context, dst any) error {
	if dst == nil {
		return errors.New("BindBody dst can't be nil")
	}
	return c.ShouldBindBodyWith(dst, binding.JSON)
}

// BindBodyPartial 绑定之后只是校验指定的字段
// fields 如果有结构体属性，那么可以通过dot来指定嵌套结构体需要校验的字段
func BindBodyPartial(c *gin.Context, dst any, fields ...string) error {
	if dst == nil {
		return errors.New("BindBody dst can't be nil")
	}
	err := BindBodyWithoutValidate(c, dst)
	if err != nil {
		return err
	}
	if len(fields) > 0 {
		return ValidateStructFields(dst, false, fields...)
	}
	return nil
}

// BindBodyExcept 绑定body内容到dst
// 绑定成功后，除了指定的fields 之外，其余的所有字段都需要根据field的tag描述来进行合法性校验
func BindBodyExcept(c *gin.Context, dst any, fields ...string) error {
	if dst == nil {
		return errors.New("BindBody dst can't be nil")
	}
	err := BindBodyWithoutValidate(c, dst)
	if err != nil {
		return err
	}
	if len(fields) > 0 {
		return ValidateStructFields(dst, true, fields...)
	}
	return nil
}

func BindQuery(c *gin.Context, dst any) error {
	if dst == nil {
		return errors.New("BindQuery dst can't be nil")
	}
	err := c.ShouldBindQuery(dst)
	if err != nil {
		return err
	}
	return nil
}

func BindUri(c *gin.Context, dst any) error {
	if dst == nil {
		logger.Warn("BindUri argus dst is nil")
		return errors.New("dst is nil")
	}
	err := c.ShouldBindUri(dst)
	if err != nil {
		return errors.New(fmt.Sprintf("bind uri error %s", err.Error()))
	}
	return nil
}

func BindBodyWithoutValidate(c *gin.Context, obj any) (err error) {
	var body []byte
	if cb, ok := c.Get(gin.BodyBytesKey); ok {
		if cbb, ok1 := cb.([]byte); ok1 {
			body = cbb
		}
	}
	if body == nil {
		body, err = io.ReadAll(c.Request.Body)
		if err != nil {
			return err
		}
		c.Set(gin.BodyBytesKey, body)
	}

	err = BindBodyByteWithoutValidate(body, obj)
	return err
}
func BindBodyByteWithoutValidate(body []byte, obj any) (err error) {
	r := bytes.NewReader(body)

	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(obj)
	return err
}

func validIgnore(obj any) bool {
	switch obj.(type) {
	case int8, int16, int32, int, int64, float32, float64, bool, string, byte, uint16, uint32, uint64, uint:
		return true
	case *int8, *int16, *int32, *int, *int64, *float32, *float64, *bool, *string, *byte, *uint16, *uint32, *uint64, *uint:
		return true
	case ctype.Int64Array, ctype.StringArray, ctype.Boolean, ctype.Float, ctype.Geometry, ctype.Integer, ctype.Json, ctype.Map, ctype.String, ctype.LocalTime:
		return true
	case *ctype.Int64Array, *ctype.StringArray, *ctype.Boolean, *ctype.Float, *ctype.Geometry, *ctype.Integer, *ctype.Json, *ctype.Map, *ctype.String, *ctype.LocalTime:
		return true
	case *time.Time, time.Time:
	case []string, []int, []int64, []int32, []float64, []float32, []bool, []byte:
		return true
	case *[]string, *[]int, *[]int64, *[]int32, *[]float64, *[]float32, *[]bool, *[]byte:
		return true
	}
	return false
}

func ValidateStructFields(obj any, except bool, fields ...string) error {
	if obj == nil || validIgnore(obj) {
		return nil
	}
	var (
		value = reflect.ValueOf(obj)
	)
	switch value.Kind() {
	case reflect.Pointer:
		if !value.Elem().IsValid() || !value.Elem().CanInterface() { //表示没有此字段
			return nil
		}
		return ValidateStructFields(value.Elem().Interface(), except, fields...)
	case reflect.Struct:
		currentFields, validChild := resolveFields2(fields...)
		var (
			err    error
			errSet = make(binding.SliceValidationError, 0)
		)
		validate := global.GetValidator()

		tp := value.Type()
		if except {
			exceptFields := make(map[string]string)
			for _, v := range currentFields {
				exceptFields[v] = v
			}
			currentFields = make([]string, 0, 10)
			for i := 0; i < value.NumField(); i++ {
				fd := tp.Field(i)
				tg := fd.Tag.Get(validity.ValidateTagName)

				_, ok := exceptFields[fd.Name]
				if ok {
					continue
				}
				if strings.Contains(tg, validity.IgnoreTagName) {
					//exceptFields[fd.Name] = fd.Name
					continue
				}
				currentFields = append(currentFields, fd.Name)
			}
		}

		err = validate.StructPartial(obj, currentFields...)

		if err != nil {
			errSet = append(errSet, err)
		}

		for i := 0; i < value.NumField(); i++ {
			val := value.Field(i)
			if !val.IsValid() && !val.CanInterface() { //表示没有此字段
				continue
			}
			switch val.Kind() {
			case reflect.Pointer, reflect.Struct, reflect.Array, reflect.Slice, reflect.Interface:
				fd := tp.Field(i)
				if strings.Contains(fd.Tag.Get(validity.ValidateTagName), validity.IgnoreTagName) || !fd.IsExported() {
					continue
				}
				err = ValidateStructFields(val.Interface(), except, validChild...)
			default:
				continue
			}
			if err != nil {
				errSet = append(errSet, err)
			}
		}
		if len(errSet) == 0 {
			return nil
		}
		return errSet

	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(binding.SliceValidationError, 0)
		for i := 0; i < count; i++ {
			tmp := value.Index(i)
			switch tmp.Kind() {
			case reflect.Pointer, reflect.Struct, reflect.Array, reflect.Slice, reflect.Interface:
				if err := ValidateStructFields(tmp.Interface(), except, fields...); err != nil {
					validateRet = append(validateRet, err)
				}
			default:
				continue
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

func resolveFields2(fields ...string) ([]string, []string) {
	var (
		child         = make([]string, 0, len(fields))
		current       = make([]string, 0, len(fields))
		childFields   = make(map[string]string, len(fields))
		currentFields = make(map[string]string, len(fields))
	)
	if len(fields) < 1 {
		return child, current
	}
	for i, v := range fields {
		before, after, found := strings.Cut(v, ".")
		if found {
			if before == "*" {
				childFields[v] = fields[i]
				currentFields[after] = after
			} else {
				currentFields[before] = before
				childFields[after] = after
			}
		} else {
			currentFields[v] = fields[i]
		}
	}

	for k, _ := range childFields {
		child = append(child, k)
	}
	for k, _ := range currentFields {
		current = append(current, k)
	}

	return current, child
}
func resolveFields(val reflect.Value, fields ...string) (validChild map[string][]string, currentFields []string) {
	validChild = make(map[string][]string, len(fields))
	currentFields = make([]string, 0, len(fields))
	if len(fields) < 1 {
		return
	}
	for _, v := range fields {
		before, after, found := strings.Cut(v, ".")
		if found {
			if before != "" {
				child, ok := validChild[before]
				if !ok {
					child = make([]string, 0, len(fields))
					validChild[before] = child
				}
				if after != "" {
					child = append(child, after)
					validChild[before] = child
				}
			}
		} else {
			currentFields = append(currentFields, v)
		}
	}
	return
}
