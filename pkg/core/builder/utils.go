package builder

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/duke-git/lancet/v2/slice"
	"strings"
)

func BuildInSql(values ...string) string {

	if len(values) < 1 {
		return ""
	}
	bd := strings.Builder{}
	bd.WriteByte('(')
	for i, v := range values {
		if i > 0 {
			bd.WriteByte(',')
		}
		bd.WriteString(fmt.Sprintf(`'%s'`, v))
	}
	bd.WriteByte(')')
	return bd.String()
}

// BuildValueToTableSql 把 values 转成一张表
func BuildValueToTableSql(dbType dorm.DbType, distinct bool, values ...string) string {
	var (
		bd         = strings.Builder{}
		valueAlias = "id"
	)
	switch dbType {
	case dorm.DaMeng:
		dis := distinct
		//TODO 达梦 当大于等于4680的时候就会报错，错误是类型不匹配
		size := 4000
		if len(values) > size {
			dis = false
		}
		ids := make([]string, 0, len(values))
		for i, v := range values {
			ids = append(ids, v)
			if (i > 0 && i%size == 0) || i == (len(values)-1) {
				if i > size {
					bd.WriteString(" union ")
				}
				bd.WriteString(buildValueToTableDm(dis, valueAlias, ids...))
				ids = ids[0:0]
			}

		}
	case dorm.PostgreSQL:
		bd.WriteString(`select `)
		if distinct {
			bd.WriteString(` distinct `)
		}
		bd.WriteString(` unnest(array[`)
		for i, v := range values {
			if i > 0 {
				bd.WriteByte(',')
			}
			bd.WriteString(fmt.Sprintf(`'%s'`, v))
		}
		bd.WriteString(fmt.Sprintf(`]) as "%s"`, valueAlias))
	case dorm.Mysql:
		bd.WriteString(`select `)
		if distinct {
			bd.WriteString(` distinct `)
		}
		bd.WriteString(fmt.Sprintf("`%s` from ", valueAlias))
		bd.WriteString(` JSON_TABLE('[`)
		vs := slice.Map(values, func(index int, item string) string {
			return fmt.Sprintf(`{"v":"%s"}`, item)
		})
		bd.WriteString(strings.Join(vs, ","))
		bd.WriteString(`]', "$[*]" COLUMNS(id VARCHAR(255) PATH "$.v")`)
		bd.WriteString(`) AS id`)
	}
	return bd.String()
}
