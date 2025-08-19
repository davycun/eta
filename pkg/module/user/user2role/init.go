package user2role

import (
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
)

func InitModule() {
	hook.AddModifyCallback(constants.TableRole, modifyCallback)
}
