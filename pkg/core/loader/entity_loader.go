package loader

import (
	"errors"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"gorm.io/gorm"
)

type (
	EntityLoaderConfig struct {
		ids       []string //where条件字段的值
		columns   []string // select 哪些字段
		idColumn  string   //where条件字段
		tableName string   // entity对应哪张表
	}

	EntityLoader struct {
		EntityLoaderConfig
		Err    error
		Schema string
		DB     *gorm.DB
		DbType dorm.DbType
	}

	LoadOption func(*EntityLoaderConfig)
)

func (l *EntityLoaderConfig) SetIdColumn(idColumn string) *EntityLoaderConfig {
	l.idColumn = idColumn
	return l
}
func (l *EntityLoaderConfig) SetIds(ids ...string) *EntityLoaderConfig {
	if len(ids) > 0 {
		l.ids = utils.Merge(l.ids, ids...)
	}
	return l
}
func (l *EntityLoaderConfig) SetColumns(cols ...string) *EntityLoaderConfig {
	if len(cols) > 0 {
		l.columns = utils.Merge(l.columns, cols...)
	}
	return l
}
func (l *EntityLoaderConfig) SetTableName(tbName string) *EntityLoaderConfig {
	l.tableName = tbName
	return l
}

func (l *EntityLoaderConfig) Clone() EntityLoaderConfig {
	el := EntityLoaderConfig{
		idColumn:  l.idColumn,
		tableName: l.tableName,
	}
	copy(el.ids, l.ids)
	copy(el.columns, l.columns)
	return el
}

func NewEntityLoader(db *gorm.DB, opts ...LoadOption) *EntityLoader {
	l := &EntityLoader{
		DB:     db,
		DbType: dorm.GetDbType(db),
		Schema: dorm.GetDbSchema(db),
	}
	for _, opt := range opts {
		opt(&l.EntityLoaderConfig)
	}
	return l
}

func (l *EntityLoader) AddColumns(col ...string) *EntityLoader {
	l.columns = append(l.columns, col...)
	return l
}
func (l *EntityLoader) AddId(id ...string) *EntityLoader {
	l.ids = utils.Merge(l.ids, id...)
	return l
}
func (l *EntityLoader) check() *EntityLoader {
	if l.idColumn == "" {
		l.idColumn = entity.IdDbName
	}

	if l.tableName == "" {
		l.Err = errors.New("tableName is empty")
	}

	if len(l.ids) < 1 {
		l.Err = NoNeedLoadError
	}
	return l
}
func (l *EntityLoader) resolveColumns() []string {
	if len(l.columns) < 1 {
		l.columns = append(l.columns, "*")
	} else {
		l.columns = utils.Merge(l.columns, l.idColumn)
	}
	return l.columns
}

func (l *EntityLoader) Load(rs any) error {

	if l.check().Err != nil {
		if errors.Is(l.Err, NoNeedLoadError) {
			logger.Errorf("%s, for %s", NoNeedLoadError, l.tableName)
			l.Err = nil
		}
		return l.Err
	}

	var (
		dbType    = dorm.GetDbType(l.DB)
		scm       = dorm.GetDbSchema(l.DB)
		tableName = l.tableName
		idColumn  = l.idColumn
		columns   = l.resolveColumns()
	)

	flt := filter.Filter{
		LogicalOperator: filter.And,
		Operator:        filter.IN,
		Column:          l.idColumn,
		Value:           l.ids,
	}

	cte := builder.NewCteSqlBuilder(dbType, scm, tableName)
	cte.AddColumn(columns...)
	if len(l.ids) == 1 {
		flt.Value = l.ids[0]
		flt.Operator = filter.Eq
		cte.AddFilter(flt)
	} else if len(l.ids) < 6 {
		flt.Value = l.ids
		flt.Operator = filter.IN
		cte.AddFilter(flt)
	} else {
		vs := builder.NewValueBuilder(dbType, idColumn, l.ids...)
		cte.With("ids", vs)
		cte.Join("", "ids", idColumn, tableName, idColumn)
	}

	listSql, _, err := cte.Build()

	if err != nil {
		l.Err = err
		return err
	}
	l.Err = dorm.RawFetch(listSql, l.DB, rs)
	return l.Err
}
