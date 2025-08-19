package history

import "github.com/davycun/eta/pkg/core/service/hook"

func init() {
	hook.AddModifyCallback(hook.CallbackForAll, DeleteHistoryCallback)
}
