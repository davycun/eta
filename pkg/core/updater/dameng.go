package updater

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/db_table"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"gorm.io/gorm"
	"strings"
)

func createDmUpdaterTrigger(db *gorm.DB, tableName string) error {
	var (
		scm         = dorm.GetDbSchema(db)
		dbUser      = dorm.GetDbUser(db)
		scmTbName   = fmt.Sprintf(`"%s"."%s"`, scm, tableName)
		triggerName = fmt.Sprintf(`"%s"."trigger_%s_updater"`, scm, tableName)
		bd          = strings.Builder{}
		bdF         = strings.Builder{}
		cols        []db_table.Column
		err         = db_table.FetchColumns(db, tableName, &cols)
	)
	if err != nil {
		return err
	}

	for _, v := range cols {
		bd.WriteString(`
			updater := jsonb_set(updater,'{"` + v.ColName + `"}','"'||:NEW."` + entity.UpdaterIdDbName + `"||'"');`)
		bd.WriteString("\n")

		//更新的时候由于自定义类型不能用比较所以去掉不支持
		if utils.ContainAny(entity.DefaultVertexColumns, v.ColName) || utils.ContainAny(notSupportType, strings.ToLower(v.ColType)) {
			continue
		}
		bdF.WriteString(`
			IF :OLD."` + v.ColName + `" != :NEW."` + v.ColName + `" or (:OLD."` + v.ColName + `" is null and NEW."` + v.ColName + `" is not null) THEN
				flag := 1;
				IF :NEW."` + entity.UpdaterIdDbName + `" NOT MEMBER OF updater_ids THEN
					ids_flag := 1;
					updater_ids.EXTEND(1);
					updater_ids(updater_ids.COUNT()) := :NEW."` + entity.UpdaterIdDbName + `";
				END IF;
				updater := jsonb_set(updater,'{` + v.ColName + `}','"'||:NEW."` + entity.UpdaterIdDbName + `"||'"');
			END IF;`)
		bdF.WriteString("\n")
	}

	trigger := `CREATE OR REPLACE TRIGGER ` + triggerName + `
    BEFORE INSERT OR UPDATE ON ` + scmTbName + ` FOR EACH ROW
DECLARE
	updater CLOB;
	updater_ids ` + dbUser + `.ARR_STR;
    flag bit;
	ids_flag bit;
BEGIN
    IF INSERTING THEN
		IF :NEW."` + entity.FieldUpdaterDbName + `" is null THEN
			updater := '{}';
			` + bd.String() + `
			:NEW."` + entity.FieldUpdaterDbName + `" := updater;
		END IF;
		IF :NEW."` + entity.FieldUpdaterIdsDbName + `" is null and :NEW."` + entity.UpdaterIdDbName + `" is not null THEN
			updater_ids := ` + dbUser + `.ARR_STR(:NEW."` + entity.UpdaterIdDbName + `");
			:NEW."` + entity.FieldUpdaterIdsDbName + `" := ` + dbUser + `.ARR_STR_CLS(updater_ids);
		END IF;

    ELSEIF UPDATING THEN
        flag := 0;
		updater := :OLD."` + entity.FieldUpdaterDbName + `";
		IF updater is null or updater='' THEN
			updater := '{}';
		END IF;
        updater_ids := ` + dbUser + `.ARR_STR();
		IF :OLD."` + entity.FieldUpdaterIdsDbName + `" is not null THEN
			updater_ids := :OLD."` + entity.FieldUpdaterIdsDbName + `".V;
		END IF;

		` + bdF.String() + `
		IF flag THEN 
			:NEW."` + entity.FieldUpdaterDbName + `" := updater;
		END IF;
		IF ids_flag THEN 
			:NEW."` + entity.FieldUpdaterIdsDbName + `" := ` + dbUser + `.ARR_STR_CLS(updater_ids);
		END IF;

    END IF;
END;`
	return db.Exec(trigger).Error
}
