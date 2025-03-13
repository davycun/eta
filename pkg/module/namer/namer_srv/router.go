package namer_srv

import (
	"github.com/davycun/eta/pkg/common/global"
)

func Router() {
	global.GetGin().POST("/id_name", HandlerIdName)
}
