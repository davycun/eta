package dorm

import (
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
)

func FetchTableColumns(db *gorm.DB, scm, tbName string) []string {
	var (
		sq = ""
		rs []string
	)
	switch db.Dialector.Name() {
	case DaMeng.String():
		sq = buildDmTableColumnsSql(scm, tbName)
	case PostgreSQL.String():
		sq = buildPgTableColumnsSql(scm, tbName)
	case Mysql.String():
		sq = buildMySqlTableColumnsSql(scm, tbName)
	default:
		return rs
	}
	err := db.Raw(sq).Find(&rs).Error
	if err != nil {
		logger.Errorf("fetch table columns err %s", err)
	}
	return rs
}
func buildDmTableColumnsSql(scm, tbName string) string {
	return `with scm as (
    SELECT ID FROM SYS.SYSOBJECTS WHERE TYPE$ = 'SCH' AND NAME = '` + scm + `'
),tb as (
        SELECT ID,SCHID FROM SYS.SYSOBJECTS WHERE TYPE$ = 'SCHOBJ' AND SUBTYPE$ IN ('UTAB', 'STAB', 'VIEW') AND NAME = '` + tbName + `'
    )
select  col.NAME from SYS.SYSCOLUMNS col,scm,tb where col.ID=tb.ID and tb.SCHID=scm.ID`
}
func buildPgTableColumnsSql(scm, tbName string) string {
	return `SELECT column_name FROM INFORMATION_SCHEMA.columns WHERE table_schema = '` + scm + `' AND table_name = '` + tbName + `'`
}
func buildMySqlTableColumnsSql(scm, tbName string) string {
	return `SELECT column_name FROM INFORMATION_SCHEMA.columns WHERE table_schema = '` + scm + `' AND table_name = '` + tbName + `'`
}

func FetchById(id string, db *gorm.DB, data any, columns ...string) error {
	tx := db.Model(data)
	if len(columns) > 0 {
		tx = tx.Select(columns)
	}
	return tx.Where(map[string]any{"id": id}).First(data).Error
}
