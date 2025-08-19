package entity

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"reflect"
)

type Setter interface {
	Set(field string, val any)
}

// Set
// src 需要传入结构体指针，否则set无效
func Set(src interface{}, key string, val interface{}) error {

	if x, ok := src.(Setter); ok {
		x.Set(key, val)
		return nil
	}
	if x, ok := src.(reflect.Value); ok {
		return SetValue(x, key, val)
	}
	return SetValue(reflect.ValueOf(src), key, val)
}

func SetValue(target reflect.Value, fieldName string, value any) error {
	if !target.IsValid() {
		return nil
	}

	switch target.Kind() {
	case reflect.Pointer:
		return SetValue(target.Elem(), fieldName, value)
	case reflect.Struct:
		field := GetValue(target, fieldName)
		if !field.IsValid() {
			field = target.FieldByName(utils.Column2StructFieldName(fieldName))
		}
		if !field.IsValid() || !field.CanInterface() {
			return nil
		}
		valFieldInter := field.Interface()
		if field.Kind() == reflect.Pointer && !field.CanSet() {
			field = field.Elem()
		}
		switch valFieldInter.(type) {
		case ctype.String:
			field.Set(reflect.ValueOf(ctype.NewString(ctype.ToString(value), true)))
		case *ctype.String:
			field.Set(reflect.ValueOf(ctype.NewStringPrt(ctype.ToString(value))))
		case ctype.Text:
			field.Set(reflect.ValueOf(ctype.NewText(ctype.ToString(value))))
		case *ctype.Text:
			field.Set(reflect.ValueOf(ctype.NewTextPrt(ctype.ToString(value))))
		case ctype.Integer:
			field.Set(reflect.ValueOf(ctype.NewInt(ctype.ToInt64(value))))
		case *ctype.Integer:
			field.Set(reflect.ValueOf(ctype.NewIntPrt(ctype.ToInt64(value))))
		case ctype.Float:
			field.Set(reflect.ValueOf(ctype.NewFloat(ctype.ToFloat(value))))
		case *ctype.Float:
			field.Set(reflect.ValueOf(ctype.NewFloatPrt(ctype.ToFloat(value))))
		case ctype.Boolean:
			field.Set(reflect.ValueOf(ctype.NewBoolean(ctype.Bool(value), true)))
		case *ctype.Boolean:
			field.Set(reflect.ValueOf(ctype.NewBooleanPrt(ctype.Bool(value))))
		case ctype.Json:
			field.Set(reflect.ValueOf(ctype.NewJson(value)))
		case *ctype.Json:
			field.Set(reflect.ValueOf(ctype.NewJsonPrt(value)))
		default:
			field.Set(reflect.ValueOf(value))
		}

	case reflect.Map:
		if !target.IsValid() || target.IsZero() || !target.CanInterface() {
			return nil
		}
		//TODO 理论上要根据jsonTag或者gormTag来确定名字
		jsonKey := utils.HumpToUnderline(fieldName)
		valInter := target.Interface()
		switch x := valInter.(type) {
		case ctype.Map:
			x[jsonKey] = value
		case *ctype.Map:
			(*x)[jsonKey] = value
		case gin.H:
			x[jsonKey] = value
		case *gin.H:
			(*x)[jsonKey] = value
		case map[string]any:
			x[jsonKey] = value
		case *map[string]any:
			(*x)[jsonKey] = value
		}
	default:
		return errs.NewServerError("the target is not a struct or a map when use SetValue")
	}
	return nil
}
func setValue(src reflect.Value, target any) {

}
