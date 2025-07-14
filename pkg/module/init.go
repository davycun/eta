package module

import (
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/eta/router"
)

func InitModules() {
	iface.Registry(entityConfig()...)
	router.Registry(Router)
	initSetting()
}
