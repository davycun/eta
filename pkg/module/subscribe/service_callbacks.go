package subscribe

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/core/service/hook"
)

// 修改自己的时候需要做的回调
func selfModifyCallback(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []Subscriber) error {
				allData.SetHasAll(cfg.TxDB, false)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterUpdate(cfg, pos, func(cfg *hook.SrvConfig, oldValues []Subscriber, newValues []Subscriber) error {
				DelCache(cfg.TxDB, oldValues...)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []Subscriber) error {
				DelCache(cfg.TxDB, oldValues...)
				return nil
			})
		}).Err
	return err
}
