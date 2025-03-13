package dept_srv

import (
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/eta/constants"
	"github.com/davycun/eta/pkg/module/dept"
)

func Router() {
	controller.Publish(constants.TableDept, "/list", controller.ApiConfig{
		Handler: func(srv iface.Service, args any, rs any) error {
			ss := srv.(*Service)
			return ss.Retrieve(args.(*dto.Param), rs.(*dto.Result), iface.MethodList)
		},
		GetParam: func() any {
			return &dto.Param{ModifyParam: dto.ModifyParam{Data: &dept.Department{}}}
		},
		GetResult: func() any {
			return &dto.Result{}
		},
	})
}
