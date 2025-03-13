package dorm

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/utils"
	"gorm.io/gorm"
	"strings"
)

func createIndex(db *gorm.DB, tbName string, idxType string, cols ...string) error {

	if len(cols) < 1 {
		return nil
	}

	var (
		scm         = GetDbSchema(db)
		dbType      = GetDbType(db)
		scmTbName   = Quote(dbType, scm, tbName)
		idxTypeList = []string{"ARRAY", "CONTEXT"}
	)
	if utils.ContainAny(idxTypeList, strings.ToUpper(idxType)) && dbType != DaMeng {
		return errs.NewClientError(fmt.Sprintf("not support idxType %s yet for database %s", idxType, dbType.String()))
	}

	colList := make([]string, 0, len(cols))
	idxCol := strings.Builder{}
	for i, v := range cols {
		if i > 0 {
			idxCol.WriteByte(',')
		}
		if strings.Contains(v, "(") {
			split := strings.Split(v, "(")
			colList = append(colList, split[0])
			idxCol.WriteString(Quote(dbType, split[0]))
			idxCol.WriteByte('(')
			idxCol.WriteString(split[1])
			continue
		}
		colList = append(colList, v)
		idxCol.WriteString(Quote(dbType, v))
	}

	idxName := fmt.Sprintf("idx_%s_%s", tbName, strings.Join(colList, "_"))
	idxSql := fmt.Sprintf(`CREATE %s INDEX IF NOT EXISTS %s ON %s(%s)`, idxType, Quote(dbType, idxName), scmTbName, idxCol.String())

	if strings.ToUpper(idxType) == "CONTEXT" {
		idxSql = idxSql + " SYNC TRANSACTION"
	}

	switch dbType {
	case Mysql:
		exists, err := mysqlIdxExists(db, tbName, idxName)
		if !exists {
			//不支持IF NOT EXISTS
			idxSql = fmt.Sprintf(`CREATE %s INDEX %s ON %s(%s)`, idxType, Quote(dbType, idxName), scmTbName, idxCol.String())
			return db.Exec(idxSql).Error
		}
		return err
	default:
		return db.Exec(idxSql).Error
	}

}

func CreateIndex(db *gorm.DB, tbName string, cols ...string) error {
	return createIndex(db, tbName, "", cols...)
}
func CreateArrayIndex(db *gorm.DB, tbName string, cols ...string) error {
	return createIndex(db, tbName, "ARRAY", cols...)
}
func CreateContextIndex(db *gorm.DB, tbName string, cols ...string) error {
	return createIndex(db, tbName, "CONTEXT", cols...)
}
func CreateUniqueIndex(db *gorm.DB, tbName string, cols ...string) error {
	return createIndex(db, tbName, "UNIQUE", cols...)
}

func mysqlIdxExists(db *gorm.DB, tbName, idxName string) (bool, error) {
	var (
		scm = GetDbSchema(db)
	)
	ctSql := fmt.Sprintf(`SELECT COUNT(*) FROM information_schema.statistics 
                WHERE TABLE_SCHEMA = '%s'
                  AND TABLE_NAME = '%s' 
                  AND INDEX_NAME = '%s'`, scm, tbName, idxName)
	count := 0
	err := db.Raw(ctSql).Find(&count).Error
	return count > 0, err
}
