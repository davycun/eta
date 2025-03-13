package user_srv

import (
	"github.com/davycun/eta/pkg/core/iface"
	"github.com/davycun/eta/pkg/core/migrate/mig_hook"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/service/sqlbd"
	"github.com/davycun/eta/pkg/eta/constants"
)

func init() {
	hook.AddModifyCallback(constants.TableApp, modifyCallbackApp)
	hook.AddModifyCallback(constants.TableUser, modifyCallbackUser)
	hook.AddModifyCallback(constants.TableUser2Dept, modifyCallbackUser2Dept)
	hook.AddModifyCallback(constants.TableUser2App, modifyCallbackUser2App)
	hook.AddModifyCallback(constants.TableDept, modifyCallbackDept)
	hook.AddRetrieveCallback(constants.TableUser, retrieveCallbackUser)

	sqlbd.AddSqlBuilder(constants.TableUser, buildListSql, iface.MethodList)
	mig_hook.AddCallback(constants.TableUser, afterMigratorUser)
	mig_hook.AddCallback(constants.TableUser2App, afterMigratorUser2App)
}
