package user_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/user"
	"github.com/davycun/eta/pkg/module/user/user2app"
	"gorm.io/gorm"
)

func modifyCallbackUser2App(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.AfterModify(cfg, pos, func(cfg *hook.SrvConfig, oldValues []user.User, newValues []user.User) error {
				user2app.CleanCache(cfg.TxDB)
				return nil
			})
		}).Err
	return err
}

func cleanUser2AppCache(txDb *gorm.DB, u2aList []user2app.User2App) {
	user2app.CleanCache(txDb)
}
