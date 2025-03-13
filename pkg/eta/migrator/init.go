package migrator

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/core/migrate"
	"github.com/davycun/eta/pkg/eta/migrator/mig_type"
	"gorm.io/gorm"
)

func init() {
	AddCallback(MigAll, createSchema)
	AddCallback(MigAll, initDb)
}

// 提前创建schema
func createSchema(cfg *MigConfig, pos CallbackPosition) error {
	if pos != BeforeCallback {
		return nil
	}
	if cfg.TxDB == nil {
		return nil
	}
	var (
		scm = dorm.GetDbSchema(cfg.TxDB)
		mg  = migrate.NewMigrator(cfg.TxDB, cfg.C)
	)
	if !mg.Schema().SchemaExists(scm) {
		return mg.Schema().CreateSchema(scm)
	}
	return nil
}

func initDb(cfg *MigConfig, pos CallbackPosition) error {
	if pos != BeforeCallback {
		return nil
	}
	return mig_type.MigrateTypeAndFunction(cfg.TxDB)
}

// 达梦的表空间,这个回调不启用，在这里只是一个示例，可以供有需要的用户使用
func createDmTableSpace(cfg *MigConfig, pos CallbackPosition) error {
	if dorm.GetDbType(cfg.TxDB) != dorm.DaMeng || pos != BeforeCallback {
		return nil
	}
	var (
		appScm   = dorm.GetDbSchema(cfg.TxDB)
		localScm = dorm.GetDbSchema(global.GetLocalGorm())
		pattern  = `CREATE TABLESPACE IF NOT EXISTS "%s_%s" DATAFILE '/opt/dmdbms/data_ts/%s/%s/%s/data.dbf' SIZE 128 AUTOEXTEND ON NEXT 128 MAXSIZE UNLIMITED`
		ts       = []string{
			fmt.Sprintf(pattern, appScm, "tableName", localScm, appScm, "tableName"),
		}
	)

	for _, v := range ts {
		err := cfg.TxDB.Transaction(func(tx *gorm.DB) error {
			return tx.Exec(v).Error
		})
		if err != nil {
			return err
		}
	}
	return nil
}
