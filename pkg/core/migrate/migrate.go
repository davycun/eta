package migrate

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/entity"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

func NewSchemaInterface(db *gorm.DB) SchemaInterface {
	switch db.Name() {
	case dorm.Mysql.String():
		return &mysqlSchema{
			orm: db,
		}
	case dorm.PostgreSQL.String():
		return &pgSchema{
			orm: db,
		}
	case dorm.DaMeng.String():
		return &dmSchema{
			orm: db,
		}
	case dorm.Doris.String():
		return &dorisSchema{
			orm: db,
		}
	default:
		return nil
	}
}

func NewMigrator(db *gorm.DB, c *ctx.Context) Migrator {

	dbType := dorm.GetDbType(db)
	switch dbType {
	case dorm.Mysql:
		mg := &baseMigrator{
			dialect: dbType,
			orm:     db,
			si: &mysqlSchema{
				orm: db,
			},
			c: c,
		}
		return mg
	case dorm.Doris:
		mg := &baseMigrator{
			dialect: dbType,
			orm:     db,
			si: &dorisSchema{
				orm: db,
			},
			c: c,
		}
		return mg
	case dorm.PostgreSQL:
		return &baseMigrator{
			dialect: dbType,
			orm:     db,
			si: &pgSchema{
				orm: db,
			},
			c: c,
		}
	case dorm.DaMeng:
		return &baseMigrator{
			dialect: dbType,
			orm:     db,
			si: &dmSchema{
				orm: db,
			},
			c: c,
		}
	default:
		return nil
	}
}

type baseMigrator struct {
	dialect dorm.DbType
	orm     *gorm.DB
	err     error
	si      SchemaInterface
	c       *ctx.Context
}

func (m *baseMigrator) Migrate(dst ...interface{}) error {
	toList := make([]entity.Table, 0, len(dst))
	for _, v := range dst {
		tp := reflect.TypeOf(v)
		if tp.Kind() == reflect.Pointer {
			tp = tp.Elem()
		}
		toList = append(toList, entity.Table{EntityType: tp})
	}
	return m.MigrateOption(toList...)
}

// MigrateOption
// 如果是事务处理，那么当有字段更改的时候，pg会报错"ERROR: cached plan must not change result type"
// 原因是pg会的Migrator实现在ColumnTypes方法中会调用表的select * from xxx_target limit 1
// 但是又因为会去对当前xxx_target的字段类型进行修改(alter)所以导致问题。如果不在事务中就没问题。

func (m *baseMigrator) MigrateOption(options ...entity.Table) error {

	var (
		schemas    = make(map[string]string, 10)
		dbType     = dorm.GetDbType(m.orm)
		scm        = dorm.GetDbSchema(m.orm)
		entityList = make([]entity.Table, 0, len(options))
	)

	schemaName := dorm.GetDbSchema(m.orm)
	if schemaName != "" {
		schemas[schemaName] = schemaName
	}
	for _, v := range schemas {
		err := m.si.CreateSchema(v)
		if err != nil {
			return err
		}
	}
	for _, tb := range options {
		tbOption := ""
		switch dbType {
		case dorm.Mysql:
			tbOption = tb.Options[dorm.Mysql]
			if tbOption == "" {
				// https://dev.mysql.com/doc/refman/8.0/en/charset-applications.html
				tbOption = "ENGINE=InnoDB ROW_FORMAT=DYNAMIC"
			}
		case dorm.PostgreSQL:
			tbOption = tb.Options[dorm.PostgreSQL]
		case dorm.DaMeng:
			tbOption = tb.Options[dorm.DaMeng]
		case dorm.Doris:
			tbOption = tb.Options[dorm.Doris]
		}
		//if len(tb.EnableDbType) > 0 {
		//	mgFlag := false
		//	for _, v := range tb.EnableDbType {
		//		if v == dbType {
		//			mgFlag = true
		//		}
		//	}
		//	if !mgFlag {
		//		continue
		//	}
		//}
		if tbOption != "" {
			if strings.Contains(tbOption, "%s") {
				tbOption = fmt.Sprintf(tbOption, scm)
			}
			m.err = m.orm.Set("gorm:table_options", tbOption).AutoMigrate(tb.NewEntityPointer())
		} else {
			m.err = m.orm.AutoMigrate(tb.NewEntityPointer())
		}
		if m.err != nil {
			return m.err
		}
		entityList = append(entityList, tb)
	}
	//只有表都创建完了才能做after
	for _, tb := range entityList {
		mc := NewMigConfig(m.c, m.orm, tb)
		m.err = mc.after()
		if m.err != nil {
			return m.err
		}
	}
	return m.err
}
func (m *baseMigrator) Schema() SchemaInterface {
	return m.si
}
