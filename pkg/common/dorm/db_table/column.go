package db_table

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type Column struct {
	SchName   string `json:"sch_name,omitempty" gorm:"column:sch_name"`
	TbName    string `json:"tb_name,omitempty" gorm:"column:tb_name"`
	ColName   string `json:"col_name,omitempty" gorm:"column:col_name"`
	ColType   string `json:"col_type,omitempty" gorm:"column:col_type"`
	Precision int    `json:"precision,omitempty" gorm:"column:numeric_precision"`
	Scale     int    `json:"scale,omitempty" gorm:"column:numeric_scale"`
}

func FetchColumns(db *gorm.DB, tableName string, cols *[]Column) error {
	var (
		sq     = ""
		scm    = dorm.GetDbSchema(db)
		dbType = dorm.GetDbType(db)
	)
	switch dbType {
	case dorm.DaMeng:
		sq = buildDmTableColumnsSql(scm, tableName)
	case dorm.PostgreSQL:
		sq = buildPgTableColumnsSql(scm, tableName)
	case dorm.Mysql:
		sq = buildMySqlTableColumnsSql(scm, tableName)
	default:
		return nil
	}

	err := db.Raw(sq).Find(cols).Error
	if err != nil {
		return err
	}

	//TODO 这段代码是为了获取类型名称，比如CLASSxxxx对应的具体类型，如果放开的话再updater包需要进行修改
	if dbType == dorm.DaMeng {
		for i, v := range *cols {
			tp := strings.ToLower(v.ColType)
			if !strings.Contains(tp, "class") {
				continue
			}
			after, found := strings.CutPrefix(tp, "class")
			if !found {
				continue
			}
			id, err1 := strconv.ParseInt(after, 10, 64)
			if err1 != nil {
				logger.Errorf("parse err %s", err1)
				continue
			}
			obj := fetchObjectById(db, id)
			(*cols)[i].ColType = obj.Name
		}
	}

	return nil
}

func buildDmTableColumnsSql(scm, tbName string) string {
	return `with scm as (
    SELECT ID,NAME FROM SYS.SYSOBJECTS WHERE TYPE$ = 'SCH' AND NAME = '` + scm + `'
),tb as (
        SELECT ID,SCHID,NAME FROM SYS.SYSOBJECTS WHERE TYPE$ = 'SCHOBJ' AND SUBTYPE$ IN ('UTAB', 'STAB', 'VIEW') AND NAME = '` + tbName + `'
    )
select  scm.NAME as "sch_name",tb.NAME as "tb_name", col.NAME as "col_name",col.TYPE$ as "col_type",LENGTH$ as "numeric_precision"
from SYS.SYSCOLUMNS col,scm,tb where col.ID=tb.ID and tb.SCHID=scm.ID`
}
func buildPgTableColumnsSql(scm, tbName string) string {
	//return `SELECT column_name FROM INFORMATION_SCHEMA.columns WHERE table_schema = '` + scm + `' AND table_name = '` + tbName + `'`
	return `SELECT table_schema as "sch_name",table_name as "tb_name",column_name as "col_name",udt_name as "col_type",numeric_precision as "numeric_precision","numeric_scale" 
FROM INFORMATION_SCHEMA.columns WHERE table_schema = '` + scm + `' AND table_name = '` + tbName + `'`
}
func buildMySqlTableColumnsSql(scm, tbName string) string {
	return `SELECT 
				table_schema as ` + "`sch_name`" + `,
				table_name as ` + "`tb_name`" + `,
				column_name as ` + "`col_name`" + `,
				column_type as ` + "`col_type`" + `,
				numeric_precision as ` + "`numeric_precision`" + `,
				numeric_scale as ` + "`numeric_scale`" + `
			FROM INFORMATION_SCHEMA.columns 
			WHERE table_schema = '` + scm + `' AND table_name = '` + tbName + `'`
}
