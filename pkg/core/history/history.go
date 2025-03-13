package history

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"strings"
)

func CreateTrigger(orm *gorm.DB, scm, tableName string) error {

	switch dorm.GetDbType(orm) {
	case dorm.DaMeng:
		return createDmTrigger(orm, scm, tableName)
	case dorm.PostgreSQL:
		return createPostgresTrigger(orm, scm, tableName)
	case dorm.Mysql:
		return createMysqlTrigger(orm, scm, tableName)
	}
	return nil
}

func DropTrigger(db *gorm.DB, tableName string) error {
	var (
		err         error
		scm         = dorm.GetDbSchema(db)
		scmFuncName = fmt.Sprintf(`"%s"."func_%s%s"`, scm, tableName, constants.TableHistorySubFix)
		triggerName = fmt.Sprintf(`"trigger_%s%s"`, tableName, constants.TableHistorySubFix)
	)
	switch db.Dialector.Name() {
	case dorm.DaMeng.String():
		err = db.Exec(fmt.Sprintf(`drop trigger if exists "%s".%s`, scm, triggerName)).Error
	case dorm.PostgreSQL.String():
		err = caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return db.Exec(fmt.Sprintf(`drop function if exists %s() cascade`, scmFuncName)).Error
			}).
			Call(func(cl *caller.Caller) error {
				return db.Exec(fmt.Sprintf(`drop trigger if exists %s on "%s"."%s"`, triggerName, scm, tableName)).Error
			}).Err
	}
	if err != nil {
		return err
	}
	return nil
}

func createPostgresTrigger(db *gorm.DB, scm, tableName string) error {
	var (
		scmFuncName     = fmt.Sprintf(`"%s"."func_%s%s"`, scm, tableName, constants.TableHistorySubFix)
		triggerName     = fmt.Sprintf(`"trigger_%s%s"`, tableName, constants.TableHistorySubFix)
		scmHisTbName    = fmt.Sprintf(`"%s"."%s%s"`, scm, tableName, constants.TableHistorySubFix)
		scmOriginTbName = fmt.Sprintf(`"%s"."%s"`, scm, tableName)
	)

	col := dorm.FetchTableColumns(db, scm, tableName)
	colBd := strings.Builder{}
	newColBd := strings.Builder{}
	oldColBd := strings.Builder{}

	for i, v := range col {
		if i > 0 {
			colBd.WriteString(",")
			newColBd.WriteString(",")
			oldColBd.WriteString(",")
		}
		colBd.WriteString(`"h_` + v + `"`)
		newColBd.WriteString(`NEW."` + v + `"`)
		oldColBd.WriteString(`OLD."` + v + `"`)
	}
	colStr := colBd.String()
	newColStr := newColBd.String()
	//oldColStr := oldColBd.String()
	//elseif (tg_op = 'DELETE') then
	//INSERT INTO ` + scmHisTbName + `("id","created_at","op_type",` + colStr + `) VALUES(nextval('` + scm + `.` + dorm.SequenceIdName + `')||'',now(),3,` + oldColStr + `);

	funcSql := `create or replace function ` + scmFuncName + `() returns trigger as $$
    begin
        if (tg_op = 'INSERT') then
            INSERT INTO ` + scmHisTbName + `("id","created_at","op_type","opt_user_id","opt_dept_id",` + colStr + `) VALUES(nextval('` + scm + `.` + dorm.SequenceIdName + `')||'',now(),1,NEW."creator_id",NEW."creator_dept_id",` + newColStr + `);
        elsif (tg_op = 'UPDATE') then
            INSERT INTO ` + scmHisTbName + `("id","created_at","op_type","opt_user_id","opt_dept_id",` + colStr + `) VALUES(nextval('` + scm + `.` + dorm.SequenceIdName + `')||'',now(),2,NEW."updater_id",NEW."updater_dept_id",` + newColStr + `);
        
        end if;
        RETURN NULL;
    end;
$$ language plpgsql;`

	triggerSql1 := fmt.Sprintf(`create or replace trigger %s after insert or update or delete on %s for each row execute function %s()`, triggerName, scmOriginTbName, scmFuncName)
	triggerSql2 := fmt.Sprintf(`create trigger %s after insert or update or delete on %s for each row execute function %s()`, triggerName, scmOriginTbName, scmFuncName)

	var vs int
	db.Raw(`SELECT current_setting('server_version_num')::integer`).Scan(&vs)
	logger.Infof("current database version is : %d", vs)

	return db.Transaction(func(tx *gorm.DB) error {
		tx2 := tx.Exec(funcSql)
		if tx2.Error != nil {
			return tx2.Error
		}
		if vs >= 140000 {
			return tx.Exec(triggerSql1).Error
		}
		return tx.Exec(triggerSql2).Error
	})

}

