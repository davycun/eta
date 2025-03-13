package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/davycun/eta/pkg/core/controller"
)

// BindBodyExcept 绑定body内容到dst
// 绑定成功后，除了指定的fields 之外，其余的所有字段都需要根据field的tag描述来进行合法性校验
func BindBodyExcept(body []byte, dst any, fields ...string) error {
	if dst == nil {
		return errors.New("BindBody dst can't be nil")
	}
	err := BindBodyByteWithoutValidate(body, dst)
	if err != nil {
		return err
	}
	if len(fields) > 0 {
		return controller.ValidateStructFields(dst, true, fields...)
	}
	return nil
}

// BindBodyPartial 绑定之后只是校验指定的字段
// fields 如果有结构体属性，那么可以通过dot来指定嵌套结构体需要校验的字段
func BindBodyPartial(body []byte, dst any, fields ...string) error {
	if dst == nil {
		return errors.New("BindBody dst can't be nil")
	}
	err := BindBodyByteWithoutValidate(body, dst)
	if err != nil {
		return err
	}
	if len(fields) > 0 {
		return controller.ValidateStructFields(dst, false, fields...)
	}
	return nil
}

func BindBodyByteWithoutValidate(body []byte, obj any) (err error) {
	r := bytes.NewReader(body)
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(obj)
	return err
}
