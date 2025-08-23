package user2app

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/history"
	"github.com/davycun/eta/pkg/core/migrate"
	"github.com/davycun/eta/pkg/eta/constants"
)

func afterMigratorUser2App(mc *migrate.MigConfig, pos migrate.CallbackPosition) error {

	if pos != migrate.CallbackAfter {
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
