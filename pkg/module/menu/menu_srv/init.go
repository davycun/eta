package menu_srv

import (
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/eta/plugin/plugin_tree"
	"github.com/davycun/eta/pkg/module/menu"
)

func InitModule() {
	hook.AddModifyCallback(constants.TableMenu, modifyCallback)
	hook.AddRetrieveCallback(constants.TableMenu, retrieveCallback)
	sqlbd.AddSqlBuilder(constants.TableMenu, buildListSql, iface.MethodList)

	hook.AddRetrieveCallback(constants.TableMenu, plugin_tree.TreeResult[menu.Menu](), func(option *hook.CallbackOption) {
		option.Order = 10000
	})

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
