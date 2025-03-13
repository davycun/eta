package authorize

import (
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
)

func init() {
	hook.AddModifyCallback(constants.TableRole, afterDeleteRoleDeleteAuth2Role)
	hook.AddModifyCallback(constants.TableDept, afterDeleteDeptDeleteAuth2Role)
	hook.AddModifyCallback(constants.TableAuth2Role, modifyCallbackAuth2Role)
	hook.AddModifyCallback(constants.TablePermission, modifyCallbackPermission)
}
