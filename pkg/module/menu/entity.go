package menu

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	Function Type = "function"
	Page     Type = "page"
	Classify Type = "classify"
)

var (
	DefaultColumns = append(entity.DefaultVertexColumns, "namespace", "parent_id", "type", "name", "api", "alias", "icon", "order")
)

type Type string

type Menu struct {
	entity.BaseEntity
	Namespace string  `json:"namespace" gorm:"column:namespace;comment:命名空间;not null" binding:"required"`
	ParentId  string  `json:"parent_id,omitempty" gorm:"column:parent_id;comment:菜单父ID"`
	Type      Type    `json:"type,omitempty" gorm:"column:type;comment:菜单类型;not null"`
	Api       ApiList `json:"api,omitempty" gorm:"column:api;serializer:json;comment:菜单访问地址"`
	Name      string  `json:"name,omitempty" gorm:"column:name;comment:菜单名称;not null"`
	Alias     string  `json:"alias,omitempty" gorm:"column:alias;comment:菜单别名"`
	Icon      string  `json:"icon,omitempty" gorm:"column:icon;comment:菜单图标"`
	Order     int     `json:"order,omitempty" gorm:"column:order;comment:菜单排序"`
	Parent    *Menu   `json:"parent,omitempty" gorm:"-:all"`
	Children  []*Menu `json:"children,omitempty" gorm:"-:all"`
}

func (m Menu) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TableMenu
	}
	return namer.TableName(constants.TableMenu)
}
func (m Menu) MustColumns() []string {
	return []string{entity.IdDbName, entity.UpdatedAtDbName, "parent_id"}
}
func (m Menu) DefaultColumns() []string {
	return DefaultColumns
}

func (m *Menu) GetId() string {
	return m.ID
}
func (m *Menu) GetParentId() string {
	return m.ParentId
}
func (m *Menu) SetChildren(cd any) {
	if x, ok := cd.([]*Menu); ok {
		m.Children = x
	} else {
		logger.Error("Menu.SetChildren 参数不是[]*Menu类型")
	}
}
func (m *Menu) GetChildren() any {
	return m.Children
}
func (m *Menu) GetParentIds(db *gorm.DB) []string {
	return GetParentIds(db, m.ID)
}

type Api struct {
	Uri    string `json:"uri,omitempty"`
	Method string `json:"method,omitempty"`
}
type ApiList []Api

func (j ApiList) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return dorm.JsonGormDBDataType(db, field)
}

func (j ApiList) GormDataType() string {
	return dorm.JsonGormDataType()
}
