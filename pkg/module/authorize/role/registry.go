package role

import (
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
)

var (
	defaultRoleList = []Role{
		{
			BaseEntity: entity.BaseEntity{ID: constants.SystemAdminRoleID},
			Namespace:  "eta",
			Target:     "menu",
			Name:       "系统管理员",
			Category:   CategoryDefault,
		},
		{
			BaseEntity: entity.BaseEntity{ID: constants.DeptAdminRoleID},
			Namespace:  "eta",
			Target:     "menu",
			Name:       "部门管理员",
			Category:   CategoryPosition,
		},
		{
			BaseEntity: entity.BaseEntity{ID: constants.BasicRoleID},
			Namespace:  "eta",
			Target:     "menu",
			Name:       "用户基础功能",
			Category:   CategoryDefault,
		},
	}
)

func Registry(roleList ...Role) {
	defaultRoleList = append(defaultRoleList, roleList...)
}
