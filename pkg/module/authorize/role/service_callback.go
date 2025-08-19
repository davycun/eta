package role

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/errs"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/duke-git/lancet/v2/slice"
)

func modifyCallbacks(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.BeforeCreate(cfg, pos, func(config *hook.SrvConfig, newValues []Role) error {
				for _, v := range newValues {
					if v.Category == "" {
						v.Category = CategoryOther
					}
				}
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			// 固定角色不能删除
			return hook.BeforeDelete(cfg, pos, func(cfg *hook.SrvConfig, oldValues []Role) error {
				fixedIds := slice.Map(defaultRoleList, func(i int, v Role) string {
					return v.ID
				})
				oldIds := slice.Map(oldValues, func(i int, v Role) string {
					return v.ID
				})
				if utils.ContainAny(oldIds, fixedIds...) {
					return errs.NewClientError("不能删除固定角色")
				}
				return nil
			})
		}).Err

	return err
}
