package template

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

// GetTable
// 获取实际模版数据的表名
func (p *Template) GetTable() *entity.Table {
	if p.Table.TableName == "" && p.Code != "" {
		p.Table.TableName = constants.TableTemplatePrefix + p.Code
	}
	return &p.Table
}

func getGormType(db *gorm.DB, tp string) (string, error) {
	if tp == ctype.TypeJsonName {
		return "serializer:json", nil
	}
	rs, err := ctype.GetDbTypeName(db, tp)
	return "type:" + rs, err
}

func getBinding(bind string) string {
	if bind != "" {
		return fmt.Sprintf(`binding:"%s"`, bind)
	}
	return ""
}
