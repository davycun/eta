package app_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/history"
	"github.com/davycun/eta/pkg/core/migrate/mig_hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/app"
)

var (
	defaultApp = app.App{
		Name:      "默认",
		Valid:     ctype.Boolean{Valid: true, Data: true},
		IsDefault: ctype.NewBooleanPrt(true),
	}
)

func afterMigrate(cfg *mig_hook.MigConfig, pos mig_hook.CallbackPosition) error {
	var (
		db      = cfg.TxDB
		c       = cfg.C
		ids     []string
		err     error
		appList = []app.App{defaultApp}
	)
	if pos != mig_hook.CallbackAfter {
		return err
	}

	if db == nil {
		db = global.GetLocalGorm()
	}
	if c == nil {
		c = ctx.NewContext()
	}

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			scm := dorm.GetDbSchema(db)
			return history.CreateTrigger(db, scm, constants.TableApp)
		}).
		Call(func(cl *caller.Caller) error {
			//如果已经存在APP，就不创建默认app，Stop会组织后续的Call继续调用
			err = dorm.Table(db, constants.TableApp).Select(entity.IdDbName).Limit(1).Find(&ids).Error
			if len(ids) > 0 {
				cl.Stop()
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			return beforeCreateApp(c, db, appList)
		}).
		Call(func(cl *caller.Caller) error {
			return dorm.Table(db, constants.TableApp).Create(&appList).Error
		}).
		Call(func(cl *caller.Caller) error {
			return afterCreateApp(c, db, appList)
		}).Err

	return err
}
