package router

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/iface"
)

type (
	RouteFunc func()
)

var (
	routerFuncList = make([]RouteFunc, 0)
)

// Registry 注册路由函数，一个简单的空函数，实际添加路由由调用者自行决定
func Registry(rf RouteFunc) {
	routerFuncList = append(routerFuncList, rf)
}

// InitRouter 录入注册入口
// 在cmd/server.go中调用之后，都会调用本包的init()函数完成自动注册
func InitRouter() {
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
}
