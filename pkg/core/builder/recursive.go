package builder

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"slices"
)

type RecursiveSqlBuilder struct {
	CteSqlBuilder
	cteName          string
	cteRsName        string
	isUp             bool
	idColumnName     string
	parentColumnName string
	recursiveFilters []filter.Filter
	depth            int //递归的深度
}

func NewRecursiveSqlBuilder(dbType dorm.DbType, scm string, tableName string) *RecursiveSqlBuilder {
	b := &RecursiveSqlBuilder{
		CteSqlBuilder: *NewCteSqlBuilder(dbType, scm, tableName),
		cteName:       "cte",
		cteRsName:     "cte_rs",
	}
	return b
}

func (b *RecursiveSqlBuilder) SetUp(up bool) *RecursiveSqlBuilder {
	b.isUp = up
	return b
}
func (b *RecursiveSqlBuilder) SetCteName(cteName string) *RecursiveSqlBuilder {
	b.cteName = cteName
	b.cteRsName = cteName + "_rs"
	return b
}
func (b *RecursiveSqlBuilder) SetIdColumnName(idColumnName string) *RecursiveSqlBuilder {
	b.idColumnName = idColumnName
	return b
}
func (b *RecursiveSqlBuilder) SetParentColumnName(parentColumnName string) *RecursiveSqlBuilder {
	b.parentColumnName = parentColumnName
	return b
}
func (b *RecursiveSqlBuilder) SetDepth(depth int) *RecursiveSqlBuilder {
	b.depth = depth
	return b
}
func (b *RecursiveSqlBuilder) AddRecursiveFilter(filters ...filter.Filter) *RecursiveSqlBuilder {
	b.recursiveFilters = append(b.recursiveFilters, filters...)
	return b
}
func (b *RecursiveSqlBuilder) check() *RecursiveSqlBuilder {
	b.SqlBuilder.check()
	if b.idColumnName == "" {
		b.idColumnName = "id"
	}
	if b.parentColumnName == "" {
		b.parentColumnName = "parent_id"
	}
	if b.cteName == "" {
		b.cteName = "cte"
		b.cteRsName = "cte_rs"
	}
	if len(b.columns) < 1 {
		b.AddTableColumn(b.GetTableName(), b.idColumnName, b.parentColumnName)
	}

	return b
}
func (b *RecursiveSqlBuilder) Clone() RecursiveSqlBuilder {
	cs := b.CteSqlBuilder.Clone()
	rs := RecursiveSqlBuilder{}
	rs.cteName = b.cteName
	rs.cteRsName = b.cteRsName
	rs.isUp = b.isUp
	rs.idColumnName = b.idColumnName
	rs.parentColumnName = b.parentColumnName
	rs.recursiveFilters = slices.Clone(b.recursiveFilters)
	rs.depth = b.depth

	rs.CteSqlBuilder = cs
	return rs
}
func (b *RecursiveSqlBuilder) Build() (listSql, countSql string, err error) {

	if b.check().Err != nil {
		return "", "", b.Err
	}
	var (
		depthName = "depth"
		bs        = b.Clone() //为了支持Build幂等需要用Copy出来的Builder进行构建
	)
	if len(bs.recursiveFilters) > 0 {
		//递归部分的with
		pBd := NewSqlBuilder(bs.dbType, bs.schema, bs.GetTableName()).AddFilter(bs.recursiveFilters...).AddColumn(bs.idColumnName, bs.parentColumnName)
		cBd := NewSqlBuilder(bs.dbType, bs.schema, bs.GetTableName()).AddColumn(bs.idColumnName, bs.parentColumnName)

		if bs.depth > 0 {
			pBd.AddExprColumn(expr.ExpColumn{
				Expression: expr.Expression{
					Expr: "?",
					Vars: []expr.ExpVar{
						{Type: expr.VarTypeValue, Value: 1},
					},
				},
				Alias: depthName,
			})
			cBd.AddTableExprColumn(bs.cteName, expr.ExpColumn{
				Expression: expr.Expression{
					Expr: "? + 1",
					Vars: []expr.ExpVar{
						{Type: expr.VarTypeColumn, Value: depthName},
					},
				},
				Alias: depthName,
			})
			cBd.AddTableFilter(bs.cteName, filter.Filter{Column: depthName, Operator: filter.LTE, Value: bs.depth})
		}
		if bs.isUp {
			cBd.Join("", bs.cteName, bs.parentColumnName, bs.GetTableName(), bs.idColumnName)
		} else {
			cBd.Join("", bs.cteName, bs.idColumnName, bs.GetTableName(), bs.parentColumnName)
		}
		pBd.UnionAll(cBd)

		if b.depth > 0 {
			bs.CteSqlBuilder.WithRecursive(bs.cteName, pBd, bs.idColumnName, bs.parentColumnName, depthName)
		} else {
			bs.CteSqlBuilder.WithRecursive(bs.cteName, pBd, bs.idColumnName, bs.parentColumnName)
		}

		//从递归部分取出ID字段，并且去重的with
		rsBd := NewSqlBuilder(bs.dbType, "", bs.cteName).SetDistinct(true).AddColumn(bs.idColumnName)
		bs.CteSqlBuilder.With(bs.cteRsName, rsBd)

		//join上RS的ID取出最后的结果
		bs.CteSqlBuilder.Join("", bs.cteRsName, bs.idColumnName, bs.GetTableName(), bs.idColumnName)
	}

	listSql, countSql, err = bs.CteSqlBuilder.Build()
	err = bs.Err

	return
}
