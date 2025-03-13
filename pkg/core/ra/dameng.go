package ra

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/duke-git/lancet/v2/slice"
	"gorm.io/gorm"
	"strings"
)

func createDmTrigger(db *gorm.DB, tableName string, raFields []string) error {
	if len(raFields) <= 0 {
		return nil
	}
	var (
		scm         = dorm.GetDbSchema(db)
		scmTbName   = fmt.Sprintf(`"%s"."%s"`, scm, tableName)
		triggerName = fmt.Sprintf(`"%s"."trigger_%s_ra"`, scm, tableName)
		raString    = strings.Join(slice.Map(raFields, func(_ int, v string) string { return `:NEW."` + v + `"` }), `||' '||`)
	)

	trigger := `CREATE OR REPLACE TRIGGER ` + triggerName + `
    BEFORE INSERT OR UPDATE ON ` + scmTbName + ` FOR EACH ROW
BEGIN
    IF INSERTING THEN
		IF :NEW."` + entity.RaContentDbName + `" is null THEN
			:NEW."` + entity.RaContentDbName + `" := ` + raString + `;
		END IF;

    ELSEIF UPDATING THEN
		:NEW."` + entity.RaContentDbName + `" := ` + raString + `;

    END IF;
END;`
	return db.Exec(trigger).Error
}
