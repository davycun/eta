package builder

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"golang.org/x/exp/slices"
	"strings"
)

const (
	Union          = "union"
	UnionAll       = "union all"
	UnionIntersect = "intersect"
	UnionExcept    = "except"
)

type whereExpr struct {
	operator string
	expr     string
}

type SqlBuilder struct {
	Err          error
	schema       string
	tableName    string
	newTableName string //可以通过SetTableName来改变创建时候传入的TableName，使用tableName的时候必须用GetTableName
	dbType       dorm.DbType

	distinct    string
	countColumn string
	joins       []string
	//columns     map[string][]expr.ExpColumn // tableName -> columns
	//filters     map[string][]filter.Filter  //tableName -> filters
	columns          []tableColumns // tableName -> columns
	partitionColumns []dorm.PartitionColumn
	filters          []tableFilters //tableName -> filters
	unionList        []unionBuilder

	mustPage bool
	offSet   int
	limit    int
	orderBy  []dorm.OrderBy

	where []whereExpr // where string
}

func NewSqlBuilder(dbType dorm.DbType, schemaName, tableName string) *SqlBuilder {
	b := &SqlBuilder{
		schema:      schemaName,
		dbType:      dbType,
		tableName:   tableName,
		countColumn: "*",
	}
	return b
}

func (b *SqlBuilder) GetTableName() string {
	if b.newTableName == "" {
		return b.tableName
	}
	return b.newTableName
}
func (b *SqlBuilder) SetTableName(tableName string) *SqlBuilder {
	b.newTableName = tableName
	//b.tableName = tableName //原始的表名不能更改
	return b
}
func (b *SqlBuilder) Union(bd Builder) *SqlBuilder {
	b.unionList = append(b.unionList, unionBuilder{union: Union, builder: bd})
	return b
}
func (b *SqlBuilder) UnionAll(bd Builder) *SqlBuilder {
	b.unionList = append(b.unionList, unionBuilder{union: UnionAll, builder: bd})
	return b
}
func (b *SqlBuilder) UnionIntersect(bd Builder) *SqlBuilder {
	b.unionList = append(b.unionList, unionBuilder{union: UnionIntersect, builder: bd})
	return b
}
func (b *SqlBuilder) UnionExcept(bd Builder) *SqlBuilder {
	b.unionList = append(b.unionList, unionBuilder{union: UnionExcept, builder: bd})
	return b
}
func (b *SqlBuilder) AddColumn(col ...string) *SqlBuilder {
	return b.AddTableColumn(b.GetTableName(), col...)
}
func (b *SqlBuilder) AddPartitionColumn(col ...dorm.PartitionColumn) *SqlBuilder {
	b.partitionColumns = append(b.partitionColumns, col...)
	return b
}
func (b *SqlBuilder) AddColumnPrefixAlias(prefix string, col ...string) *SqlBuilder {
	return b.AddTableColumnPrefixAlias(b.GetTableName(), prefix, col...)
}
func (b *SqlBuilder) AddExprColumn(col ...expr.ExpColumn) *SqlBuilder {
	return b.AddTableExprColumn(b.GetTableName(), col...)
}
func (b *SqlBuilder) AddTableColumn(tableName string, col ...string) *SqlBuilder {

	if tableName == "" || len(col) < 1 {
		return b
	}
	ecs := make([]expr.ExpColumn, 0, len(col))
	for _, v := range col {
		ec := expr.ExpColumn{}
		ec.Expr = v
		ecs = append(ecs, ec)
	}
	return b.AddTableExprColumn(tableName, ecs...)
}
func (b *SqlBuilder) AddTableColumnPrefixAlias(tableName string, prefix string, col ...string) *SqlBuilder {

	if tableName == "" || len(col) < 1 {
		return b
	}
	ecs := make([]expr.ExpColumn, 0, len(col))
	for _, v := range col {
		ec := expr.NewAliasColumn(v, prefix+v)
		ecs = append(ecs, ec)
	}
	return b.AddTableExprColumn(tableName, ecs...)
}
func (b *SqlBuilder) AddTableExprColumn(tableName string, col ...expr.ExpColumn) *SqlBuilder {
	//允许TableName为空
	if len(col) < 1 {
		return b
	}

	flag := false
	for i, v := range b.columns {
		if v.tableName == tableName {
			flag = true
			b.columns[i].columns = append(b.columns[i].columns, col...)
		}
	}
	if !flag {
		b.columns = append(b.columns, tableColumns{tableName: tableName, columns: col})
	}
	return b
}
func (b *SqlBuilder) AddFilter(filters ...filter.Filter) *SqlBuilder {
	return b.AddTableFilter(b.GetTableName(), filters...)
}
func (b *SqlBuilder) AddTableFilter(tableName string, filters ...filter.Filter) *SqlBuilder {
	if tableName == "" || len(filters) < 1 {
		return b
	}
	flag := false
	for i, v := range b.filters {
		if v.tableName == tableName {
			flag = true
			b.filters[i].filters = append(b.filters[i].filters, filters...)
		}
	}
	if !flag {
		b.filters = append(b.filters, tableFilters{tableName: tableName, filters: filters})
	}
	return b
}
func (b *SqlBuilder) AddWhere(operator string, wh ...string) *SqlBuilder {
	if operator == "" {
		operator = "and"
	}
	for _, v := range wh {
		if v == "" {
			continue
		}
		b.where = append(b.where, whereExpr{operator: operator, expr: v})
	}
	return b
}

