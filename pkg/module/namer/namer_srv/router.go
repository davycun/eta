package namer_srv

import (
	"github.com/davycun/eta/pkg/common/global"
)

func InitModule() {
	global.GetGin().POST("/id_name", HandlerIdName)
}
