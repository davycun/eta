package user_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/history"
	"github.com/davycun/eta/pkg/core/migrate/mig_hook"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/user"
)

// 在系统初始化的时候，需要默认创建root用户
func afterMigratorUser(mc *mig_hook.MigConfig, pos mig_hook.CallbackPosition) error {

	if pos != mig_hook.CallbackAfter {
		return nil
	}
	var (
		usList = []user.User{user.GetRootUser()}
		db     = mc.TxDB
	)

	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			scm := dorm.GetDbSchema(db)
			return history.CreateTrigger(db, scm, constants.TableUser)
		}).
		Call(func(cl *caller.Caller) error {
			err := beforeCreateFillUserField(mc.C, mc.TxDB, usList)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			var ct int64
			err := dorm.Table(mc.TxDB, constants.TableUser).
				Where(&user.User{Account: ctype.NewStringPrt(user.RootUserAccount)}).
				Count(&ct).Error
			if err != nil {
				return err
			}
			//如果存在就不创建了
			if ct > 0 {
				return nil
			}
			return service.NewSrvWrapper(constants.TableUser, mc.C, mc.TxDB).SetData([]user.User{user.GetRootUser()}).Create()
		}).Err

	return err
}
