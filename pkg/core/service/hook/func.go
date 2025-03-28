package hook

import (
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/iface"
	"reflect"
)

func BeforeCreate[T any](cfg *SrvConfig, pos CallbackPosition, f func(cfg *SrvConfig, newValues []T) error) error {
	if pos != CallbackBefore {
		return nil
	}
	return modifyCallback(cfg, func(cfg *SrvConfig, oldValues []T, newValues []T) error {
		return f(cfg, newValues)
	}, iface.MethodCreate)
}
func BeforeCreateAny(cfg *SrvConfig, pos CallbackPosition, fc func(cfg *SrvConfig) error) error {
	if pos != CallbackBefore {
		return nil
	}
	if cfg.Method != iface.MethodCreate {
		return nil
	}
	return fc(cfg)
}
func AfterCreate[T any](cfg *SrvConfig, pos CallbackPosition, f func(cfg *SrvConfig, newValues []T) error) error {
	if pos != CallbackAfter {
		return nil
	}
	return modifyCallback(cfg, func(cfg *SrvConfig, oldValues []T, newValues []T) error {
		return f(cfg, newValues)
	}, iface.MethodCreate)
}
func AfterCreateAny(cfg *SrvConfig, pos CallbackPosition, fc func(cfg *SrvConfig) error) error {
	if pos != CallbackAfter {
		return nil
	}
	if cfg.Method != iface.MethodCreate {
		return nil
	}
	return fc(cfg)
}
func BeforeUpdate[T any](cfg *SrvConfig, pos CallbackPosition, f func(cfg *SrvConfig, oldValues []T, newValues []T) error) error {
	if pos != CallbackBefore {
		return nil
	}
	return modifyCallback(cfg, f, iface.MethodUpdate, iface.MethodUpdateByFilters)
}
func BeforeUpdateAny(cfg *SrvConfig, pos CallbackPosition, fc func(cfg *SrvConfig) error) error {
	if pos != CallbackBefore {
		return nil
	}
	if cfg.Method != iface.MethodUpdate && cfg.Method != iface.MethodUpdateByFilters {
		return nil
	}
	return fc(cfg)
}
func AfterUpdate[T any](cfg *SrvConfig, pos CallbackPosition, f func(cfg *SrvConfig, oldValues []T, newValues []T) error) error {
	if pos != CallbackAfter {
		return nil
	}
	return modifyCallback(cfg, f, iface.MethodUpdate, iface.MethodUpdateByFilters)
}
func AfterUpdateAny(cfg *SrvConfig, pos CallbackPosition, fc func(cfg *SrvConfig) error) error {
	if pos != CallbackAfter {
		return nil
	}
	if cfg.Method != iface.MethodUpdate && cfg.Method != iface.MethodUpdateByFilters {
		return nil
	}
	return fc(cfg)
}
func BeforeDelete[T any](cfg *SrvConfig, pos CallbackPosition, f func(cfg *SrvConfig, oldValues []T) error) error {
	if pos != CallbackBefore {
		return nil
	}
	return modifyCallback(cfg, func(cfg *SrvConfig, oldValues []T, newValues []T) error {
		return f(cfg, oldValues)
	}, iface.MethodDelete, iface.MethodDeleteByFilters)
}
func BeforeDeleteAny(cfg *SrvConfig, pos CallbackPosition, fc func(cfg *SrvConfig) error) error {
	if pos != CallbackBefore {
		return nil
	}
	if cfg.Method != iface.MethodDelete && cfg.Method != iface.MethodDeleteByFilters {
		return nil
	}
	return fc(cfg)
}
func AfterDelete[T any](cfg *SrvConfig, pos CallbackPosition, f func(cfg *SrvConfig, oldValues []T) error) error {
	if pos != CallbackAfter {
		return nil
	}
	return modifyCallback(cfg, func(cfg *SrvConfig, oldValues []T, newValues []T) error {
		return f(cfg, oldValues)
	}, iface.MethodDelete, iface.MethodDeleteByFilters)
}
func AfterDeleteAny(cfg *SrvConfig, pos CallbackPosition, fc func(cfg *SrvConfig) error) error {
	if pos != CallbackAfter {
		return nil
	}
	if cfg.Method != iface.MethodDelete && cfg.Method != iface.MethodDeleteByFilters {
		return nil
	}
	return fc(cfg)
}
func BeforeModify[T any](cfg *SrvConfig, pos CallbackPosition, f func(cfg *SrvConfig, oldValues []T, newValues []T) error, methodList ...iface.Method) error {
	if pos != CallbackBefore {
		return nil
	}
	if cfg.CurdType != iface.CurdModify {
		return nil
	}
	if len(methodList) < 1 {
		methodList = append(methodList, iface.MethodCreate, iface.MethodUpdate, iface.MethodUpdateByFilters, iface.MethodDelete, iface.MethodDeleteByFilters)
	}
	return modifyCallback(cfg, f, methodList...)
}
func BeforeModifyAny(cfg *SrvConfig, pos CallbackPosition, fc func(cfg *SrvConfig) error, methodList ...iface.Method) error {
	if pos != CallbackBefore {
		return nil
	}
	if cfg.CurdType != iface.CurdModify {
		return nil
	}
	if len(methodList) > 0 && !utils.ContainAny(methodList, cfg.Method) {
		return nil
	}
	return fc(cfg)
}
func AfterModify[T any](cfg *SrvConfig, pos CallbackPosition, f func(cfg *SrvConfig, oldValues []T, newValues []T) error, methodList ...iface.Method) error {
	if pos != CallbackAfter {
		return nil
	}
	if cfg.CurdType != iface.CurdModify {
		return nil
	}
	if len(methodList) < 1 {
		methodList = append(methodList, iface.MethodCreate, iface.MethodUpdate, iface.MethodUpdateByFilters, iface.MethodDelete, iface.MethodDeleteByFilters)
	}
	return modifyCallback(cfg, f, methodList...)
}
func AfterModifyAny(cfg *SrvConfig, pos CallbackPosition, fc func(cfg *SrvConfig) error, methodList ...iface.Method) error {
	if pos != CallbackAfter {
		return nil
	}
	if cfg.CurdType != iface.CurdModify {
		return nil
	}
	if len(methodList) > 0 && !utils.ContainAny(methodList, cfg.Method) {
		return nil
	}
	return fc(cfg)
}
func modifyCallback[T any](cfg *SrvConfig, f func(cfg *SrvConfig, oldValues []T, newValues []T) error, methodList ...iface.Method) error {

	if len(methodList) > 0 {
		if !utils.ContainAny(methodList, cfg.Method) {
			return nil
		}
	}

	var (
		t         T
		oldValues []T
		newValues []T
	)
	if cfg.OldValues == nil && cfg.NewValues == nil {
		return f(cfg, oldValues, newValues)
	}
	if _, ok := any(t).(reflect.Value); ok {
		if cfg.OldValues != nil {
			vs := transValues(reflect.ValueOf(cfg.OldValues))
			utils.ConvertToSlice(vs, &oldValues)
		}
		if cfg.NewValues != nil {
			vs := transValues(reflect.ValueOf(cfg.NewValues))
			utils.ConvertToSlice(vs, &newValues)
		}

	} else {
		if cfg.OldValues != nil {
			utils.ConvertToSlice(cfg.OldValues, &oldValues)
		}
		if cfg.NewValues != nil {
			utils.ConvertToSlice(cfg.NewValues, &newValues)
		}
	}

	return f(cfg, oldValues, newValues)
}

