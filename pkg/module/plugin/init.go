package plugin

import (
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/module/app"
	"github.com/davycun/eta/pkg/module/plugin/plugin_auth"
	"github.com/davycun/eta/pkg/module/plugin/plugin_crypt"
	"github.com/davycun/eta/pkg/module/plugin/plugin_es"
	"github.com/davycun/eta/pkg/module/plugin/plugin_geo"
	"github.com/davycun/eta/pkg/module/plugin/plugin_his"
	"github.com/davycun/eta/pkg/module/plugin/plugin_push"
	"github.com/davycun/eta/pkg/module/plugin/plugin_tree"
)

func InitPlugin() {
	hook.AddModifyCallback(hook.CallbackForAll, plugin_es.ModifyCallbackForEs)
	hook.AddModifyCallback(hook.CallbackForAll, plugin_crypt.StoreSign)
	hook.AddModifyCallback(hook.CallbackForAll, plugin_crypt.StoreEncrypt)
	hook.AddModifyCallback(hook.CallbackForAll, plugin_push.PublishModifyCallbacks)
	hook.AddRetrieveCallback(hook.CallbackForAll, plugin_geo.ProcessGeometryResult)
	hook.AddRetrieveCallback(hook.CallbackForAll, plugin_crypt.EncryptQueryParam)
	hook.AddRetrieveCallback(hook.CallbackForAll, plugin_crypt.VerifySign)
	hook.AddRetrieveCallback(hook.CallbackForAll, plugin_crypt.StoreDecrypt)
	hook.AddRetrieveCallback(hook.CallbackForAll, plugin_his.DeleteHistoryCallback)
	hook.AddRetrieveCallback(hook.CallbackForAll, plugin_tree.TreeResult[app.App](), func(option *hook.CallbackOption) {
		option.Order = 10000
	})
	hook.AddAuthCallback(hook.CallbackForAll, plugin_auth.AuthFilter, func(option *hook.CallbackOption) {
		option.IsAuth = true
		option.Methods = []iface.Method{iface.MethodAll}
		option.Order = 1
	})
}
