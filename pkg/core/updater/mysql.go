package updater

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/db_table"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"gorm.io/gorm"
	"strings"
)

func createMysqlUpdaterTrigger(db *gorm.DB, scm, tableName string) error {
	var (
		scmTbName    = fmt.Sprintf("`%s`.`%s`", scm, tableName)
		tgInsertName = fmt.Sprintf("`%s`.`trigger_%s_insert`", scm, tableName)
		tgUpdateName = fmt.Sprintf("`%s`.`trigger_%s_updater`", scm, tableName)
		cols         []db_table.Column
		dbType       = dorm.GetDbType(db)
		err          = db_table.FetchColumns(db, tableName, &cols)
	)
	if err != nil {
		return err
	}
	dropTgInsert := `DROP TRIGGER IF EXISTS ` + tgInsertName
	dropTgUpdate := `DROP TRIGGER IF EXISTS ` + tgUpdateName
	tgInsert := mysqlTgCreate(scmTbName, tgInsertName, cols, dbType)
	tgUpdate := mysqlTgUpdate(scmTbName, tgUpdateName, cols, dbType)

	cl := func(tx *gorm.DB) error {
		return caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return tx.Exec(dropTgInsert).Error
			}).
			Call(func(cl *caller.Caller) error {
				return tx.Exec(tgInsert).Error
			}).
			Call(func(cl *caller.Caller) error {
				return tx.Exec(dropTgUpdate).Error
			}).
			Call(func(cl *caller.Caller) error {
				return tx.Exec(tgUpdate).Error
			}).Err
	}

	if dorm.InTransaction(db) {
		return cl(db)
	} else {
		return db.Transaction(func(tx *gorm.DB) error {
			return cl(tx)
		})
	}
}

func mysqlTgCreate(scmTbName, tgName string, cols []db_table.Column, dbType dorm.DbType) string {
	var (
		bd                    = strings.Builder{}
		updaterIdDbName       = dorm.Quote(dbType, entity.UpdaterIdDbName)
		fieldUpdaterDbName    = dorm.Quote(dbType, entity.FieldUpdaterDbName)
		fieldUpdaterIdsDbName = dorm.Quote(dbType, entity.FieldUpdaterIdsDbName)
	)

	for _, v := range cols {
		bd.WriteString(fmt.Sprintf("SET updater = JSON_SET(updater, '$.%s', NEW.%s);", v.ColName, updaterIdDbName))
		bd.WriteString("\n")
	}

	sql := `
	CREATE TRIGGER ` + tgName + ` BEFORE INSERT ON ` + scmTbName + `
	FOR EACH ROW
	BEGIN
		DECLARE updater json;
        IF NEW.` + fieldUpdaterDbName + ` is null THEN
			SET updater = JSON_OBJECT();
			` + bd.String() + `
			SET NEW.` + fieldUpdaterDbName + ` = updater;
  		END IF;

		IF NEW.` + fieldUpdaterIdsDbName + ` is null and NEW.` + updaterIdDbName + ` is not null THEN
			SET NEW.` + fieldUpdaterIdsDbName + ` = JSON_ARRAY(NEW.` + updaterIdDbName + `);
		END IF;
	END`
	return sql
}

func mysqlTgUpdate(scmTbName, tgName string, cols []db_table.Column, dbType dorm.DbType) string {
	var (
		bdF                   = strings.Builder{}
		updaterIdDbName       = dorm.Quote(dbType, entity.UpdaterIdDbName)
		fieldUpdaterDbName    = dorm.Quote(dbType, entity.FieldUpdaterDbName)
		fieldUpdaterIdsDbName = dorm.Quote(dbType, entity.FieldUpdaterIdsDbName)
	)

	for _, v := range cols {
		//更新的时候由于自定义类型不能用比较所以去掉不支持
		if utils.ContainAny(entity.DefaultVertexColumns, v.ColName) || utils.ContainAny(notSupportType, strings.ToLower(v.ColType)) {
			continue
		}
		dbColName := dorm.Quote(dbType, v.ColName)

		bdF.WriteString(`IF OLD.` + dbColName + ` != NEW.` + dbColName + ` THEN
			SET flag = true;
			SET updater = JSON_SET(updater, '$.` + v.ColName + `', NEW.` + updaterIdDbName + `);
		END IF;`)
		bdF.WriteString("\n")
	}

	sql := `
	CREATE TRIGGER ` + tgName + ` BEFORE UPDATE ON ` + scmTbName + `
	FOR EACH ROW
	BEGIN
		DECLARE updater json;
		DECLARE updater_ids json;
    	DECLARE flag bit;

        SET flag = 0;
		SET updater = JSON_OBJECT();
		SET updater_ids = JSON_ARRAY();
		IF OLD.` + fieldUpdaterDbName + ` is not null THEN
			SET updater = OLD.` + fieldUpdaterDbName + `;
		END IF;
		IF OLD.` + fieldUpdaterIdsDbName + ` is not null THEN
			SET updater_ids = OLD.` + fieldUpdaterIdsDbName + `;
		END IF;

		` + bdF.String() + `

		IF flag THEN 
			SET NEW.` + fieldUpdaterDbName + ` = updater;

			IF not JSON_CONTAINS(updater_ids, concat('"',NEW.` + updaterIdDbName + `,'"')) THEN
				SET NEW.` + fieldUpdaterIdsDbName + ` = JSON_ARRAY_APPEND(updater_ids, '$', NEW.` + updaterIdDbName + `);
			END IF;

		END IF;
	END`
	return sql
}

func dropMysqlUpdaterTrigger(db *gorm.DB, scm, tableName string) error {
	var (
		dropTgInsert = fmt.Sprintf("DROP TRIGGER IF EXISTS `%s`.`trigger_%s_insert`", scm, tableName)
		dropTgUpdate = fmt.Sprintf("DROP TRIGGER IF EXISTS `%s`.`trigger_%s_updater`", scm, tableName)
	)

	cl := func(tx *gorm.DB) error {
		return caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return tx.Exec(dropTgInsert).Error
			}).
			Call(func(cl *caller.Caller) error {
				return tx.Exec(dropTgUpdate).Error
			}).Err
	}

	if dorm.InTransaction(db) {
		return cl(db)
	} else {
		return db.Transaction(func(tx *gorm.DB) error {
			return cl(tx)
		})
	}

}
