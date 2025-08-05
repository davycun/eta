package eta

import (
	"github.com/davycun/eta/pkg/eta/validator"
	"github.com/davycun/eta/pkg/module"
)

func InitEta() {
	//初始化模块需放第一
	module.InitModules()
	validator.AddValidate()
}