func (b *SqlBuilder) SetMustPage(mustPage bool) *SqlBuilder {
	b.mustPage = mustPage
	return b
}
func (b *SqlBuilder) SetDistinct(distinct bool) *SqlBuilder {
	if distinct {
		b.distinct = "distinct"
	} else {
		b.distinct = ""
	}
	return b
}
func (b *SqlBuilder) SetCountColumn(column string) *SqlBuilder {
	if column != "" {
		b.countColumn = column
	}
	return b
}
func (b *SqlBuilder) Joins(join ...string) *SqlBuilder {
	for _, v := range join {
		if v != "" {
			b.joins = append(b.joins, v)
		}
	}
	return b
}
func (b *SqlBuilder) Join(scm, tb1Name, tb1Col, tb2Name, tb2Col string) *SqlBuilder {
	return b.join("join", scm, tb1Name, tb1Col, tb2Name, tb2Col)
}
func (b *SqlBuilder) LeftJoin(scm, tb1Name, tb1Col, tb2Name, tb2Col string) *SqlBuilder {

	return b.join("left join", scm, tb1Name, tb1Col, tb2Name, tb2Col)
}
func (b *SqlBuilder) join(join, scm, tb1Name, tb1Col, tb2Name, tb2Col string) *SqlBuilder {

	b.joins = append(b.joins, fmt.Sprintf("%s %s on %s = %s", join,
		dorm.Quote(b.dbType, scm, tb1Name),
		dorm.Quote(b.dbType, tb1Name, tb1Col),
		dorm.Quote(b.dbType, tb2Name, tb2Col)))

	return b
}

// JoinWithFilter 如果filters不为空才进行join
func (b *SqlBuilder) JoinWithFilter(tableName, join string, filters ...filter.Filter) *SqlBuilder {
	if join == "" || len(filters) < 1 {
		return b
	}
	if tableName == "" {
		tableName = b.GetTableName()
	}
	return b.Joins(join).AddTableFilter(tableName, filters...)
}

func (b *SqlBuilder) AddOrderBy(orderBy ...dorm.OrderBy) *SqlBuilder {
	if len(orderBy) < 1 {
		return b
	}
	b.orderBy = append(b.orderBy, orderBy...)
	return b
}
func (b *SqlBuilder) Offset(offSet int) *SqlBuilder {
	b.offSet = offSet
	return b
}
func (b *SqlBuilder) Limit(limit int) *SqlBuilder {
	b.limit = limit
	return b
}

