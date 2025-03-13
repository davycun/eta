package migrate

import (
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
)

type pgSchema struct {
	orm *gorm.DB
}

func (s pgSchema) SchemaExists(schema string) bool {
	var total int64
	err := s.orm.Table(`pg_catalog.pg_namespace`).Where(`nspname = ?`, schema).Count(&total).Error
	if err != nil {
		logger.Errorf("schemaExists error %s", err.Error())
		return false
	}
	return total > 0
}
func (s pgSchema) CreateSchema(schema string) error {
	return s.orm.Exec("create schema if not exists " + schema + " authorization current_user").Error
}

func (s pgSchema) DefaultSchema() string {
	return "public"
}
