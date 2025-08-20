package dict_srv

import (
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/eta/plugin/plugin_tree"
	"github.com/davycun/eta/pkg/module/dict"
	"github.com/gin-gonic/gin"
)

func InitModule() {

	//加载回调
	hook.AddModifyCallback(constants.TableDictionary, modifyCallback)
	hook.AddRetrieveCallback(constants.TableDictionary, retrieveCallbacks)
	sqlbd.AddSqlBuilder(constants.TableDictionary, buildListSql, iface.MethodList)

	hook.AddRetrieveCallback(constants.TableDictionary, plugin_tree.TreeResult[dict.Dictionary](), func(option *hook.CallbackOption) {
		option.Order = 10000
	})

	//自我注册默认字典
	dict.Registry(defaultCommonDictionary...)
	dict.Registry(industryCategoryDictionary...)
	dict.Registry(labelCategoriesDictionary...)
	dict.Registry(labelColorDictionary...)

	//添加路由
	controller.Publish(constants.TableDictionary, "/list", controller.ApiConfig{
		Handler: func(srv iface.Service, args any, rs any) error {
			ss := srv.(*Service)
			return ss.Retrieve(args.(*dto.Param), rs.(*dto.Result), iface.MethodList)
		},
	})
	controller.Publish(constants.TableDictionary, "/tree_delete", controller.ApiConfig{
		Handler: func(srv iface.Service, args any, rs any) error {
			ss := srv.(*Service)
			return ss.TreeDelete(args.(*dto.Param), rs.(*dto.Result))
		},
		GetParam: func() any {
			return &dto.Param{ModifyParam: dto.ModifyParam{Data: &dict.Dictionary{}}}
		},
		Binding: func(c *gin.Context, target any) error {
			return controller.BindBodyWithoutValidate(c, target)
		},
	})

}
