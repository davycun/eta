package server

import (
	"github.com/davycun/eta/pkg/common/global"
	"github.com/davycun/eta/pkg/eta/data"
	"github.com/davycun/eta/pkg/eta/middleware"
	"github.com/davycun/eta/pkg/eta/migrator"
	"github.com/davycun/eta/pkg/eta/plugin"
	"github.com/davycun/eta/pkg/eta/router"
	"github.com/davycun/eta/pkg/eta/validator"
	"github.com/davycun/eta/pkg/module"
)

func init() {
	AddLifeCycle(InitConfig, readConfig)
	AddLifeCycle(InitPlugin, plugin.InitPlugin)
	AddLifeCycle(InitData, data.InitData)
	AddLifeCycle(InitApplication, func() error {
		return global.InitApplication(destCfg)
	})
	AddLifeCycle(InitModules, module.InitModules)
	AddLifeCycle(InitValidator, validator.InitValidate)
	AddLifeCycle(InitMiddleware, middleware.InitMiddleware)
	AddLifeCycle(InitEntityConfigRouter, router.InitRouter)
	AddLifeCycle(InitMigrator, migrator.InitMigrator)
	AddLifeCycle(Migrate, func() error {
		return migrator.MigrateLocal(global.GetLocalGorm())
	})
	AddLifeCycle(StartServer, func() error {
		startServer(global.GetApplication())
		return nil
	})
}
