package migrator

import (
	"errors"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/migrate"
	"github.com/davycun/eta/pkg/core/service/ecf"
	"gorm.io/gorm"
	"slices"
)

func MigrateLocal(db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}
	c := ctx.NewContext()
	c.SetContextGorm(db)
	c.SetContextIsManager(true)
	var (
		mig = migrate.NewMigrator(db, c)
		mc  = NewMigConfig(ctx.NewContext(), db, &dto.Param{})
	)

	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return mc.BeforeMigrateLocal()
		}).
		Call(func(cl *caller.Caller) error {
			//注意下面的需要先migrate User，因为在创建User2App的时候会用到
			tbList := ecf.GetMigrateLocalEntityConfig()
			slices.SortFunc(tbList, func(a, b entity.Table) int {
				return a.Order - b.Order
			})
			return mig.MigrateOption(tbList...)
		}).
		Call(func(cl *caller.Caller) error {
			return mc.AfterMigrateLocal()
		}).Err
	return err
}
