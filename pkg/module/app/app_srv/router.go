package app_srv

import (
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/eta/migrator"
)

func Router() {
	controller.Publish(constants.TableApp, "/migrate", controller.ApiConfig{
		Handler: func(srv iface.Service, args any, rs any) error {
			return srv.(*Service).Migrate(args.(*migrator.MigrateAppParam), rs.(*dto.Result))
		},
		GetParam: func() any {
			return &migrator.MigrateAppParam{}
		},
	})
}
