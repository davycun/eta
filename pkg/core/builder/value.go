package builder

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/duke-git/lancet/v2/slice"
	"strings"
)

type ValueBuilder struct {
	dbType     dorm.DbType
	values     []string
	distinct   bool
	valueAlias string
}

func NewValueBuilder(dbType dorm.DbType, valueAlias string, values ...string) *ValueBuilder {
	if valueAlias == "" {
		valueAlias = entity.IdDbName
	}
	bd := &ValueBuilder{
		dbType:     dbType,
		values:     values,
		valueAlias: valueAlias,
	}
	return bd
}

func (b *ValueBuilder) SetDistinct(distinct bool) *ValueBuilder {
	b.distinct = distinct
	return b
}

func (b *ValueBuilder) Build() (listSql, countSql string, err error) {

	if len(b.values) < 1 {
		return "", "", nil
	}

	bd := strings.Builder{}
	switch b.dbType {
	case dorm.DaMeng:
		dis := b.distinct
		//TODO 达梦 当大于等于4680的时候就会报错，错误是类型不匹配
		size := 4000
		if len(b.values) > size {
			dis = false
		}
		ids := make([]string, 0, len(b.values))
		for i, v := range b.values {
			ids = append(ids, v)
			if (i > 0 && i%size == 0) || i == (len(b.values)-1) {
				if i > size {
					bd.WriteString(" union ")
				}
				bd.WriteString(buildValueToTableDm(dis, b.valueAlias, ids...))
				ids = ids[0:0]
			}
		}
	case dorm.PostgreSQL:
		bd.WriteString(`select `)
		if b.distinct {
			bd.WriteString(` distinct `)
		}
		bd.WriteString(` unnest(array[`)
		for i, v := range b.values {
			if i > 0 {
				bd.WriteByte(',')
			}
			bd.WriteString(fmt.Sprintf(`'%s'`, v))
		}
		bd.WriteString(fmt.Sprintf(`]) as "%s"`, b.valueAlias))
	case dorm.Mysql, dorm.Doris:
		bd.WriteString(`select `)
		if b.distinct {
			bd.WriteString(` distinct `)
		}
		bd.WriteString(fmt.Sprintf("`%s` from ", b.valueAlias))
		bd.WriteString(` JSON_TABLE('[`)
		vs := slice.Map(b.values, func(index int, item string) string {
			return fmt.Sprintf(`{"v":"%s"}`, item)
		})
		bd.WriteString(strings.Join(vs, ","))
		bd.WriteString(fmt.Sprintf(`]', "$[*]" COLUMNS(%s VARCHAR(255) PATH "$.v")`, b.valueAlias))
		bd.WriteString(`) AS id`)
	}

	listSql = bd.String()
	tbName := dorm.Quote(b.dbType, "t")
	countSql = fmt.Sprintf("with %s as (%s) select count(*) from %s", tbName, listSql, tbName)

	return
}

func buildValueToTableDm(distinct bool, valuesAlias string, values ...string) string {
	bd := strings.Builder{}
	bd.WriteString(`select `)
	if distinct {
		bd.WriteString(` distinct `)
	}
	bd.WriteString(fmt.Sprintf(`value as "%s" from `, valuesAlias))
	bd.WriteString(` jsonb_array_elements_text('[`)
	for i, v := range values {
		if i > 0 {
			bd.WriteByte(',')
		}
		bd.WriteString(fmt.Sprintf(`"%s"`, v))
	}
	bd.WriteString(`]')`)
	return bd.String()
}
