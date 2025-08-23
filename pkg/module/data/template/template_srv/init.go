package template_srv

import (
	"github.com/davycun/eta/pkg/common/caller"
	"github.com/davycun/eta/pkg/common/dorm"
	"github.com/davycun/eta/pkg/core/migrate"
	"github.com/davycun/eta/pkg/core/service/hook"
	"github.com/davycun/eta/pkg/core/updater"
	"github.com/davycun/eta/pkg/eta/constants"
)

func InitModule() {
	hook.AddModifyCallback(constants.TableTemplate, modifyCallback)
	migrate.AddCallback(constants.TableTemplate, afterTemplateMigrate)
	migrate.AddCallback(constants.TableTemplateHistory, afterHistoryMigrate)
}

func afterTemplateMigrate(cfg *migrate.MigConfig, pos migrate.CallbackPosition) error {

	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return updater.CreateUpdaterTrigger(cfg.TxDB, constants.TableTemplate)
		}).
		Call(func(cl *caller.Caller) error {
			return dorm.CreateUniqueIndex(cfg.TxDB, constants.TableTemplate, "code")
		}).Err
}
func afterHistoryMigrate(cfg *migrate.MigConfig, pos migrate.CallbackPosition) error {

	return caller.NewCaller().
		Call(func(cl *caller.Caller) error {
			return updater.CreateUpdaterTrigger(cfg.TxDB, constants.TableTemplate)
		}).
		Call(func(cl *caller.Caller) error {
			return dorm.CreateUniqueIndex(cfg.TxDB, constants.TableTemplate, "code")
		}).Err
}
