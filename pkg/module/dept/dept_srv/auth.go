package dept_srv

import (
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/dept"
)

// BuildListSql
// 该数据权限只在后台管理中启用:
// 1. 超级管理员和系统管理员角色可以在后台管理中看到所有的部门和所有的用户（包括虚拟用户）
// 2. 部门管理员可以在后台管理中看到他有部门管理员角色的所在部门及其子部门和部门及子部门的用户。（不包括虚拟用户）
// 3. 一般人如果被分配了后台管理的功能权限，那么在后台管理中只能看到他自己，看不到部门列表。
// disable_perm_filter 控制是否开启权限过滤

func AuthRetrieve(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	uId := cfg.Ctx.GetContextUserId()

	//如果是管理员就不加任何限制
	if cfg.Ctx.GetContextIsManager() {
		return nil
	}
	// 查询该用户是不是系统管理员的角色
	if auth.IsSystemAdmin(cfg.Ctx.GetAppGorm(), uId) {
		return nil
	}

	if cfg.UseParamAuth() && cfg.Param.DisablePermFilter {
		return nil
	}

	var (
		u2d, err = dept.LoadUser2DeptByUserId(cfg.Ctx, uId)
		deptIds  = make([]string, 0, 3)
	)
	if err != nil {
		return err
	}

	for i, v := range u2d {
		if v.IsManager {
			deptIds = append(deptIds, u2d[i].ToId)
		}
	}

	if len(deptIds) < 1 {
		return errs.NoPermissionNoErr
	}

	flt := filter.Filter{
		LogicalOperator: filter.And,
		Column:          entity.IdDbName,
		Operator:        filter.IN,
		Value:           deptIds,
	}

	cfg.Param.AuthRecursiveFilters = append(cfg.Param.AuthRecursiveFilters, flt)
	return nil
}
