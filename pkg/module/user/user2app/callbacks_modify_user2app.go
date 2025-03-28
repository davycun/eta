package user2app

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/core/service/hook"
)

func modifyCallbackUser2App(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.AfterModify(cfg, pos, func(cfg *hook.SrvConfig, oldValues []User2App, newValues []User2App) error {
				CleanCache(cfg.TxDB)
				return nil
			})
		}).Err
	return err
}
