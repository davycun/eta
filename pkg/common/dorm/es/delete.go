package es

import (
	"context"
	"github.com/davycun/eta/pkg/common/dorm/es/es_api"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"reflect"
	"time"
)

// Delete
// dest 可以是切片或者切片指针
// dest可以是id的切片，也可以是实体对象的切片，但是实体对象要包含ID字段，也可以是map的切片，但是map里也需要id字段
// 最终会根据ID进行删除
func Delete(api *es_api.Api, idx string, dest any) error {

	if api == nil || api.EsApi == nil || idx == "" {
		return nil
	}

	var (
		err   error
		start = time.Now()
		bulk  = api.EsTypedApi.Bulk()
	)

	val := getValueFromPointer(reflect.ValueOf(dest))

	for _, v := range val {
		id := GetStringValue(v, entity.IdFieldName)
		co := types.DeleteOperation{Index_: &idx, Id_: &id}
		err = bulk.DeleteOp(co)
		if err != nil {
			return err
		}
	}

	resp, err := bulk.Do(context.Background())

	if err != nil {
		LatencyLog(start, idx, optDelete, utils.StringToBytes(err.Error()), 500)
		return err
	}

	msg := getBulkErrorMsg(resp, optDelete)
	LatencyLog(start, idx, optDelete, utils.StringToBytes(msg), getBulkResultCode(resp))

	if resp != nil && resp.Errors {
		return errs.NewServerError(msg)
	}
	return nil
}
