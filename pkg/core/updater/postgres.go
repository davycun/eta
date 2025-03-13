package updater

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/db_table"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"gorm.io/gorm"
	"strings"
)

func createPostgresUpdaterTrigger(db *gorm.DB, scm, tableName string) error {
	var (
		scmFuncName  = fmt.Sprintf(`"%s"."func_%s_updater"`, scm, tableName)
		triggerName  = fmt.Sprintf(`"trigger_%s_updater"`, tableName)
		scmTableName = fmt.Sprintf(`"%s"."%s"`, scm, tableName)
		cols         []db_table.Column
		bd           = strings.Builder{}
		bdF          = strings.Builder{}
		err          = db_table.FetchColumns(db, tableName, &cols)
	)
	if err != nil {
		return err
	}
	for _, v := range cols {

		bd.WriteString(fmt.Sprintf(`updater := jsonb_concat(updater,jsonb_build_object('%s',NEW."%s"));`, v.ColName, entity.UpdaterIdDbName))
		bd.WriteString("\n")

		//更新的时候由于自定义类型不能用比较所以去掉不支持
		if utils.ContainAny(entity.DefaultVertexColumns, v.ColName) || utils.ContainAny(notSupportType, strings.ToLower(v.ColType)) {
			continue
		}

		bdF.WriteString(fmt.Sprintf(`IF OLD."%s" != NEW."%s" THEN
			flag := true;
			IF not updater_ids @> tmp THEN
				ids_flag := true;
				updater_ids := array_append(updater_ids,NEW."`+entity.UpdaterIdDbName+`");
			END IF;
			updater := jsonb_concat(updater,jsonb_build_object('%s',NEW."%s"));
		END IF;`, v.ColName, v.ColName, v.ColName, entity.UpdaterIdDbName))
		bdF.WriteString("\n")
	}

	funcSql := `CREATE OR REPLACE FUNCTION ` + scmFuncName + `() RETURNS TRIGGER AS $$
	DECLARE
    	updater jsonb := '{}'::jsonb;
		updater_ids varchar[];
		tmp   varchar[];
		flag bool := false;
		ids_flag bool := false;
	BEGIN
        IF (tg_op = 'INSERT') THEN
			IF NEW."` + entity.FieldUpdaterDbName + `" is null THEN
				` + bd.String() + `
				NEW."` + entity.FieldUpdaterDbName + `" := updater;
			END IF;

			IF NEW."` + entity.FieldUpdaterIdsDbName + `" is null THEN
				updater_ids := array_append(updater_ids, NEW."` + entity.UpdaterIdDbName + `");
				NEW."` + entity.FieldUpdaterIdsDbName + `" := updater_ids;
			END IF;
			
        ELSIF (tg_op = 'UPDATE') THEN
			IF OLD."` + entity.FieldUpdaterDbName + `" is not null THEN
				updater := OLD."` + entity.FieldUpdaterDbName + `";
			END IF;
			IF OLD."` + entity.FieldUpdaterIdsDbName + `" is not null THEN
				updater_ids := OLD."` + entity.FieldUpdaterIdsDbName + `";
			END IF;
			tmp := array_append(tmp,NEW."` + entity.UpdaterIdDbName + `");
			
        	` + bdF.String() + `
			
			IF flag THEN
				NEW."` + entity.FieldUpdaterDbName + `" := updater;
			END IF;
			IF ids_flag THEN
				NEW."` + entity.FieldUpdaterIdsDbName + `" := updater_ids;
			END IF;
			
        END IF;
        RETURN NEW;
    END;
$$ language plpgsql;`

	triggerSql1 := fmt.Sprintf(`create or replace trigger %s before insert or update on %s for each row execute function %s()`, triggerName, scmTableName, scmFuncName)
	triggerSql2 := fmt.Sprintf(`create trigger %s before insert or update on %s for each row execute function %s()`, triggerName, scmTableName, scmFuncName)

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
