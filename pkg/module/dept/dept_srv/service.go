package dept_srv

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"gorm.io/gorm"
)

type Service struct {
	service.DefaultService
}

func NewService(c *ctx.Context, db *gorm.DB, ec *iface.EntityConfig) iface.Service {
	srv := &Service{}
	_ = srv.Init(c, db, ec)
	srv.SetUseParamAuth(true)
	srv.SetDisableRetrieveWithES(true)
	return srv
}
