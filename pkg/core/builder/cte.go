package builder

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/utils"
	"slices"
	"strings"
)

type CteSqlBuilder struct {
	SqlBuilder
	bdList []builder
}

func NewCteSqlBuilder(dbType dorm.DbType, schemaName, tableName string) *CteSqlBuilder {
	bd := &CteSqlBuilder{
		SqlBuilder: SqlBuilder{
			schema:      schemaName,
			dbType:      dbType,
			tableName:   tableName,
			countColumn: "*",
		},
	}
	return bd
}

func (c *CteSqlBuilder) With(cteName string, bd Builder, cteCol ...string) *CteSqlBuilder {
	c.bdList = append(c.bdList, builder{name: cteName, cols: cteCol, sqlBd: bd})
	return c
}
func (c *CteSqlBuilder) WithRecursive(cteName string, bd Builder, cteCol ...string) *CteSqlBuilder {
	//cteCol must be not empty
	if len(cteCol) < 1 {
		c.Err = fmt.Errorf("WithRecursive cteCol must be not empty")
		return c
	}
	c.bdList = append(c.bdList, builder{name: cteName, cols: cteCol, sqlBd: bd, recursive: true})
	return c
}
func (c *CteSqlBuilder) Clone() CteSqlBuilder {
	sb := c.SqlBuilder.Clone()
	cs := CteSqlBuilder{}
	cs.bdList = slices.Clone(c.bdList)
	cs.SqlBuilder = sb
	return cs
}
func (c *CteSqlBuilder) Build() (listSql, countSql string, err error) {

	sq := make([]string, 0, len(c.bdList))
	for _, v := range c.bdList {
		recursive := ""
		if v.recursive {
			switch c.dbType {
			case dorm.Mysql, dorm.PostgreSQL, dorm.Doris:
				recursive = "recursive"
			}
		}
		listSql, _, err = v.sqlBd.Build()
		if err != nil {
			c.Err = err
			return
		}
		if len(v.cols) > 0 {
			sq = utils.AppendNoEmpty(sq, fmt.Sprintf("%s %s(%s) as (%s)",
				recursive,
				dorm.Quote(c.dbType, v.name),
				dorm.JoinColumns(c.dbType, "", v.cols),
				listSql))
		} else {
			sq = utils.AppendNoEmpty(sq, fmt.Sprintf("%s as (%s)", dorm.Quote(c.dbType, v.name), listSql))
		}
	}

	withSql := ""
	if len(sq) > 0 {
		withSql = "with " + strings.Join(sq, ",") + " "
	}
	listSql, countSql, err = c.SqlBuilder.Build()

	return withSql + listSql, withSql + countSql, err
}

type builder struct {
	name      string
	cols      []string
	recursive bool
	sqlBd     Builder
}
