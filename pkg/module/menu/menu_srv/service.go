package menu_srv

import (
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/menu"
)

func init() {
	hook.AddModifyCallback(constants.TableMenu, modifyCallback)
	hook.AddRetrieveCallback(constants.TableMenu, retrieveCallback)
	sqlbd.AddSqlBuilder(constants.TableMenu, buildListSql, iface.MethodList)
}

type Service struct {
	service.DefaultService
}

func (s *Service) MyMenu(cfg *hook.SrvConfig) error {

	var (
		err     error
		c       = s.GetContext()
		listRs  = make([]menu.Menu, 0, 100)
		menuMap = make(map[string]menu.Menu)
		args    = cfg.Param
		result  = cfg.Result
	)

	if len(args.Filters) > 0 {
		err = s.GetDB().Model(&listRs).Where(filter.ResolveWhere(args.Filters, dorm.GetDbType(s.GetDB()))).Find(&listRs).Error
		if err != nil {
			return err
		}
	} else {
		menuMap, err = menu.LoadAllMenu(s.GetDB())
		if err != nil {
			return err
		}
		for _, v := range menuMap {
			listRs = append(listRs, v)
		}
	}

	// 1.0 先判断是否为超级管理员, 如果为超级管理员的话不需要判断用户部门&角色的权限
	// 1.1 如果不是超级管理员，判断use_cur_dept_auth 为true的情况下 判断是否有部门&&角色的权限
	if !s.GetContext().GetContextIsManager() {
		if args.UseCurDeptAuth {
			//ids, err1 := auth.LoadUserOnlyRoleIds(c.GetAppGorm(), c.GetContextUserId())
			ids, err1 := auth.LoadUserRoleIdsByDeptId(c.GetAppGorm(), c.GetContextUserId(), c.GetContextCurrentDeptId())
			if err1 != nil {
				return err1
			}
			ids = utils.Merge(ids, c.GetContextUserId(), c.GetContextCurrentDeptId())
			listRs, err = menu.FilterMenuByRoleIds(c, listRs, ids...)
		} else {
			listRs, err = menu.FilterMenuByUserId(c, c.GetContextUserId(), listRs)
		}
	}

	if err != nil || len(listRs) < 1 {
		return err
	}

	result.Total = int64(len(listRs))
	if args.WithTree {
		result.Data = entity.Tree(s.GetDB(), listRs)
	} else {
		if !args.IsUp {
			err = fill(cfg, listRs)
		}
		result.Data = listRs
	}
	return err
}
