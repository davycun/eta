package ra

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
	"strings"
)

func createPostgresTrigger(db *gorm.DB, scm, tableName string, raFields []string) error {
	if len(raFields) <= 0 {
		return nil
	}
	var (
		scmFuncName  = fmt.Sprintf(`"%s"."func_%s_ra"`, scm, tableName)
		triggerName  = fmt.Sprintf(`"trigger_%s_ra"`, tableName)
		scmTableName = fmt.Sprintf(`"%s"."%s"`, scm, tableName)
		raString     = strings.Join(slice.Map(raFields, func(_ int, v string) string { return `NEW."` + v + `"` }), `||' '||`)
	)

	funcSql := `CREATE OR REPLACE FUNCTION ` + scmFuncName + `() RETURNS TRIGGER AS $$
	BEGIN
        IF (tg_op = 'INSERT') THEN
			IF NEW."` + entity.RaContentDbName + `" is null THEN
				NEW."` + entity.RaContentDbName + `" := ` + raString + `;
			END IF;
			
        ELSIF (tg_op = 'UPDATE') THEN
			NEW."` + entity.RaContentDbName + `" := ` + raString + `;
			
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
