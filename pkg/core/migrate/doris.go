package migrate

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"gorm.io/gorm"
)

type dorisSchema struct {
	orm *gorm.DB
}

func (s dorisSchema) SchemaExists(schema string) bool {
	var total int64
	//err := s.orm.Table(`information_schema.SCHEMATA`).Where(`SCHEMA_NAME = ?`, schema).Count(&total).Error
	//err := s.orm.Raw(fmt).Count(&total).Error
	err := dorm.RawFetch(fmt.Sprintf("select count(*) from `information_schema`.`SCHEMATA` where `SCHEMA_NAME`='%s'", schema), s.orm, &total)
	if err != nil {
		logger.Errorf("schemaExists error %s", err.Error())
		return false
	}

	return total > 0
}
func (s dorisSchema) CreateSchema(scm string) error {
	if !s.SchemaExists(scm) {
		err := s.orm.Exec("create database if not exists " + scm).Error
		//if err != nil {
		//	return err
		//}
		//自己创建的应该是有权限的，所以不需要再次调用者这个
		//err = s.orm.Exec("grant all privileges on " + scm + ".* to current_user").Error
		return err
	}
	return nil
}
func (s dorisSchema) DefaultSchema() string {
	return dorm.GetDbSchema(s.orm)
}
