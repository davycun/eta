package authorize

import (
	"errors"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/common/dorm/filter"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/loader"
	"github.com/davycun/eta/pkg/core/service"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/ws"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/davycun/eta/pkg/module/authorize/role"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/duke-git/lancet/v2/slice"
)

func modifyCallbackAuth2Role(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []auth.Auth2Role) error {
				return delCacheAuth2Role(dorm.GetDbSchema(cfg.CurDB), newValues...)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []auth.Auth2Role) error {
				return delCacheAuth2Role(dorm.GetDbSchema(cfg.CurDB), oldValues...)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterUpdate(cfg, pos, func(cfg *hook.SrvConfig, oldValues []auth.Auth2Role, newValues []auth.Auth2Role) error {
				return delCacheAuth2Role(dorm.GetDbSchema(cfg.CurDB), oldValues...)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterModify(cfg, pos, notifyMenuChanged)
		}).Err

	return err
}

func delCacheAuth2Role(scm string, dt ...auth.Auth2Role) error {
	slice.ForEach(dt, func(index int, item auth.Auth2Role) {
		auth.DelAuth2RoleCache(scm, item.ToId)
	})
	return nil
}

func afterDeleteRoleDeleteAuth2Role(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []role.Role) error {
		return deleteAuth2Role(cfg, oldValues)
	})
}

func afterDeleteDeptDeleteAuth2Role(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	return hook.AfterDelete(cfg, pos, func(config *hook.SrvConfig, oldValues []dept.Department) error {
		return deleteAuth2Role(config, oldValues)
	})
}

func deleteAuth2Role[T role.Role | dept.Department](cfg *hook.SrvConfig, dt []T) error {
	if len(dt) < 1 {
		return nil
	}
	var (
		a2rArgs   dto.Param
		a2rRes    dto.Result
		a2rSvc    = service.NewService(constants.TableAuth2Role, cfg.Ctx, cfg.TxDB)
		batchSize = 100
	)
	for _, chunk := range slice.Chunk(dt, batchSize) {
		a2rFilters := make([]filter.Filter, 0, batchSize)
		slice.ForEach(chunk, func(index int, item T) {
			var toID string
			switch any(item).(type) {
			case role.Role:
				toID = any(item).(role.Role).ID
			case dept.Department:
				toID = any(item).(dept.Department).ID
			}
			a2rFilters = append(a2rFilters, filter.Filter{
				LogicalOperator: filter.Or,
				Filters: filter.Filters{
					{
						LogicalOperator: filter.And,
						Column:          "to_id",
						Operator:        filter.Eq,
						Value:           toID,
					},
					{
						LogicalOperator: filter.And,
						Column:          "to_table",
						Operator:        filter.Eq,
						Value:           cfg.GetTableName(),
					},
				},
			})
		})
		a2rArgs.Filters = a2rFilters
		a2rArgs.Data = &auth.Auth2Role{}
		err := a2rSvc.DeleteByFilters(&a2rArgs, &a2rRes)
		if err != nil && !errors.Is(err, errs.NoRecordAffected) {
			return err
		}
	}
	return nil
}

// 菜单变更，给部门和角色相关人发送 websocket 通知
func notifyMenuChanged(cfg *hook.SrvConfig, oldValues []auth.Auth2Role, newValues []auth.Auth2Role) error {
	depts := make([]string, 0) // 部门ID
	roles := make([]string, 0) // 角色ID

	vs := slice.Concat(oldValues, newValues)
	slice.ForEach(vs, func(i int, v auth.Auth2Role) {
		if v.FromTable != constants.TableMenu || v.AuthTable != constants.TableMenu {
			return
		}
		if v.ToTable == constants.TableDept {
			depts = append(depts, v.ToId)
		} else if v.ToTable == constants.TableRole {
			roles = append(roles, v.ToId)
		}
	})
	depts = slice.Unique(depts)
	roles = slice.Unique(roles)

	deptUserIds := make([]string, 0)
	roleUserIds := make([]string, 0)
	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			ld := loader.NewEntityLoader(cfg.OriginDB, loader.EntityLoaderConfig{
				TableName:            constants.TableUser2Dept,
				IdColumn:             entity.ToIdDbName,
				Ids:                  depts,
				DefaultEntityColumns: []string{entity.FromIdDbName},
			})
			return ld.Load(&deptUserIds)
		}).
		Call(func(cl *caller.Caller) error {
			ld := loader.NewEntityLoader(cfg.OriginDB, loader.EntityLoaderConfig{
				TableName:            constants.TableUser2Role,
				IdColumn:             entity.ToIdDbName,
				Ids:                  roles,
				DefaultEntityColumns: []string{entity.FromIdDbName},
			})
			return ld.Load(&roleUserIds)
		}).Err
	if err != nil {
		return err
	}

	uids := slice.Unique(slice.Concat(deptUserIds, roleUserIds))
	slice.ForEach(uids, func(i int, v string) {
		ws.SendMessage(constants.WsKeyUserMenuChanged, "", v)
	})
	return nil
}
