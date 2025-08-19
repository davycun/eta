package migrator

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dto"
	"gorm.io/gorm"
)

func NewCallbackCaller(c *ctx.Context, db *gorm.DB, param *dto.Param) *CallbackCaller {
	return &CallbackCaller{
		TxDB:  db,
		C:     c,
		Param: param,
	}
}

type CallbackCaller struct {
	Param *dto.Param
	TxDB  *gorm.DB
	C     *ctx.Context
}

func (mc *CallbackCaller) BeforeMigrateApp() error {
	return mc.callMigrate(MigApp, BeforeCallback)
}
func (mc *CallbackCaller) BeforeMigrateLocal() error {
	return mc.callMigrate(MigLocal, BeforeCallback)
}
func (mc *CallbackCaller) AfterMigrateApp() error {
	return mc.callMigrate(MigApp, AfterCallback)
}
func (mc *CallbackCaller) AfterMigrateLocal() error {
	return mc.callMigrate(MigLocal, AfterCallback)
}

func (mc *CallbackCaller) callMigrate(migType MigrateType, pos CallbackPosition) error {
	for _, fc := range mc.getCallback(migType) {
		err := fc(mc, pos)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mc *CallbackCaller) getCallback(migType MigrateType) []Callback {
	switch migType {
	case MigLocal:
		return migratorLocalCallbacks
	case MigApp:
		return migratorAppCallbacks
	default:
		logger.Errorf("get callback not support migrate type %d", migType)
		return []Callback{}
	}
}
