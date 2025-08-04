package builder

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"slices"
	"strings"
)

type AggregateSqlBuilder struct {
	CteSqlBuilder
	cteName          string
	aggregateColumns []dorm.AggregateColumn
	groupColumns     []string
	havingFilters    []filter.Having
}

func NewAggregateSqlBuilder(dbType dorm.DbType, schemaName, tableName string) *AggregateSqlBuilder {
	b := &AggregateSqlBuilder{
		CteSqlBuilder: *NewCteSqlBuilder(dbType, schemaName, tableName),
	}
	return b
}

func (b *AggregateSqlBuilder) SetCteName(cteName string) *AggregateSqlBuilder {
	b.cteName = cteName
	return b
}
func (b *AggregateSqlBuilder) AddAggregateColumn(aggColumn ...dorm.AggregateColumn) *AggregateSqlBuilder {
	b.aggregateColumns = append(b.aggregateColumns, aggColumn...)
	return b
}
func (b *AggregateSqlBuilder) AddGroupColumn(column ...string) *AggregateSqlBuilder {
	b.groupColumns = append(b.groupColumns, column...)
	return b
}
func (b *AggregateSqlBuilder) AddHavingFilter(filters ...filter.Having) *AggregateSqlBuilder {
	b.havingFilters = append(b.havingFilters, filters...)
	return b
}

func (b *AggregateSqlBuilder) check() *AggregateSqlBuilder {

	if b.GetTableName() == "" {
		b.Err = errors.New("tableName is empty")
		return b
	}
	if len(b.groupColumns) < 1 {
		b.Err = errors.New("groupColumns is empty")
		return b
	}
	if b.cteName == "" {
		b.cteName = "r"
	}

	return b
}

func (b *AggregateSqlBuilder) resolveOrderByString() string {
	return dorm.ResolveOrderByString(b.orderBy, "", "", true)
}

func (b *AggregateSqlBuilder) Clone() AggregateSqlBuilder {
	cs := b.CteSqlBuilder.Clone()
	as := AggregateSqlBuilder{}
	as.cteName = b.cteName
	as.aggregateColumns = slices.Clone(b.aggregateColumns)
	as.groupColumns = slices.Clone(b.groupColumns)
	as.havingFilters = slices.Clone(b.havingFilters)
	as.CteSqlBuilder = cs
	return as
}
func (b *AggregateSqlBuilder) Build() (listSql, countSql string, err error) {

	if b.check().Err != nil {
		return "", "", b.Err
	}

	bc := b.Clone()

	var (
		groupBy = " group by " + dorm.JoinColumns(bc.dbType, bc.GetTableName(), bc.groupColumns)
		orderBy = bc.resolveOrderByString()
		lmtSql  = bc.resolveLimitString()
	)

	//开始构建listSql
	//1、重置columns，聚合sql里不能有选择非聚合函数字段
	bc.columns = make([]tableColumns, 0, 1)
	for _, v := range bc.groupColumns {
		bc.AddColumn(v)
	}
	for _, v := range bc.aggregateColumns {
		exp := expr.ExpColumn{}
		exp.Expr = fmt.Sprintf("%s(?)", v.AggFunc)
		exp.Alias = v.Alias
		exp.Vars = []expr.ExpVar{{Type: expr.VarTypeColumn, Value: v.Column}}
		bc.AddExprColumn(exp)
	}

	//2、分页和limit内容需要放在group之后
	bc.SetMustPage(false)
	bc.Limit(0)
	bc.orderBy = make([]dorm.OrderBy, 0) //清空order by
	listSql, _, err = bc.CteSqlBuilder.Build()
	if err != nil {
		return
	}
	listSql = listSql + groupBy + orderBy + lmtSql

	bc = b.Clone()
	//开始构建countSql
	//1、重置columns
	bc.columns = make([]tableColumns, 0, 1)
	for _, v := range bc.groupColumns {
		//取一个聚合字段即可
		bc.AddColumn(v)
		break
	}
	//2、不能有分页
	bc.SetMustPage(false)
	bc.Limit(0)
	bc.orderBy = make([]dorm.OrderBy, 0) //清空order by
	listSqlForCount, _, err := bc.CteSqlBuilder.Build()
	listSqlForCount = listSqlForCount + groupBy
	if err != nil {
		return
	}

	//3、拼接CountSql
	bd := strings.Builder{}
	cte := dorm.Quote(bc.dbType, bc.cteName)
	bd.WriteString(fmt.Sprintf("with %s as (", cte))
	bd.WriteString(listSqlForCount)
	bd.WriteString(")")
	bd.WriteString(fmt.Sprintf("select count(*) from %s", cte))
	countSql = bd.String()

	return
}
