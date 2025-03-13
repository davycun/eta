package optlog

import (
	"fmt"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/core/builder"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/duke-git/lancet/v2/slice"
)

func init() {
	hook.AddAuthCallback(constants.TableOperateLog, authRead)
}

func authRead(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	//操作记录权限：
	//1、超管或系统管理员这个角色看所有
	//2、部门管理员看自己和所有子部门的数据，注意这里是限制部门
	//3、普通用户只能看自己的

	uId := cfg.Ctx.GetContextUserId()
	//如果是管理员就不加任何限制
	if cfg.Ctx.GetContextIsManager() {
		return nil
	}
	// 查询该用户是不是系统管理员的角色
	if auth.IsSystemAdmin(cfg.Ctx.GetAppGorm(), uId) {
		return nil
	}

	u2d, err := dept.LoadUser2DeptByUserId(cfg.Ctx, uId)
	if err != nil {
		return err
	}

	deptIds := slice.FilterMap(u2d, func(i int, v dept.RelationDept) (string, bool) { return v.ToId, v.IsManager })

	if len(deptIds) <= 0 {
		//只能看到自己
		cfg.Param.Filters = []filter.Filter{{
			LogicalOperator: filter.And,
			Column:          "opt_user_id",
			Operator:        filter.Eq,
			Value:           uId,
			Filters: []filter.Filter{{
				LogicalOperator: filter.And,
				Filters:         cfg.Param.Filters,
			}},
		}}
		return err
	}

	//部门管理员
	// ES

	if cfg.RetrieveEnableEs() {
		cfg.Param.Filters = []filter.Filter{{
			LogicalOperator: filter.And,
			Column:          "parent_dept_ids",
			Operator:        filter.IN,
			Value:           deptIds,
			Filters: []filter.Filter{{
				LogicalOperator: filter.And,
				Filters:         cfg.Param.Filters,
			}},
		}}
		return nil
	}
	// DB
	dbType := dorm.GetDbType(cfg.Ctx.GetAppGorm())
	scm := dorm.GetDbSchema(cfg.Ctx.GetAppGorm())
	listSql, _, err := builder.NewRecursiveSqlBuilder(dbType, scm, constants.TableDept).
		AddRecursiveFilter(filter.Filter{Column: entity.IdDbName, Operator: filter.IN, Value: deptIds}).AddColumn(entity.IdDbName).
		Build()
	cfg.CurDB = cfg.CurDB.Where(fmt.Sprintf(`%s in (%s)`, dorm.Quote(dbType, constants.TableOperateLog, "opt_dept_id"), listSql))
	return err
}
