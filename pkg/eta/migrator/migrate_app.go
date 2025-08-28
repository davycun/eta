package migrator

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/migrate"
	"github.com/davycun/eta/pkg/eta/ecf"
	"gorm.io/gorm"
)

func getMigrateTable(db *gorm.DB) []entity.Table {
	return ecf.GetMigrateAppTable(db, global.GetConfig().Server.MigratePkg...)
}

// MigrateApp
// dbs 支持不传或者db、doris、es中的任意个数
func MigrateApp(db *gorm.DB, c *ctx.Context, param *MigrateAppParam) error {
	if param == nil {
		param = &MigrateAppParam{}
	}
	var (
		err  error
		dbs  = param.Dbs
		txDb = db
		mg   = migrate.NewMigrator(db, c) //这里注意达梦进行ddl的时候默认会把前面的事务进行提交，所以不用事务db
		pm   = &dto.Param{RetrieveParam: dto.RetrieveParam{Extra: param}}
		mc   = NewCallbackCaller(c, db, pm)
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return mc.BeforeMigrateApp()
		}).
		Call(func(cl *caller.Caller) error {
			if len(dbs) < 1 || utils.ContainAny(dbs, "db") {
				return mg.MigrateOption(getMigrateTable(db)...)
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			if len(dbs) < 1 || utils.ContainAny(dbs, "es") {
				return MigrateElasticsearch(txDb, param.Es)
			}
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			return mc.AfterMigrateApp()
		}).Err

	return err
}
