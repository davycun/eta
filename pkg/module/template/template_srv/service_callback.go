package template_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/template"
)

func modifyCallback(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, beforeCreateFillField)
		}).
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []template.Template) error {
				return SignValidator(newValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []template.Template) error {
				return EncryptValidator(newValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []template.Template) error {
				return RaDbFieldsValidator(newValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.BeforeUpdate(cfg, pos, beforeUpdateValidate)
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, afterCreate)
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterDelete(cfg, pos, afterDelete)
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterUpdate(cfg, pos, afterUpdate)
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterModify(cfg, pos, func(cfg *hook.SrvConfig, oldValues []template.Template, newValues []template.Template) error {
				template.DelCache(cfg.TxDB, oldValues...)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterModify(cfg, pos, func(cfg *hook.SrvConfig, oldValues []template.Template, newValues []template.Template) error {
				template.DelCache(cfg.TxDB, oldValues...)
				if cfg.Method == iface.MethodCreate {
					template.SetCacheHasAll(cfg.TxDB, false)
				}
				return nil
			})
		}).Err
	return err
}
