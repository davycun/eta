package user2app

import (
	"github.com/davycun/eta/pkg/core/migrate"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
)

func InitModule() {
	migrate.AddCallback(constants.TableUser2App, afterMigratorUser2App)
	hook.AddModifyCallback(constants.TableUser2App, modifyCallbackUser2App)
}
