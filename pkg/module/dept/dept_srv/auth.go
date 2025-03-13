package dept_srv

import (
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/dept"
)

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

	if cfg.UseParamAuth && cfg.Param.DisablePermFilter {
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