func (b *SqlBuilder) check() *SqlBuilder {
	if b.GetTableName() == "" {
		b.Err = errors.New("tableName is empty")
	}
	if len(b.columns) < 1 {
		b.AddColumn("*")
	}
	if b.dbType == "" {
		b.dbType = dorm.DaMeng
	}
	if b.countColumn == "" {
		b.countColumn = "*"
	}

	return b
}
func (b *SqlBuilder) resolveColStr() string {
	var (
		cols = make([]string, 0, 10)
	)

	for _, v := range b.columns {
		tbName := v.tableName
		if b.tableName == tbName && b.newTableName != "" {
			tbName = b.newTableName
		}
		cl, er := expr.JoinExprColumn(b.dbType, tbName, v.columns...)
		if er != nil {
			b.Err = er
			return ""
		}
		cols = append(cols, cl)
	}

	if len(b.partitionColumns) > 0 {
		cols = append(cols, dorm.ResolvePartitionColumn(b.dbType, b.GetTableName(), b.partitionColumns))
	}

	return strings.Join(cols, ",")
}
func (b *SqlBuilder) resolveCountCol() string {
	cc := b.countColumn
	if cc != "*" {
		cc = dorm.Quote(b.dbType, b.schema, cc)
	}
	return fmt.Sprintf(`count(%s %s)`, b.distinct, cc)
}
func (b *SqlBuilder) resolveFilterWhere(withWhere bool) (wh string) {

	bd := strings.Builder{}
	whs := make([]string, 0, len(b.filters))

	for _, v := range b.filters {
		tbName := v.tableName
		if b.tableName == tbName && b.newTableName != "" {
			tbName = b.newTableName
		}
		tmp := filter.ResolveWhereTable(tbName, v.filters, b.dbType)
		if tmp != "" {
			whs = append(whs, tmp)
		}
	}
	if len(whs) > 0 {
		bd.WriteString(strings.Join(whs, " and "))
	}

	str := bd.String()
	for i, v := range b.where {
		if v.operator == "" {
			v.operator = "and"
		}
		if str != "" || i > 0 {
			bd.WriteString(fmt.Sprintf(` %s `, v.operator))
		}
		bd.WriteString(v.expr)
	}

	if withWhere && bd.String() != "" {
		return "where " + bd.String()
	}
	return bd.String()
}
func (b *SqlBuilder) resolveLimitString() string {
	if !b.mustPage && b.limit == 0 {
		return ""
	}
	return dorm.ResolveLimitString(b.offSet, b.limit)
}
func (b *SqlBuilder) resolveOrderByString() string {
	return dorm.ResolveOrderByString(b.orderBy, b.GetTableName(), "", true)
}
func (b *SqlBuilder) Clone() SqlBuilder {
	nb := SqlBuilder{}
	nb.Err = b.Err
	nb.schema = b.schema
	nb.tableName = b.tableName
	nb.newTableName = b.newTableName
	nb.dbType = b.dbType
	nb.distinct = b.distinct
	nb.countColumn = b.countColumn
	nb.joins = make([]string, 0, len(b.joins))
	for _, v := range b.joins {
		nb.joins = append(nb.joins, v)
	}
	nb.columns = slices.Clone(b.columns)
	nb.filters = slices.Clone(b.filters)
	nb.unionList = slices.Clone(b.unionList)
	nb.mustPage = b.mustPage
	nb.offSet = b.offSet
	nb.limit = b.limit
	nb.orderBy = slices.Clone(b.orderBy)
	nb.where = slices.Clone(b.where)
	return nb
}

func (b *SqlBuilder) Build() (listSql, countSql string, err error) {

	if err = b.check().Err; err != nil {
		return
	}

	var (
		bd        = strings.Builder{}
		where     = b.resolveFilterWhere(true)
		lmt       = b.resolveLimitString()
		odb       = b.resolveOrderByString()
		joins     = strings.Join(b.joins, " ")
		scmTbName = dorm.Quote(b.dbType, b.schema, b.GetTableName())
	)

	if b.Err != nil {
		return "", "", b.Err
	}

	listSql = fmt.Sprintf(`select %s %s from %s %s %s %s %s`, b.distinct, b.resolveColStr(), scmTbName, joins, where, odb, lmt)
	countSql = fmt.Sprintf(`select %s from %s %s %s`, b.resolveCountCol(), scmTbName, joins, where)

	bd.WriteString(listSql)
	for _, v := range b.unionList {
		ls, _, err1 := v.builder.Build()
		if err1 != nil {
			return "", "", err1
		}
		bd.WriteString(fmt.Sprintf(" %s ", v.union))
		bd.WriteString("(" + ls + ")")
	}

	if len(b.unionList) > 0 {
		listSql = bd.String()
		tb := dorm.Quote(b.dbType, "r")
		countSql = fmt.Sprintf(`with %s as (%s) select count(*) from %s`, tb, listSql, tb)
	}
	return
}

type unionBuilder struct {
	union   string
	builder Builder
}

type tableFilters struct {
	tableName string
	filters   []filter.Filter
}
type tableColumns struct {
	tableName string
	columns   []expr.ExpColumn
}
