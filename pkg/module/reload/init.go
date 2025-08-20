package reload

import (
	"github.com/davycun/eta/pkg/common/dsync"
	"github.com/davycun/eta/pkg/core/controller"
	"github.com/davycun/eta/pkg/core/dto"
	"github.com/davycun/eta/pkg/core/iface"
)

func InitModule() {

	controller.Publish("reload", "/db2es", controller.ApiConfig{
		Handler: func(srv iface.Service, args any, rs any) error {
			ss := srv.(*Service)
			return ss.ReloadDb2Es(args.(*dto.Param), rs.(*dto.Result))
		},
		GetParam: dto.NewParamWithExtra[dsync.SyncOption](),
	})

}
