package migrate

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
)

type dmSchema struct {
	orm *gorm.DB
}

func (s dmSchema) SchemaExists(schema string) bool {
	var total int64
	err := s.orm.Table(`SYS.SYSOBJECTS`).Where(`TYPE$ = ? and NAME = ?`, "SCH", schema).Count(&total).Error
	if err != nil {
		logger.Errorf("schemaExists error %s", err.Error())
		return false
	}
	return total > 0
}
func (s dmSchema) CreateSchema(schema string) error {
	if !s.SchemaExists(schema) {
		return s.orm.Exec(`create schema "` + schema + `" authorization ` + dorm.GetDbUser(s.orm)).Error
	}
	return nil
}

func (s dmSchema) DefaultSchema() string {
	var name string
	s.orm.Raw("SELECT SYS_CONTEXT('USERENV', 'CURRENT_SCHEMA');").Row().Scan(&name)
	return name
}
