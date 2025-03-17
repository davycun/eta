package auth

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/cache"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/eta/constants"
	"gorm.io/gorm"
	"slices"
)

// DelUserRoleCache 简单的清空下缓存
func DelUserRoleCache(userId ...string) {
	for _, v := range userId {
		_, _ = cache.Del(constants.RedisKey(constants.UserRoleIdsKey, v))
	}
}
func LoadUserOnlyRoleIds(db *gorm.DB, userId string) (roleIds []string, err error) {
	//获取用户的角色，角色表可能是t_role 或者t_department，将来可能支持其他
	var (
		tbl    = constants.TableUser2Role
		tmpIds []string
	)

	err = dorm.Table(db, tbl).
		Select(dorm.GetDbColumn(db, tbl, entity.ToIdDbName)).
		Where(fmt.Sprintf(`%s = ?`, dorm.GetDbColumn(db, tbl, entity.FromIdDbName)), userId).
		Find(&tmpIds).Error
	return tmpIds, err
}
func LoadUserRoleIdsByParentId(db *gorm.DB, userId string, parentId string) (roleIds []string, err error) {

	var (
		scm              = dorm.GetDbSchema(db)
		dbType           = dorm.GetDbType(db)
		idsTmpTableAlias = "allIds"
		recursiveAlias   = "recur"
		listSql          = ""
	)
	cte := builder.NewCteSqlBuilder(dbType, "", recursiveAlias)

	//获取部门的ID和ParentId
	deptBd := builder.NewSqlBuilder(dbType, scm, constants.TableDept).AddColumn(entity.IdDbName, "parent_id")
	//获取角色表的Id和ParentId
	roleBd := builder.NewSqlBuilder(dbType, scm, constants.TableRole).AddColumn(entity.IdDbName, "parent_id")
	//取合集
	deptBd.UnionAll(roleBd)
	cte.With(idsTmpTableAlias, deptBd)

	//通过递归sql获取所有的子数据
	recursiveBd := builder.NewRecursiveSqlBuilder(dbType, "", idsTmpTableAlias)
	recursiveBd.AddRecursiveFilter(filter.Filter{
		LogicalOperator: filter.And,
		Column:          entity.IdDbName,
		Operator:        filter.Eq,
		Value:           parentId,
	})
	recursiveBd.AddColumn(entity.IdDbName)
	cte.With(recursiveAlias, recursiveBd)

	cte.AddColumn(entity.IdDbName)
	listSql, _, err = cte.Build()
	if err != nil {
		return
	}

	err = dorm.RawFetch(listSql, db, &roleIds)
	return
}

// LoadUserRoleIdsByDeptId
// 找出用户所有的角色ID，同时去除掉掉角色与用户同时与部门有关系的角色
func LoadUserRoleIdsByDeptId(db *gorm.DB, userId string, deptId string) (roleIds []string, err error) {

	var (
		scm              = dorm.GetDbSchema(db)
		dbType           = dorm.GetDbType(db)
		userRoleIds      []string //只包含user2role中to_id
		user2DeptRoleIds []string //用户所在部门的所有部门ID
	)

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			//只包含roleIds
			userRoleIds, err = LoadUserOnlyRoleIds(db, userId)
			return err
		}).
		Call(func(cl *caller.Caller) error {
			type TmpUser2Dept struct {
				entity.BaseEdgeEntity
				RoleIds ctype.StringArray `json:"role_ids,omitempty" gorm:"column:role_ids;comment:角色IDS"`
			}
			var u2dList []TmpUser2Dept
			//找出用户所有的所在部门，除了当前部门
			bd := builder.NewSqlBuilder(dbType, scm, constants.TableUser2Dept).
				AddFilter(filter.Filter{LogicalOperator: filter.And, Column: entity.FromIdDbName, Operator: filter.Eq, Value: userId}).
				AddColumn(entity.IdDbName, entity.FromIdDbName, entity.ToIdDbName, "role_ids")
			listSql, _, err1 := bd.Build()
			if err1 != nil {
				return err1
			}
			err = dorm.RawFetch(listSql, db, &u2dList)
			for _, v := range u2dList {
				//把用户当前部门的角色放到角色列表中
				if v.ToId == deptId {
					userRoleIds = utils.Merge(userRoleIds, v.ToId)
				} else {
					//需要被最终排除的角色ID，因为这些角色ID当用户切换到对应的部门才生效
					user2DeptRoleIds = utils.Merge(user2DeptRoleIds, v.RoleIds.Data...)
				}
			}
			return err
		}).
		Call(func(cl *caller.Caller) error {
			//找出需要被过滤的所有角色：这些角色被同一个用户在其他部门有角色关联
			roleIds = slices.DeleteFunc(userRoleIds, func(s string) bool {
				if utils.ContainAny(user2DeptRoleIds, s) {
					return true
				}
				return false
			})
			return nil
		}).Err
	return
}

// LoadUserAllRoleIds 包括用户自己的ID，部门ID，角色ID
func LoadUserAllRoleIds(db *gorm.DB, userId string) (roleIds []string, err error) {

	//获取用户的角色，角色表可能是t_role 或者t_department，将来可能支持其他
	var (
		user2RoleRelationTables = []string{constants.TableUser2Role, constants.TableUser2Dept}
	)
	defer func() {
		roleIds = append(roleIds, userId)
	}()

	b, err := cache.Get(constants.RedisKey(constants.UserRoleIdsKey, userId), &roleIds)
	if b || err != nil {
		return
	}

	for _, v := range user2RoleRelationTables {
		var tmpIds []string
		err = dorm.Table(db, v).
			Select(dorm.GetDbColumn(db, v, entity.ToIdDbName)).
			Where(fmt.Sprintf(`%s = ?`, dorm.GetDbColumn(db, v, entity.FromIdDbName)), userId).
			Find(&tmpIds).Error
		if err != nil {
			return
		}
		roleIds = append(roleIds, tmpIds...)
	}
	err = cache.Set(constants.RedisKey(constants.UserRoleIdsKey, userId), roleIds)
	//这里不管是否为空都得存redis，否则会一直查数据库
	return roleIds, err
}

// IsSystemAdmin 是否是系统管理员
func IsSystemAdmin(db *gorm.DB, userID string) (isSystemAdmin bool) {
	var count int64

	// 查询该用户是不是系统管理员的角色
	err := dorm.Table(db, constants.TableUser2Role).
		Select(`count(*)`).
		Where(
			fmt.Sprintf(
				`%s = ? and %s = ?`,
				dorm.GetDbColumn(db, "", entity.FromIdDbName),
				dorm.GetDbColumn(db, "", entity.ToIdDbName),
			),
			userID,
			constants.SystemAdminRoleID,
		).
		Find(&count).Error

	if err != nil {
		logger.Warnf("IsSystemAdmin error: %v", err)
		return
	}

	return count > 0
}
