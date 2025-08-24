package data_test

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/http_tes"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/template"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

type Book struct {
	Name  string
	Price float64
}

func TestCallback(t *testing.T) {

	var (
		code = fmt.Sprintf("my_test_table_%d", time.Now().UnixMilli())
	)

	hook.AddModifyCallback(code, func(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
		return modifyCallback(t, cfg, pos)
	})

	tmp := template.Template{
		Code:   code,
		Status: "ready",
		Table: entity.Table{
			Fields: []entity.TableField{
				{
					Name:    "name",
					Title:   "姓名",
					Type:    "string",
					Comment: "姓名",
				},
				{
					Name:    "年龄",
					Title:   "年龄",
					Type:    "integer",
					Comment: "年龄",
				},
			},
		},
	}
	http_tes.Call(t, http_tes.HttpCase{
		Method: "POST",
		Path:   "/template/create",
		Body: dto.ModifyParam{
			Data: []template.Template{tmp},
		},
	})
	http_tes.Call(t, http_tes.HttpCase{
		Method: "POST",
		Path:   fmt.Sprintf(`/data/%s/create`, code),
		Body:   `{"data":[{"name":"davy","年龄":23}]}`,
	})

}

func modifyCallback(t *testing.T, cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	return hook.BeforeModify(cfg, pos, func(cfg *hook.SrvConfig, oldValues []reflect.Value, newValues []reflect.Value) error {
		for i, v := range newValues {
			if i == 0 {
				fd := v.FieldByName(utils.Column2StructFieldName("name"))
				nm := fd.Interface()
				if r, ok := nm.(ctype.String); ok {
					assert.Equal(t, "davy", r.Data)
					fd.Set(reflect.ValueOf(ctype.NewString("davy_新的值", true)))
				}
				if r, ok := nm.(*ctype.String); ok {
					assert.Equal(t, "davy", r.Data)
					fd.Set(reflect.ValueOf(ctype.NewStringPrt("davy_新的值")))
				}
				fd1 := v.FieldByName(utils.Column2StructFieldName("年龄"))
				nm1 := fd1.Interface()
				if r, ok := nm1.(ctype.Integer); ok {
					assert.Equal(t, int64(23), r.Data)
				}
				if r, ok := nm.(*ctype.Integer); ok {
					assert.Equal(t, int64(23), r.Data)
				}
			}
		}
		return nil
	}, iface.MethodCreate)

}

func TestTypeConvert(t *testing.T) {
	data := []Book{{Name: "b1", Price: 9.9}, {Name: "b2", Price: 8.8}}
	logger.Infof("%v", data)
	var res any
	res = data
	if s, ok := res.([]interface{}); ok {
		logger.Infof("[]interface{} :%v", s)
	}
	if s, ok := res.([]any); ok {
		logger.Infof("[]any :%v", s)
	}
	switch res.(type) {
	case []any:
		logger.Infof("[]any :%v", res)
	}
	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Slice:
		logger.Infof("slice :%v", v.Len())
	default:
		logger.Infof("not a slice")
	}

}
