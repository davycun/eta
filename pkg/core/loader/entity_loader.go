package loader

import (
	"errors"
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"gorm.io/gorm"
)

type EntityLoader struct {
	Err     error
	Schema  string
	DB      *gorm.DB
	DbType  dorm.DbType
	Config  EntityLoaderConfig
	columns []string
}

type EntityLoaderConfig struct {
	Ids                  []string //where条件字段的值
	DefaultEntityColumns []string // select 哪些字段
	IdColumn             string   //where条件字段
	TableName            string   // entity对应哪张表
}

func NewEntityLoader(db *gorm.DB, config EntityLoaderConfig) *EntityLoader {
	l := &EntityLoader{
		DB:     db,
		DbType: dorm.GetDbType(db),
		Schema: dorm.GetDbSchema(db),
		Config: config,
	}
	return l
}

func (l *EntityLoader) SetTableName(tableName string) *EntityLoader {
	l.Config.TableName = tableName
	return l
}

func (l *EntityLoader) AddColumns(col ...string) *EntityLoader {
	l.columns = append(l.columns, col...)
	return l
}
func (l *EntityLoader) AddId(id ...string) *EntityLoader {
	if l.Config.Ids == nil {
		l.Config.Ids = make([]string, 0, 10)
	}
	l.Config.Ids = append(l.Config.Ids, id...)
	return l
}
func (l *EntityLoader) check() *EntityLoader {
	if l.Config.IdColumn == "" {
		l.Config.IdColumn = entity.IdDbName
	}

	if l.Config.TableName == "" {
		l.Err = errors.New("tableName is empty")
	}

	if len(l.Config.Ids) < 1 {
		l.Err = NoNeedLoadError
	}
	return l
}
func (l *EntityLoader) resolveColumns() *EntityLoader {
	l.columns = utils.Merge(l.columns, l.Config.DefaultEntityColumns...)
	if len(l.columns) < 1 {
		l.columns = append(l.columns, "*")
	}
	return l
}

func (l *EntityLoader) Load(rs any) error {

	if l.check().resolveColumns().Err != nil {
		if errors.Is(l.Err, NoNeedLoadError) {
			logger.Errorf("%s, for %s", NoNeedLoadError, l.Config.TableName)
			l.Err = nil
		}
		return l.Err
	}

	var (
		dbType    = dorm.GetDbType(l.DB)
		_, tbName = dorm.Quote(dbType, l.Schema), dorm.Quote(dbType, l.Config.TableName)
		scmTbName = fmt.Sprintf("%s.%s", l.Schema, l.Config.TableName)
		cols      = dorm.JoinColumns(dbType, l.Config.TableName, l.columns)
		idCol     = dorm.Quote(dbType, l.Config.IdColumn)
	)

	if len(l.Config.Ids) == 1 {
		l.Err = dorm.Table(l.DB, l.Config.TableName).Select(cols).
			Where(fmt.Sprintf(`%s = ?`, idCol), l.Config.Ids[0]).Find(rs).Error
		return l.Err
	}
	if len(l.Config.Ids) < 6 {
		l.Err = dorm.Table(l.DB, l.Config.TableName).Select(cols).
			Where(fmt.Sprintf(`%s in ?`, idCol), l.Config.Ids).Find(rs).Error
		return l.Err
	}

	//raw sql 需要自己包
	scmTbName = fmt.Sprintf("%s.%s", dorm.Quote(dbType, l.Schema), dorm.Quote(dbType, l.Config.TableName))
	rSql := builder.BuildValueToTableSql(l.DbType, true, l.Config.Ids...)
	sq := fmt.Sprintf(`with r as (%s) select %s from r, %s where r.%s=%s.%s `,
		rSql, cols, scmTbName, dorm.Quote(dbType, "id"), tbName, idCol)

	l.Err = dorm.RawFetch(sq, l.DB, rs)
	return l.Err
}
