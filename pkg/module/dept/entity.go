package dept

import (
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DefaultColumns = append(entity.DefaultVertexColumns, "name", "parent_id", "display_order", "top_level")
)

type Department struct {
	entity.BaseEntity
	Namespace    string        `json:"namespace,omitempty" gorm:"column:namespace;comment:命名空间针对,同一个app下不同定制化" binding:"required"`
	Name         string        `json:"name,omitempty" gorm:"column:name;comment:部门名称" binding:"required"`
	ParentId     string        `json:"parent_id,omitempty" gorm:"type:varchar;column:parent_id;comment:父部门ID"`
	Status       string        `json:"status,omitempty" gorm:"column:status;comment:状态"`
	DisplayOrder int           `json:"display_order,omitempty" gorm:"column:display_order;comment:排序"`
	Children     []*Department `json:"children,omitempty" gorm:"-:all"`
	Parent       *Department   `json:"parent,omitempty" gorm:"-:all"`

	//市民的定制，先放下
	TopLevel int `json:"top_level,omitempty" gorm:"column:top_level;comment:地址展示的level层级"`
}

func (d *Department) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableDept
	}
	return namer.TableName(constants.TableDept)
}

func (d Department) MustColumns() []string {
	return []string{entity.IdDbName, entity.UpdatedAtDbName, "parent_id"}
}
func (d Department) DefaultColumns() []string {
	return DefaultColumns
}

func (d *Department) GetId() string {
	return d.ID
}
func (d *Department) GetParentId() string {
	return d.ParentId
}
func (d *Department) SetChildren(cd any) {
	if x, ok := cd.([]*Department); ok {
		d.Children = x
	} else {
		logger.Error("Department.SetChildren 参数不是[]*Department类型")
	}
}
func (d *Department) GetChildren() any {
	return d.Children
}
func (d *Department) GetParentIds(db *gorm.DB) []string {
	return []string{d.ParentId}
}
