package dorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

func dmBatchUpdate(db *gorm.DB, values clause.Values, cols ...string) {

	db.Statement.SQL = strings.Builder{}
	db.Statement.WriteString("UPDATE ")
	//schema的支持
	if db.Statement.TableExpr != nil && db.Statement.TableExpr.SQL != "" {
		db.Statement.TableExpr.Build(db.Statement)
	} else {
		db.Statement.WriteQuoted(db.Statement.Table)
	}

	db.Statement.WriteString(" set ")

	st := clause.AssignmentColumns(cols)
	st.Build(db.Statement)

	db.Statement.WriteString(" from (")

	for idx, value := range values.Values {

		if idx > 0 {
			db.Statement.WriteString(" UNION ALL ")
		}

		db.Statement.WriteString("SELECT ")
		db.Statement.AddVar(db.Statement, value...)
		db.Statement.WriteString(" FROM DUAL")
	}

	db.Statement.WriteString(`) AS "excluded" (`)
	for idx, column := range values.Columns {
		if idx > 0 {
			db.Statement.WriteByte(',')
		}
		db.Statement.WriteQuoted(column.Name)
	}
	db.Statement.WriteString(") where ")

	var where clause.Where
	for _, column := range []string{"id"} {
		where.Exprs = append(where.Exprs, clause.Eq{
			Column: clause.Column{Table: db.Statement.Table, Name: column},
			Value:  clause.Column{Table: "excluded", Name: column},
		})
	}
	where.Build(db.Statement)

}
