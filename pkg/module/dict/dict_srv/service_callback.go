package dict_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/ws"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/dict"
	"github.com/duke-git/lancet/v2/slice"
	jsoniter "github.com/json-iterator/go"
)

func modifyCallback(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []dict.Dictionary) error {
				dict.DataCache.SetHasAll(cfg.TxDB, false)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterUpdate(cfg, pos, func(cfg *hook.SrvConfig, oldValues []dict.Dictionary, newValues []dict.Dictionary) error {
				dict.DelCache(cfg.TxDB, oldValues...)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []dict.Dictionary) error {
				dict.DelCache(cfg.TxDB, oldValues...)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterUpdate(cfg, pos, notifyDictChanged)
		}).Err

	return err
}

func retrieveCallbacks(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	return hook.AfterRetrieveAny(cfg, pos, func(cfg *hook.SrvConfig) error {
		return processResult(cfg, cfg.Result.Data)
	})
}

func notifyDictChanged(cfg *hook.SrvConfig, oldValues []dict.Dictionary, newValues []dict.Dictionary) error {
	notifyBody := slice.Map(newValues, func(_ int, v dict.Dictionary) dict.Dictionary {
		return dict.Dictionary{
			Namespace: v.Namespace,
			Category:  v.Category,
			Name:      v.Name,
		}
	})

	msg, err := jsoniter.MarshalToString(notifyBody)
	if err != nil {
		return err
	}
	ws.SendMessage(constants.WsKeyDictChanged, msg)
	return nil
}
