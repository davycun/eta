package menu_srv

import (
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/eta/constants"
)

func Router() {

	controller.Publish(constants.TableMenu, "/list", controller.ApiConfig{
		Handler: func(srv iface.Service, args any, rs any) error {
			s := srv.(*Service)
			return s.Retrieve(args.(*dto.Param), rs.(*dto.Result), iface.MethodList)
		},
	})
	controller.Publish(constants.TableMenu, "/my_menu", controller.ApiConfig{
		Handler: func(srv iface.Service, args any, rs any) error {
			s := srv.(*Service)
			return s.RetrieveWrapper(args.(*dto.Param), rs.(*dto.Result), "my_menu", s.MyMenu)
		},
	})
}
