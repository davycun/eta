package setting_srv

import (
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
)

func InitModule() {
	hook.AddModifyCallback(constants.TableSetting, modifyCallback)
}
