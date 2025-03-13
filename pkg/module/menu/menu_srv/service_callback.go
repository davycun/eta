package menu_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/menu"
)

func modifyCallback(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []menu.Menu) error {
				menu.DataCache.SetHasAll(cfg.TxDB, false)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterUpdate(cfg, pos, func(cfg *hook.SrvConfig, oldValues []menu.Menu, newValues []menu.Menu) error {
				menu.DelCache(cfg.TxDB, oldValues...)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []menu.Menu) error {
				menu.DelCache(cfg.TxDB, oldValues...)
				return nil
			})
		}).Err
}

func retrieveCallback(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	return hook.AfterRetrieveAny(cfg, pos, func(config *hook.SrvConfig) error {

		switch listRs := cfg.Result.Data.(type) {
		case []menu.Menu:
			return fill(cfg, listRs)
		case []ctype.Map:
			return nil
		}
		return nil
	})
}
