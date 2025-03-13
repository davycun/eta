package migrate

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/entity"
	"gorm.io/gorm"
)

type MigratorAfter interface {
	AfterMigrator(db *gorm.DB, c *ctx.Context) error
}

type Migrator interface {
	Migrate(dst ...interface{}) error
	MigrateOption(tableOptions ...entity.Table) error
	Schema() SchemaInterface
}
type SchemaInterface interface {
	SchemaExists(schema string) bool
	CreateSchema(schema string) error
	DefaultSchema() string
}
