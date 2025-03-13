package dict_srv

import (
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/dict"
	"github.com/gin-gonic/gin"
)

func Router() {

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
