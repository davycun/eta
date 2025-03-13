package dept_srv

import (
	"errors"
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/dept"
	"github.com/davycun/eta/pkg/module/namer"
	"github.com/duke-git/lancet/v2/slice"
)

func deptCallbacks(cfg *hook.SrvConfig, pos hook.CallbackPosition) (err error) {

	err = caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return hook.BeforeDelete(cfg, pos, validateHasChildBeforeDelete)
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterCreate(cfg, pos, func(cfg *hook.SrvConfig, newValues []dept.Department) error {
				namer.DelIdNameCacheByContext(cfg.Ctx)
				return nil
			})
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterModify(cfg, pos, func(cfg *hook.SrvConfig, oldValues []dept.Department, newValues []dept.Department) error {

				ids := make([]string, 0, len(oldValues))
				slice.ForEach(oldValues, func(index int, item dept.Department) {
					ids = append(ids, item.ID)
				})
				err = dept.DelDeptAndUser2DeptCache(cfg.OriginDB, ids...)
				namer.DelIdNameCacheByContext(cfg.Ctx)
				return nil
			}, iface.MethodUpdate, iface.MethodUpdateByFilters, iface.MethodDelete, iface.MethodDeleteByFilters)
		}).Err
	return err
}

func validateHasChildBeforeDelete(cfg *hook.SrvConfig, deps []dept.Department) error {
	// 如果有子部门不能删除
	var (
		childId   string
		parentIds = make([]string, 0)
		db        = cfg.OriginDB
	)
	for _, item := range deps {
		parentIds = append(parentIds, item.ID)
	}
	err := db.Model(&dept.Department{}).Select("id").Where(`"parent_id" in ?`, parentIds).Limit(1).Find(&childId).Error
	if err != nil {
		return err
	}
	if childId != "" {
		return errors.New("当前部门存在子部门不允许删除,请使用tree_delete接口")
	}
	return nil
}

func retrieveCallbacks(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {
	return hook.AfterRetrieveAny(cfg, pos, func(cfg *hook.SrvConfig) error {
		switch listRs := cfg.Result.Data.(type) {
		case []dept.Department:
			return fill(cfg, listRs)
		case []ctype.Map:
			logger.Errorf("部门fill暂未实现ctype.Map填充")
			return nil
		}
		return nil
	})
}
