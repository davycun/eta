package user_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/history"
	"github.com/davycun/eta/pkg/core/migrate/mig_hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/user"
)

// 1.在系统初始化的时候，创建了用户表之后，需要默认创建root用户
// 2.注意这个root用户的创建不走service，直接走数据库操作
// 3.创建用户的时候需要创建root与默认app的关系
func afterMigratorUser2App(mc *mig_hook.MigConfig, pos mig_hook.CallbackPosition) error {

	if pos != mig_hook.CallbackAfter {
		return nil
	}
	var (
		err error
		db  = mc.TxDB
		us  user.User
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			us, err = user.LoadDefaultUser(db)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			scm := dorm.GetDbSchema(db)
			return history.CreateTrigger(db, scm, constants.TableUser2App)
		}).
		Call(func(cl *caller.Caller) error {
			scm := dorm.GetDbSchema(db)
			return history.CreateTrigger(db, scm, constants.TableUser2App)
		}).
		Call(func(cl *caller.Caller) error {
			return dorm.CreateUniqueIndex(db, constants.TableUser2App, entity.FromIdDbName, entity.ToIdDbName)
		}).
		Call(func(cl *caller.Caller) error {
			//放在这里目的是，为了先创建history的trigger
			return afterCreateUserNewUser2App(mc.C, mc.TxDB, []user.User{us})
		}).Err
	return err
}
