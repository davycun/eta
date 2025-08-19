package user2role

import (
	"errors"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/authorize/role"
	"github.com/duke-git/lancet/v2/slice"
)

func modifyCallback(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.AfterDelete(cfg, pos, deleteUser2Role)
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterModify(cfg, pos, func(cfg *hook.SrvConfig, oldValues []User2Role, newValues []User2Role) error {
				if cfg.Method == iface.MethodCreate {
					cleanCache(newValues)
				} else {
					cleanCache(oldValues)
				}
				return nil
			})
		}).Err
	return err
}

// 新增修改和删除角色的时候，对用户与角色的缓存进行清理
func cleanCache(dt []User2Role) {

	if len(dt) < 1 {
		return
	}
	uIds := make([]string, 0, len(dt))
	for _, v := range dt {
		uIds = append(uIds, v.FromId)
	}
	auth.DelUserRoleCache(uIds...)
}

// 删除角色的时候同时删除用户与角色的关系
func deleteUser2Role(cfg *hook.SrvConfig, dt []role.Role) error {
	if len(dt) < 1 {
		return nil
	}
	var (
		u2rArgs     dto.Param
		u2rRes      dto.Result
		u2rSvc, err = service.NewService(constants.TableUser2Role, cfg.Ctx, cfg.TxDB)
		batchSize   = 1000
	)
	if err != nil {
		return err
	}
	for _, roles := range slice.Chunk(dt, batchSize) {
		u2rArgs.Filters = []filter.Filter{
			{
				LogicalOperator: filter.And,
				Column:          "to_id",
				Operator:        filter.IN,
				Value:           slice.Map(roles, func(i int, item role.Role) string { return item.ID }),
			},
		}
		u2rArgs.Data = &User2Role{}
		err = u2rSvc.DeleteByFilters(&u2rArgs, &u2rRes)
		if err != nil && !errors.Is(err, errs.NoRecordAffected) {
			return err
		}
	}
	return nil
}
