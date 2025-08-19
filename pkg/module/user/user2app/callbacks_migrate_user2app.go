package user2app

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/history"
	"github.com/davycun/eta/pkg/core/migrate/mig_hook"
	"github.com/davycun/eta/pkg/eta/constants"
)

func afterMigratorUser2App(mc *mig_hook.MigConfig, pos mig_hook.CallbackPosition) error {

	if pos != mig_hook.CallbackAfter {
		return nil
	}
	var (
		err error
		db  = mc.TxDB
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			scm := dorm.GetDbSchema(db)
			return history.CreateTrigger(db, scm, constants.TableUser2App)
		}).
		Call(func(cl *caller.Caller) error {
			return dorm.CreateUniqueIndex(db, constants.TableUser2App, entity.FromIdDbName, entity.ToIdDbName)
		}).Err
	return err
}
