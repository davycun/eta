package reload

import (
	"github.com/davycun/eta/pkg/common/global"
)

func InitModule() {

	group := global.GetGin().Group("/reload")
	group.POST("/db2es", NewReloadController(RdTypeDb2Es))

}
