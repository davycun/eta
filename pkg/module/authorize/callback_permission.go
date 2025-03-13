package authorize

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/authorize/auth"
	"github.com/duke-git/lancet/v2/slice"
)

func modifyCallbackPermission(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []auth.Permission) error {
				return fillField(newValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterUpdate(cfg, pos, func(cfg *hook.SrvConfig, oldValues []auth.Permission, newValues []auth.Permission) error {
				return delCachePermission(oldValues)
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []auth.Permission) error {
				return delCachePermission(oldValues)
			})
		}).Err

	return err
}

func fillField(dt []auth.Permission) error {

	var (
		err            error
		fillPermission = func(perm *auth.Permission) error {
			if perm.TbName == "" {
				return errs.NewClientError("tb_name 不能为空")
			}

			if len(perm.Filters) < 1 && len(perm.RecursiveFilters) < 1 && perm.Type == auth.PermissionFilter {
				return errs.NewClientError("filters 或者 RecursiveFilters 不能都为空")
			}

			if (len(perm.Filters) > 0 || len(perm.RecursiveFilters) > 0) && perm.Type == "" {
				perm.Type = auth.PermissionFilter
			}
			return nil
		}
	)
	for i, _ := range dt {
		err = errs.Cover(err, fillPermission(&dt[i]))
	}
	return err
}

func delCachePermission(data []auth.Permission) error {
	slice.ForEach(data, func(index int, item auth.Permission) {
		auth.DelPermissionCache(item.ID)
	})
	return nil
}
