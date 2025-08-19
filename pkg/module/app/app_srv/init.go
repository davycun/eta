package app_srv

import (
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/migrate/mig_hook"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/eta/migrator"
)

func InitModule() {
	hook.AddModifyCallback(constants.TableApp, modifyCallback)
	hook.AddRetrieveCallback(constants.TableApp, retrieveCallbacks)
	mig_hook.AddCallback(constants.TableApp, afterMigrate)

	controller.Publish(constants.TableApp, "/migrate", controller.ApiConfig{
		Handler: func(srv iface.Service, args any, rs any) error {
			return srv.(*Service).Migrate(args.(*migrator.MigrateAppParam), rs.(*dto.Result))
		},
		GetParam: func() any {
			return &migrator.MigrateAppParam{}
		},
	})
}
