package router

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/common/utils"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/service/ecf"
	"github.com/davycun/eta/pkg/module"
)

type (
	RouteFunc func()
)

var (
	routerFuncList = []RouteFunc{module.Router}
)

// Registry 注册路由函数，一个简单的空函数，实际添加路由由调用者自行决定
func Registry(rf RouteFunc) {
	routerFuncList = append(routerFuncList, rf)
}

// InitRouter 录入注册入口
// 在cmd/server.go中调用之后，都会调用本包的init()函数完成自动注册
func InitRouter() {
	module.RegistryEntityConfig() //加载配置

	cfg := global.GetConfig()
	ecList := ecf.GetEntityConfigList()
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
