package user_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/user"
)

func retrieveCallbackUser(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

	err := caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			//return service.BeforeRetrieve(cfg, pos, func(cfg *hook.SrvConfig) error {
			//	if cfg.Ctx.GetContextIsManager() {
			//		cfg.Param.Filters = []filter.Filter{{
			//			LogicalOperator: filter.And,
			//			Column:          "app_id",
			//			Operator:        filter.Eq,
			//			Value:           cfg.Ctx.GetContextAppId(),
			//			Filters: []filter.Filter{{
			//				LogicalOperator: filter.And,
			//				Filters:         cfg.Param.Filters,
			//			}},
			//		}}
			//	}
			//	return nil
			//})
			return nil
		}).
		Call(func(cl *caller.Caller) error {
			return hook.AfterRetrieveAny(cfg, pos, func(cfg *hook.SrvConfig) error {
				switch listRs := cfg.Result.Data.(type) {
				case []user.ListResult:
					return fill(cfg, listRs)
				case *[]user.ListResult:
					return fill(cfg, *listRs)
				case []user.User:
					return fill(cfg, listRs)
				case *[]user.User:
					return fill(cfg, *listRs)
				case []ctype.Map:
					return fill(cfg, listRs)
				case *[]ctype.Map:
					return fill(cfg, *listRs)
				}
				return nil
			})
		}).Err

	return err
}
