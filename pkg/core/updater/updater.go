package updater

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/mysql"
	"github.com/davycun/eta/pkg/core/entity"
	"gorm.io/gorm"
)

var (
	//第二行是pg的，第一行是达梦的，不支持updater功能的类型（主要是还数组和自定义类型），没有枚举完，下划线代表数组
	notSupportType = []string{
		"st_geometry", "arr_int_cls", "arr_int", "arr_str_cls", "arr_str",
		"geometry", "array", "jsonb", "_int", "_varchar", "_int4", "_int8", "_text", "clob", "text",
	}
)

func CreateUpdaterTrigger(db *gorm.DB, tableName string) error {
	var (
		scm    = dorm.GetDbSchema(db)
		dbType = dorm.GetDbType(db)
	)
	switch dbType {
	case dorm.DaMeng:
		return caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return createDmUpdaterTrigger(db, tableName)
			}).
			Call(func(cl *caller.Caller) error {
				idx := fmt.Sprintf(`CREATE ARRAY INDEX IF NOT EXISTS "idx_%s_%s" ON "%s"."%s"("%s")`, tableName, entity.FieldUpdaterIdsDbName, scm, tableName, entity.FieldUpdaterIdsDbName)
				return db.Exec(idx).Error
			}).Err
	case dorm.PostgreSQL:
		return caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return createPostgresUpdaterTrigger(db, scm, tableName)
			}).
			Call(func(cl *caller.Caller) error {
				idx := fmt.Sprintf(`CREATE INDEX IF NOT EXISTS "idx_%s_%s" ON "%s"."%s" USING GIN ("%s")`, tableName, entity.FieldUpdaterIdsDbName, scm, tableName, entity.FieldUpdaterIdsDbName)
				return db.Exec(idx).Error
			}).Err
	case dorm.Mysql:
		return caller.NewCaller().
			// 创建触发器
			Call(func(cl *caller.Caller) error {
				return createMysqlUpdaterTrigger(db, scm, tableName)
			}).
			// 创建 entity.FieldUpdaterIdsDbName 字段索引
			Call(func(cl *caller.Caller) error {
				// CREATE INDEX `idx_d_ut_feature_1717589249444_field_updater_ids` ON `eta_default`.`d_ut_feature_1717589249444` ((cast(`field_updater_ids`->'$' as char(255) ARRAY)))
				// SELECT * FROM `eta_default`.`d_ut_feature_1717589249444` where JSON_CONTAINS(`field_updater_ids`, '"188590326540742657"');
				indexName := fmt.Sprintf("idx_%s_%s", tableName, entity.FieldUpdaterIdsDbName)
				idxSql := fmt.Sprintf("CREATE INDEX `%s` ON %s ((cast(%s->'$' as char(255) array)))",
					indexName, dorm.GetScmTableName(db, tableName), dorm.Quote(dbType, entity.FieldUpdaterIdsDbName))
				return mysql.CreateIndexIfNotExists(db, tableName, indexName, idxSql)
			}).
			Err
	}
	return nil
}

func DropUpdaterTrigger(db *gorm.DB, tableName string) error {
	var (
		err         error
		dbType      = dorm.GetDbType(db)
		scm         = dorm.GetDbSchema(db)
		scmFuncName = fmt.Sprintf(`"%s"."func_%s_updater"`, scm, tableName)
		triggerName = fmt.Sprintf(`"trigger_%s_updater"`, tableName)
	)
	switch dbType {
	case dorm.DaMeng:
		err = db.Exec(fmt.Sprintf(`drop trigger if exists "%s".%s`, scm, triggerName)).Error
	case dorm.PostgreSQL:
		err = caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return db.Exec(fmt.Sprintf(`drop function if exists %s() cascade`, scmFuncName)).Error
			}).
			Call(func(cl *caller.Caller) error {
				return db.Exec(fmt.Sprintf(`drop trigger if exists %s on "%s"."%s" cascade`, triggerName, scm, tableName)).Error
			}).Err
	case dorm.Mysql:
		return dropMysqlUpdaterTrigger(db, scm, tableName)
	}
	return err
}
