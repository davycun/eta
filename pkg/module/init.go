package module

import (
	"github.com/davycun/eta/pkg/core/service/ecf"
	"github.com/davycun/eta/pkg/eta/router"
)

func InitModules() {
	ecf.Registry(entityConfig()...)
	router.Registry(Router)
	initSetting()
}
