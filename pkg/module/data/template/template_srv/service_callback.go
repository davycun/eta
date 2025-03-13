package template_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/data/template"
)

func modifyCallback(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, beforeCreateFillField)
		}).
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []template.Template) error {
				return template.SignValidator(newValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []template.Template) error {
				return template.EncryptValidator(newValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []template.Template) error {
				return template.RaDbFieldsValidator(newValues)
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
		}).Err
	return err
}
