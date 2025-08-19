package migrator

import (
	"github.com/davycun/eta/pkg/common/logger"
)

const (
	BeforeCallback CallbackPosition = 1
	AfterCallback  CallbackPosition = 2
	MigLocal       MigrateType      = 1 //migrate本地DB
	MigApp         MigrateType      = 2 //migrate对应的APP DB
	MigAll         MigrateType      = 3
)

type (
	CallbackPosition int
	Callback         func(cfg *CallbackCaller, pos CallbackPosition) error
	MigrateType      int //主要区分回调是针对local还是app
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
