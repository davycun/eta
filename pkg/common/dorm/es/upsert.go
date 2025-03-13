package es

import (
	"bytes"
	"context"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/es/es_api"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	"reflect"
	"time"
)

// Upsert
// dest 可以是切片或者切片指针
// dest 是具体要操作的文档的切片，文档内容需要包含ID字段
// 如果对应的id数据已经存在则更新，不存在则插入
func Upsert(api *es_api.Api, idx string, dest any) error {

	if api == nil || api.EsApi == nil || idx == "" {
		return nil
	}

	var (
		err  error
		bulk = api.EsTypedApi.Bulk()
	)

	val := getValueFromPointer(reflect.ValueOf(dest))

	for _, v := range val {
		id := GetStringValue(v, entity.IdFieldName)
		co := types.IndexOperation{Index_: &idx, Id_: &id}
		err = bulk.IndexOp(co, v.Interface())
		if err != nil {
			return err
		}
	}

	start := time.Now()
	resp, err := bulk.Refresh(refresh.True).Do(context.Background())
	if err != nil {
		LatencyLog(start, idx, optUpsert, utils.StringToBytes(err.Error()), 500)
		return err
	}

	reqBody := bytes.Buffer{}
	msg := getBulkErrorMsg(resp, optUpsert)
	reqBody.WriteString(msg)
	LatencyLog(start, idx, optUpsert, reqBody.Bytes(), getBulkResultCode(resp))

	if resp != nil && resp.Errors {
		return errs.NewServerError(msg)
	}
	return nil
}

func getErrorMsg(er *types.ErrorCause) string {
	if er.CausedBy != nil && er.CausedBy.Reason != nil {
		return ctype.ToString(er.Reason) + ", cause by " + getErrorMsg(er.CausedBy)
	}
	return ctype.ToString(er.Reason)
}

func getValueFromPointer(val reflect.Value) []reflect.Value {
	switch val.Kind() {
	case reflect.Pointer:
		return getValueFromPointer(val.Elem())
	case reflect.Slice:
		vs := make([]reflect.Value, 0, val.Len())
		for i := 0; i < val.Len(); i++ {
			vs = append(vs, val.Index(i))
		}
		return vs
	default:
		return []reflect.Value{val}
	}
}

func GetStringValue(val reflect.Value, field string) string {

	if !val.IsValid() {
		return ""
	}

	switch val.Kind() {
	case reflect.Pointer:
		return GetStringValue(val.Elem(), field)
	case reflect.String:
		return val.String()
	case reflect.Struct:
		return GetStringValue(val.FieldByName(field), field)
	case reflect.Map:
		return GetStringValue(val.MapIndex(reflect.ValueOf(field)), field)
	default:

	}
	return ""
}
func GetInt64Value(val reflect.Value, field string) int64 {

	if !val.IsValid() {
		return 0
	}
	switch val.Kind() {
	case reflect.Pointer:
		return GetInt64Value(val.Elem(), field)
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return val.Int()
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		return int64(val.Uint())
	case reflect.Struct:

		if val.CanInterface() {
			switch x := val.Interface().(type) {
			case ctype.Integer:
				return x.Data
			case *ctype.Integer:
				if x != nil {
					return x.Data
				}
				return 0
			}
		}

		return GetInt64Value(val.FieldByName(field), field)

	case reflect.Map:
		return GetInt64Value(val.MapIndex(reflect.ValueOf(field)), field)
	default:

	}
	return 0
}