func BeforeRetrieve(cfg *SrvConfig, pos CallbackPosition, f func(cfg *SrvConfig) error, curds ...iface.Method) error {
	if pos != CallbackBefore {
		return nil
	}

	if len(curds) > 0 {
		if !utils.ContainAny(curds, cfg.Method) {
			return nil
		}
	}
	return f(cfg)
}
func AfterRetrieve[T any](cfg *SrvConfig, pos CallbackPosition, f func(cfg *SrvConfig, rs []T) error, curds ...iface.Method) error {
	if pos != CallbackAfter {
		return nil
	}
	return afterRetrieveCallback(cfg, f, curds...)
}

// AfterRetrieveAny
// 当查询不太确定返回何种结果的时候可以用这个，返回的data可能是[]ctype.Map或者指定的结构体[]xxx
func AfterRetrieveAny(cfg *SrvConfig, pos CallbackPosition, fc func(cfg *SrvConfig) error, methodList ...iface.Method) error {
	if pos != CallbackAfter {
		return nil
	}
	if cfg.CurdType != iface.CurdRetrieve {
		return nil
	}
	if len(methodList) > 0 && !utils.ContainAny(methodList, cfg.Method) {
		return nil
	}

	return fc(cfg)
}

func afterRetrieveCallback[T any](cfg *SrvConfig, f func(cfg *SrvConfig, rs []T) error, curds ...iface.Method) error {

	if len(curds) > 0 {
		if !utils.ContainAny(curds, cfg.Method) {
			return nil
		}
	}

	var (
		t    T
		rsDt []T
	)
	if cfg.Result.Data == nil {
		return f(cfg, rsDt)
	}
	if _, ok := any(t).(reflect.Value); ok {
		values := transValues(reflect.ValueOf(cfg.Result.Data))
		utils.ConvertToSlice(values, &rsDt)
	} else {
		utils.ConvertToSlice(cfg.Result.Data, &rsDt)
	}

	return f(cfg, rsDt)
}

func transValues(val reflect.Value) []reflect.Value {

	var (
		values = make([]reflect.Value, 0, 10)
	)

	if !val.IsValid() {
		return values
	}
	switch val.Kind() {
	case reflect.Pointer:
		values = append(values, transValues(val.Elem())...)
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			values = append(values, val.Index(i))
		}
	case reflect.Struct:
		values = append(values, val)
	default:
		values = append(values, val)
	}

	return values
}
