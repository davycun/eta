package app_srv

import (
	"github.com/davycun/eta/pkg/core/migrate/mig_hook"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/eta/constants"
)

func init() {
	hook.AddModifyCallback(constants.TableApp, modifyCallback)
	hook.AddRetrieveCallback(constants.TableApp, retrieveCallbacks)
	mig_hook.AddCallback(constants.TableApp, afterMigrate)
}
