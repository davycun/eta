package setting_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/ws"
	"github.com/davycun/eta/pkg/eta/constants"
	setting2 "github.com/davycun/eta/pkg/module/setting"
	"github.com/duke-git/lancet/v2/slice"
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

func init() {
	hook.AddModifyCallback(constants.TableSetting, modifyCallback)
}

func modifyCallback(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.BeforeUpdate(cfg, pos, func(cfg *hook.SrvConfig, oldValues []reflect.Value, newValues []reflect.Value) error {
				return removeUpdateField(cfg)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []reflect.Value) error {
				setting2.HasCacheAll(cfg.TxDB, false)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterUpdate(cfg, pos, func(cfg *hook.SrvConfig, oldValues []reflect.Value, newValues []reflect.Value) error {
				setting2.DelCache(cfg.TxDB, oldValues...)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []reflect.Value) error {
				setting2.DelCache(cfg.TxDB, oldValues...)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterUpdate(cfg, pos, notifyConfigChanged)
		}).Err
}

// removeUpdateField 移除不可更新字段
func removeUpdateField(cfg *hook.SrvConfig) error {
	fields := []string{"namespace", "category", "name"}
	entityFieldNames := make([]string, len(fields))

	cfg.Param.Columns = slice.Filter(cfg.Param.Columns, func(index int, item string) bool {
		return !slice.Contain(fields, item)
	})

	for _, v := range cfg.Values {
		for _, f := range entityFieldNames {
			var fld reflect.Value
			if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
				fld = v.Elem().FieldByName(f)
			} else {
				fld = v.FieldByName(f)
			}
			if fld.IsValid() && !fld.IsZero() {
				fld.SetZero()
			}
		}
	}
	return nil
}

func notifyConfigChanged(cfg *hook.SrvConfig, oldValues []reflect.Value, newValues []reflect.Value) error {
	notifyBody := slice.Map(newValues, func(_ int, v reflect.Value) setting2.Setting {
		return setting2.Setting{
			Namespace: entity.GetString(v, "namespace"),
			Category:  entity.GetString(v, "category"),
			Name:      entity.GetString(v, "name"),
		}
	})

	msg, err := jsoniter.MarshalToString(notifyBody)
	if err != nil {
		return err
	}
	ws.SendMessage(constants.WsKeySettingChanged, msg)
	return nil
}
