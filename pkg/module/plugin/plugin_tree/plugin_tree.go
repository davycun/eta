package plugin_tree

import (
	"github.com/davycun/eta/pkg/common/dorm/ctype"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
)

// TreeResult
// 注意，在查询类型中，TreeResult这个回调，需要放到最后
// 示例：Order放得足够大，因为树结构处理会调整cfg.Result.Data切片的内容，所以放在最后
//
//	hook.AddRetrieveCallback(hook.CallbackForAll, plugin_tree.TreeResult[Address](), func(option *hook.CallbackOption) {
//			option.Order = 10000
//		})
func TreeResult[E any]() hook.Callback {

	return func(cfg *hook.SrvConfig, pos hook.CallbackPosition) error {

		if pos != hook.CallbackAfter || cfg.CurdType != iface.CurdRetrieve {
			return nil
		}
		//树结构处理
		if cfg.Param.WithTree {
			switch x := cfg.Result.Data.(type) {
			case []E:
				cfg.Result.Data = entity.Tree(cfg.GetDB(), x)
			case []ctype.Map:
				cfg.Result.Data = entity.Tree(cfg.GetDB(), x)
			default:
				return nil
			}
		}
		return nil
	}
}
