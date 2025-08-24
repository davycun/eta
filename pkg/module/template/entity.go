package template

import (
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm/schema"
)

var (
	Ready   = "ready"
	Unready = "unready"
)

type Template struct {
	entity.BaseEntity
	Namespace string       `json:"namespace,omitempty" gorm:"column:namespace;comment:命名空间区分项目"`
	Code      string       `json:"code,omitempty" gorm:"column:code;comment:编码"`
	Title     string       `json:"title,omitempty" gorm:"column:title;comment:中文名字"`
	Alias     string       `json:"alias,omitempty" gorm:"column:alias;comment:表单别名"`
	Status    string       `json:"status,omitempty" gorm:"column:status;comment:表单状态" binding:"required,oneof=ready disable"`
	Table     entity.Table `json:"table,omitempty" gorm:"comment:表的定义;serializer:json" binding:"required"`
}

func (p Template) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableTemplate
	} else {
		return namer.TableName(constants.TableTemplate)
	}
}
func (p Template) GetTableName() string {
	if p.Table.TableName != "" {
		return p.Table.TableName
	}
	return constants.TableTemplatePrefix + p.Code
}

func (p Template) HistoryTableName() string {
	return p.GetTable().GetTableName() + constants.TableHistorySubFix
}
func (p Template) TriggerName() string {
	return constants.TableTriggerPrefix + p.GetTable().GetTableName() + constants.TableHistorySubFix
}
