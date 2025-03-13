package migrator

import (
	"github.com/davycun/eta/pkg/common/ctx"
	"github.com/davycun/eta/pkg/common/logger"
	"github.com/davycun/eta/pkg/core/dto"
	"gorm.io/gorm"
)

type (
	CallbackPosition int
	Callback         func(cfg *MigConfig, pos CallbackPosition) error
	MigrateType      int //主要区分回调是针对local还是app
)

const (
	BeforeCallback CallbackPosition = 1
	AfterCallback  CallbackPosition = 2
	MigLocal       MigrateType      = 1 //migrate本地DB
	MigApp         MigrateType      = 2 //migrate对应的APP DB
	MigAll         MigrateType      = 3
)

var (
	migratorLocalCallbacks = make([]Callback, 0)
	migratorAppCallbacks   = make([]Callback, 0)
)

func AddCallback(migType MigrateType, cb Callback) {
	switch migType {
	case MigLocal:
		migratorLocalCallbacks = append(migratorLocalCallbacks, cb)
	case MigApp:
		migratorAppCallbacks = append(migratorAppCallbacks, cb)
	case MigAll:
		migratorLocalCallbacks = append(migratorLocalCallbacks, cb)
		migratorAppCallbacks = append(migratorAppCallbacks, cb)
	default:
		logger.Errorf("add callback not support migrate type %d", migType)
	}
}

func NewMigConfig(c *ctx.Context, db *gorm.DB, param *dto.Param) *MigConfig {
	return &MigConfig{
		TxDB:  db,
		C:     c,
		Param: param,
	}
}

type MigConfig struct {
	Param *dto.Param
	TxDB  *gorm.DB
	C     *ctx.Context
}

func (mc *MigConfig) BeforeMigrateApp() error {
	return mc.callMigrate(MigApp, BeforeCallback)
}
func (mc *MigConfig) BeforeMigrateLocal() error {
	return mc.callMigrate(MigLocal, BeforeCallback)
}
func (mc *MigConfig) AfterMigrateApp() error {
	return mc.callMigrate(MigApp, AfterCallback)
}
func (mc *MigConfig) AfterMigrateLocal() error {
	return mc.callMigrate(MigLocal, AfterCallback)
}

func (mc *MigConfig) callMigrate(migType MigrateType, pos CallbackPosition) error {
	for _, fc := range mc.getCallback(migType) {
		err := fc(mc, pos)
		if err != nil {
			return err
		}
	}
	return nil
}

func (mc *MigConfig) getCallback(migType MigrateType) []Callback {
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