func createDmTrigger(db *gorm.DB, scm string, tableName string) error {
	scmHisTbName := fmt.Sprintf(`"%s"."%s%s"`, scm, tableName, constants.TableHistorySubFix)
	scmOriginTbName := fmt.Sprintf(`"%s"."%s"`, scm, tableName)
	triggerName := fmt.Sprintf(`"%s"."trigger_%s%s"`, scm, tableName, constants.TableHistorySubFix)

	col := dorm.FetchTableColumns(db, scm, tableName)
	colBd := strings.Builder{}
	newColBd := strings.Builder{}
	oldColBd := strings.Builder{}

	for i, v := range col {
		if i > 0 {
			colBd.WriteString(",")
			newColBd.WriteString(",")
			oldColBd.WriteString(",")
		}
		colBd.WriteString(`"h_` + v + `"`)
		newColBd.WriteString(`:NEW."` + v + `"`)
		oldColBd.WriteString(`:OLD."` + v + `"`)
	}
	colStr := colBd.String()
	newColStr := newColBd.String()
	//oldColStr := oldColBd.String()
	//ELSIF DELETING THEN
	//INSERT INTO ` + scmHisTbName + `("id","created_at","op_type",` + colStr + `) VALUES(TO_CHAR("` + scm + `"."` + dorm.SequenceIdName + `".NEXTVAL),CURRENT_TIMESTAMP(),3,` + oldColStr + `);

	trigger := `CREATE OR REPLACE TRIGGER ` + triggerName + `
    AFTER INSERT OR UPDATE OR DELETE ON ` + scmOriginTbName + ` FOR EACH ROW
DECLARE
BEGIN
    IF INSERTING THEN
        INSERT INTO ` + scmHisTbName + `("id","created_at","op_type","opt_user_id","opt_dept_id",` + colStr + `) VALUES(TO_CHAR("` + scm + `"."` + dorm.SequenceIdName + `".NEXTVAL),CURRENT_TIMESTAMP(),1,:NEW."creator_id",:NEW."creator_dept_id",` + newColStr + `);
    ELSIF UPDATING THEN
        INSERT INTO ` + scmHisTbName + `("id","created_at","op_type","opt_user_id","opt_dept_id",` + colStr + `) VALUES(TO_CHAR("` + scm + `"."` + dorm.SequenceIdName + `".NEXTVAL),CURRENT_TIMESTAMP(),2,:NEW."updater_id",:NEW."updater_dept_id",` + newColStr + `);
    END IF;
END;`
	return db.Exec(trigger).Error
}

func createMysqlTrigger(db *gorm.DB, scm, tableName string) error {
	var (
		tgInsertName = fmt.Sprintf("`%s`.`trigger_%s%s_insert`", scm, tableName, constants.TableHistorySubFix)
		tgUpdateName = fmt.Sprintf("`%s`.`trigger_%s%s_update`", scm, tableName, constants.TableHistorySubFix)
		//tgDeleteName    = fmt.Sprintf("`%s`.`trigger_%s%s_delete`", scm, tableName, TableSub)
		scmHisTbName    = fmt.Sprintf("`%s`.`%s%s`", scm, tableName, constants.TableHistorySubFix)
		scmOriginTbName = fmt.Sprintf("`%s`.`%s`", scm, tableName)
	)

	col := dorm.FetchTableColumns(db, scm, tableName)
	colBd := strings.Builder{}
	newColBd := strings.Builder{}
	oldColBd := strings.Builder{}

	for i, v := range col {
		if i > 0 {
			colBd.WriteString(",")
			newColBd.WriteString(",")
			oldColBd.WriteString(",")
		}
		colBd.WriteString("`h_" + v + "`")
		newColBd.WriteString("NEW.`" + v + "`")
		oldColBd.WriteString("OLD.`" + v + "`")
	}
	colStr := colBd.String()
	newColStr := newColBd.String()
	//oldColStr := oldColBd.String()

	dropTgInsert := `DROP TRIGGER IF EXISTS ` + tgInsertName
	dropTgUpdate := `DROP TRIGGER IF EXISTS ` + tgUpdateName
	//dropTgDelete := `DROP TRIGGER IF EXISTS ` + tgDeleteName

	tgInsert := `
	CREATE TRIGGER ` + tgInsertName + ` AFTER INSERT ON ` + scmOriginTbName + `
	FOR EACH ROW
	BEGIN
        INSERT INTO ` + scmHisTbName + "(`id`,`created_at`,`op_type`,`opt_user_id`,`opt_dept_id`," + colStr + `) VALUES (concat(nextval('` + scm + `.` + dorm.SequenceIdName + `'),''),now(),1,` + "NEW.`creator_id`,NEW.`creator_dept_id`," + newColStr + `);
	END`

	tgUpdate := `
	CREATE TRIGGER ` + tgUpdateName + ` AFTER UPDATE ON ` + scmOriginTbName + `
	FOR EACH ROW
	BEGIN
        INSERT INTO ` + scmHisTbName + "(`id`,`created_at`,`op_type`,`opt_user_id`,`opt_dept_id`," + colStr + `) VALUES (concat(nextval('` + scm + `.` + dorm.SequenceIdName + `'),''),now(),2,` + "NEW.`updater_id`,NEW.`updater_dept_id`," + newColStr + `);
	END`

	//tgDelete := `
	//CREATE TRIGGER ` + tgDeleteName + ` AFTER DELETE ON ` + scmOriginTbName + `
	//FOR EACH ROW
	//BEGIN
	//    INSERT INTO ` + scmHisTbName + "(`id`,`created_at`,`op_type`," + colStr + `) VALUES (concat(nextval('` + scm + `.` + dorm.SequenceIdName + `'),''),now(),3,` + oldColStr + `);
	//END`

	cl := func(tx *gorm.DB) error {
		return caller.NewCaller().
			Call(func(cl *caller.Caller) error {
				return tx.Exec(`LOCK TABLES ` + scmOriginTbName + ` WRITE`).Error
			}).
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
			}).
			//Call(func(cl *caller.Caller) error {
			//	return tx.Exec(dropTgDelete).Error
			//}).
			//Call(func(cl *caller.Caller) error {
			//	return tx.Exec(tgDelete).Error
			//}).
			Call(func(cl *caller.Caller) error {
				return tx.Exec(`UNLOCK TABLES`).Error
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
