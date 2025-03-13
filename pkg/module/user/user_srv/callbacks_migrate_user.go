package user_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/history"
	"github.com/davycun/eta/pkg/core/migrate/mig_hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/user"
	"gorm.io/gorm/clause"
)

// 1.在系统初始化的时候，创建了用户表之后，需要默认创建root用户
// 2.注意这个root用户的创建不走service，直接走数据库操作
// 3.创建用户的时候需要创建root与默认app的关系
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
			clf := clause.OnConflict{
				Columns:   []clause.Column{{Name: "account"}},
				DoNothing: true,
			}
			return dorm.Table(mc.TxDB, constants.TableUser).Clauses(clf).Create(&usList).Error
		}).Err

	return err
}
