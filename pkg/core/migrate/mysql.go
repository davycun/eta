package migrate

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
)

type mysqlSchema struct {
	orm *gorm.DB
}

func (s mysqlSchema) SchemaExists(schema string) bool {
	var total int64
	//err := s.orm.Raw("select count(*) from information_schema.SCHEMATA where SCHEMA_NAME=?", schema).Row().Scan(&total)
	err := s.orm.Table(`information_schema.SCHEMATA`).Where(`SCHEMA_NAME = ?`, schema).Count(&total).Error
	if err != nil {
		logger.Errorf("schemaExists error %s", err.Error())
		return false
	}
	return total > 0
}
func (s mysqlSchema) CreateSchema(scm string) error {
	if !s.SchemaExists(scm) {
		err := s.orm.Exec("create schema if not exists " + scm).Error
		if err != nil {
			return err
		}
		//自己创建的应该是有权限的，所以不需要再次调用者这个
		err = s.orm.Exec("grant all privileges on " + scm + ".* to current_user").Error
		return err
	}
	return nil
}
func (s mysqlSchema) DefaultSchema() string {
	return dorm.GetDbSchema(s.orm)
}
