package router

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/iface"
)

// InitRouter 录入注册入口
// 在cmd/server.go中调用之后，都会调用本包的init()函数完成自动注册
func InitRouter() error {
	cfg := global.GetConfig()
	ecList := iface.GetEntityConfigList()
	for _, v := range ecList {
		if v.DisableApi {
			continue
		}
		if len(cfg.Server.RouterPkg) > 0 {
			if utils.ContainAny(cfg.Server.RouterPkg, v.Namespace) {
				controller.Registry(v)
			}
		} else {
			controller.Registry(v)
		}
	}

	for _, rf := range routerFuncList {
		rf()
	}
	return nil
}
