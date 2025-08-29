package loader

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"gorm.io/gorm"
)

var (
	NoNeedLoadError = errors.New("No need to load because no with is true or ids is empty ")
)

type (
	OptionFunc func(opt *Options)
	Options    struct {
		db                *gorm.DB
		schema            string
		dbType            dorm.DbType
		relationTableName string
		fromTable         string
		toTable           string
		fromColumn        string
		toColumn          string

		relationCol       []string
		entityCol         []string
		entityColPrefix   string
		DisableAuthFilter bool //是否去掉权限过滤
	}
)

func (l *Options) GetFromColumn() string {
	if l.fromColumn != "" {
		return l.fromColumn
	}
	return entity.FromIdDbName
}
func (l *Options) GetToColumn() string {
	if l.toColumn != "" {
		return l.toColumn
	}
	return entity.ToIdDbName
}

func (l *Options) SetEntityColPrefix(prefix string) *Options {
	l.entityColPrefix = prefix
	return l
}
func (l *Options) SetFromColumn(fromCol string) *Options {
	l.fromColumn = fromCol
	return l
}
func (l *Options) SetToColumn(toCol string) *Options {
	l.toColumn = toCol
	return l
}

func (l *Options) AddEntityColumns(col ...string) *Options {
	l.entityCol = utils.Merge(l.entityCol, col...)
	return l
}
func (l *Options) AddRelationColumns(col ...string) *Options {
	l.relationCol = utils.Merge(l.relationCol, col...)
	return l
}

// RelationEntityLoader
// E代表实体，R打标关联了实体的关系
// 比如：E是Address，R是RelationAddr
type RelationEntityLoader[E, R any] struct {
	Err error
	Options
}

func NewRelationEntityLoader[E, R any](db *gorm.DB, fromTable, toTable string) *RelationEntityLoader[E, R] {
	l := &RelationEntityLoader[E, R]{
		Options: Options{
			db:         db,
			dbType:     dorm.GetDbType(db),
			schema:     dorm.GetDbSchema(db),
			fromTable:  fromTable,
			toTable:    toTable,
			fromColumn: entity.FromIdDbName,
			toColumn:   entity.ToIdDbName,
		},
	}
	return l
}
func (l *RelationEntityLoader[E, R]) getEmbeddedPrefix() string {

	if l.entityColPrefix != "" {
		return l.entityColPrefix
	}
	r := new(E)
	if x, ok := any(r).(entity.Embedded); ok {
		return x.EmbeddedPrefix()
	}
	return ""
}

func (l *RelationEntityLoader[E, R]) WithOption(optionFunc ...OptionFunc) *RelationEntityLoader[E, R] {
	for _, v := range optionFunc {
		v(&l.Options)
	}
	return l
}
func (l *RelationEntityLoader[E, R]) resolveEntityCols() []string {
	if len(l.entityCol) < 1 || utils.ContainAny(l.entityCol, "*") {
		return entity.GetTableColumns(new(E))
	}
	return utils.Merge(l.entityCol, entity.IdDbName)
}
func (l *RelationEntityLoader[E, R]) resolveRelationCols() []string {
	if len(l.relationCol) < 1 || utils.ContainAny(l.entityCol, "*") {
		//这里不能直接获取R的所有字段，因为R可能是各种关系表的字段融合的struct
		ec, b := iface.GetEntityConfigByKey(l.relationTableName)
		if b {
			return entity.GetTableColumns(ec.NewEntityPointer())
		}
		return []string{entity.IdDbName, l.fromColumn, l.toColumn}
	}
	return utils.Merge(l.relationCol, l.fromColumn, l.toColumn, entity.IdDbName)
}

func (l *RelationEntityLoader[E, R]) LoadToMap(fromIds ...string) (map[string][]R, error) {

	var (
		listSql = ""
		rsMap   = make(map[string][]R)
	)
	if len(fromIds) < 1 {
		return rsMap, nil
	}

	cte := builder.NewCteSqlBuilder(l.dbType, l.schema, l.fromTable)
	cte.AddColumn(l.resolveRelationCols()...)

	//如果fromIds小于5用IN条件，如果大于5，用Value，然后join。
	//最后在基治中台千万级以上数据发现达梦多表join的时候，用in也不行，还是用value join 性能可以
	//if len(fromIds) < 5 {
	//	cte.AddFilter(filter.Filter{Column: l.fromColumn, Operator: filter.IN, Value: fromIds})
	//} else {
	//}
	vb := builder.NewValueBuilder(l.dbType, entity.IdDbName, fromIds...)
	cte.With("vb", vb)
	cte.Join("", "vb", entity.IdDbName, l.fromTable, l.fromColumn)

	//添加To表（entity表）的字段
	eCols := l.resolveEntityCols()
	eCols = utils.Merge(eCols, entity.IdDbName)
	for _, v := range eCols {
		expCol := expr.ExpColumn{}
		expCol.Expr = "?"
		expCol.Vars = []expr.ExpVar{{Type: expr.VarTypeColumn, Value: v}}
		expCol.Alias = fmt.Sprintf("%s%s", l.getEmbeddedPrefix(), v)
		cte.AddTableExprColumn(l.toTable, expCol)
	}
	cte.Join(l.schema, l.toTable, entity.IdDbName, l.fromTable, l.toColumn)

	//开始构建sql及查询数据
	listSql, _, l.Err = cte.Build()
	if l.Err != nil {
		return rsMap, l.Err
	}

	var rs []R
	l.Err = dorm.RawFetch(listSql, l.db, &rs)
	if l.Err != nil {
		return rsMap, l.Err
	}

	for _, v := range rs {
		fromId := entity.GetString(v, l.fromColumn)
		tmpRs := rsMap[fromId]
		tmpRs = append(tmpRs, v)
		rsMap[fromId] = tmpRs
	}
	return rsMap, l.Err
}
