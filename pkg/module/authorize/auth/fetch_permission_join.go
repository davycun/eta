package auth

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/expr"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
)

// NeedJoinAuth2Role
// 是否需要关联到r_auth2role来查询有权限的数据，主要是判断r_auth2role表中有没有对应的数据
func NeedJoinAuth2Role(db *gorm.DB, fromTable, authTable string, authType Type, roleIds ...string) (bool, error) {
	//添加缓存
	var (
		err  error
		perm []Auth2Role
	)
	for _, v := range roleIds {
		perm, err = LoadAuth2RoleByRoleId(db, v)
		if err != nil {
			return false, err
		}
		for _, val := range perm {
			if val.FromTable != constants.TablePermission && val.FromTable == fromTable && val.AuthTable == authTable && val.AuthType&authType == authType {
				return true, nil
			}
		}
	}

	return false, err
}

// BuildJoinAuth2RoleFilter
// 当auth2role的from_id直接关联业务数据id的时候，这个filter就是对auth2role的filter筛选
func BuildJoinAuth2RoleFilter(db *gorm.DB, fromTable, authTable, userId string, authType Type) (auth2RoleFilters []filter.Filter, err error) {
	var (
		roleIds []string
		filters = make([]filter.Filter, 0, 3)
	)
	roleIds, err = LoadUserAllRoleIds(db, userId)
	if err != nil || len(roleIds) < 1 {
		return filters, err
	}

	b, err := NeedJoinAuth2Role(db, fromTable, authTable, authType, roleIds...)
	if !b {
		return filters, err
	}

	filters = append(filters, filter.Filter{LogicalOperator: filter.And, Column: "from_table", Operator: filter.Eq, Value: fromTable})
	filters = append(filters, filter.Filter{LogicalOperator: filter.And, Column: "auth_table", Operator: filter.Eq, Value: fromTable})
	filters = append(filters, filter.Filter{LogicalOperator: filter.And, Column: "to_id", Operator: filter.IN, Value: roleIds})
	filters = append(filters, filter.Filter{
		LogicalOperator: filter.And,
		Expr: expr.Expression{
			Expr: "? & ?",
			Vars: []expr.ExpVar{
				{
					Type:  expr.VarTypeColumn,
					Value: "auth_type",
				},
				{
					Type:  expr.VarTypeValue,
					Value: int(authType),
				},
			},
		},
		Operator: filter.Eq,
		Value:    int(authType),
	})

	return filters, err
}

func BuildJoinAuth2RoleSql(db *gorm.DB, fromTable string) string {
	return fmt.Sprintf(` join "%s"."%s" on "%s"."from_id"="%s"."id"`, dorm.GetDbSchema(db), constants.TableAuth2Role, constants.TableAuth2Role, fromTable)
}
