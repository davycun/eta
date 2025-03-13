package auth

import (
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm/schema"
)

const (
	PermissionForAll                = "all"
	PermissionFilter PermissionType = "filter"
)

type PermissionType string

var (
	DefaultColumnsPermission = append(entity.DefaultVertexColumns, "namespace", "name", "schema", "table", "filters")
)

type Permission struct {
	entity.BaseEntity
	Namespace        string         `json:"namespace,omitempty" gorm:"column:namespace;comment:命名空间" binding:"required"` //主要用于区分不同的项目或者定制化
	Name             string         `json:"name,omitempty" gorm:"column:name;comment:权限名字" binding:"required"`           //权限的名字
	TbName           string         `json:"tb_name,omitempty" gorm:"column:tb_name;comment:对应的table的名字"`                 //如果TbName=ALL表示针对所有表，比如这个filter是creator_id的filter
	Type             PermissionType `json:"type,omitempty" gorm:"column:type;type:varchar;comment:权限类型"`                 //预留字段
	Filters          filter.Filters `json:"filters,omitempty" gorm:"serializer:json;comment:数据的过滤条件"`
	RecursiveFilters filter.Filters `json:"recursive_filters,omitempty" gorm:"serializer:json;comment:数据的递归过滤条件"`
}

func (p Permission) TableName(namer schema.Namer) string {
	if namer == nil {
		return constants.TablePermission
	}
	return namer.TableName(constants.TablePermission)
}
