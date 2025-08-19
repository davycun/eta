package optlog

import (
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
)

func InitModule() {
	hook.AddAuthCallback(constants.TableOperateLog, authRead)
}
