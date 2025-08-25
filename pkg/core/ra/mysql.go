package ra

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/db_table"
	"github.com/davycun/eta/pkg/common/dorm/xa"
	"github.com/davycun/eta/pkg/core/entity"
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
)

func createMysqlTrigger(db *gorm.DB, scm, tableName string, raFields []string) error {
	var (
		scmTbName    = fmt.Sprintf("`%s`.`%s`", scm, tableName)
		tgInsertName = fmt.Sprintf("trigger_%s_insert_ra", tableName)
		tgUpdateName = fmt.Sprintf("trigger_%s_updater_ra", tableName)
		cols         []db_table.Column
		dbType       = dorm.GetDbType(db)
		err          = db_table.FetchColumns(db, tableName, &cols)
	)
	if err != nil {
		return err
	}
	tgInsert := mysqlTgCreate(scmTbName, dorm.Quote(dbType, scm, tgInsertName), dbType, raFields)
	tgUpdate := mysqlTgUpdate(scmTbName, dorm.Quote(dbType, scm, tgUpdateName), dbType, raFields)

	cl := func(tx *gorm.DB) error {
		return caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return dorm.TriggerDelete(tx, scm+"."+tgInsertName, "")
			}).
			Call(func(cl *caller.Caller) error {
				return tx.Exec(tgInsert).Error
			}).
			Call(func(cl *caller.Caller) error {
				return dorm.TriggerDelete(tx, scm+"."+tgUpdateName, "")
			}).
			Call(func(cl *caller.Caller) error {
				return tx.Exec(tgUpdate).Error
			}).Err
	}

	if xa.InTransaction(db) {
		return cl(db)
	} else {
		return db.Transaction(func(tx *gorm.DB) error {
			return cl(tx)
		})
	}
}

func mysqlTgCreate(scmTbName, tgName string, dbType dorm.DbType, raFields []string) string {
	var (
		raContentDbName = dorm.Quote(dbType, entity.RaContentDbName)
		raString        = strings.Join(slice.Map(raFields, func(_ int, v string) string {
			return fmt.Sprintf("NEW.%s", dorm.Quote(dbType, v))
		}), `||' '||`)
	)

	sql := `
	CREATE TRIGGER ` + tgName + ` BEFORE INSERT ON ` + scmTbName + `
	FOR EACH ROW
	BEGIN
		IF NEW.` + raContentDbName + ` is null THEN
			SET NEW.` + raContentDbName + ` = ` + raString + `;
		END IF;
	END`
	return sql
}

func mysqlTgUpdate(scmTbName, tgName string, dbType dorm.DbType, raFields []string) string {
	var (
		raContentDbName = dorm.Quote(dbType, entity.RaContentDbName)
		raString        = strings.Join(slice.Map(raFields, func(_ int, v string) string {
			return fmt.Sprintf("NEW.%s", dorm.Quote(dbType, v))
		}), `||' '||`)
	)

	sql := `
	CREATE TRIGGER ` + tgName + ` BEFORE UPDATE ON ` + scmTbName + `
	FOR EACH ROW
	BEGIN
		SET NEW.` + raContentDbName + ` = ` + raString + `;
	END`
	return sql
}

func dropMysqlTrigger(db *gorm.DB, scm, tableName string) error {
	var (
		tgInsert = fmt.Sprintf("%s.trigger_%s_insert_ra", scm, tableName)
		tgUpdate = fmt.Sprintf("%s.trigger_%s_updater_ra", scm, tableName)
	)

	cl := func(tx *gorm.DB) error {
		return caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return dorm.TriggerDelete(tx, tgInsert, "")
			}).
			Call(func(cl *caller.Caller) error {
				return dorm.TriggerDelete(tx, tgUpdate, "")
			}).Err
	}

	if xa.InTransaction(db) {
		return cl(db)
	} else {
		return db.Transaction(func(tx *gorm.DB) error {
			return cl(tx)
		})
	}

}
