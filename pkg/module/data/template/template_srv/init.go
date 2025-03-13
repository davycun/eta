package template_srv

import (
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
)

func init() {
	hook.AddModifyCallback(constants.TableTemplate, modifyCallback)
}
