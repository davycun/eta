package dept_srv

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/core/entity"
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/service"
	"gorm.io/gorm"
)

type Service struct {
	service.DefaultService
}

func NewService(c *ctx.Context, db *gorm.DB, tb *entity.Table) iface.Service {
	srv := &Service{}
	srv.SetContext(c)
	srv.SetDB(db)
	srv.SetTable(tb)
	srv.SetUseParamAuth(true)
	srv.SetDisableRetrieveWithES(true)
	return srv
}
